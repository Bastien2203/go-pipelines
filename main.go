package main

import (
	"context"
	"fmt"
	"go-pipelines/config"
	"go-pipelines/docker"
	"go-pipelines/git"
	"go-pipelines/queue"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	buildTimeout, err := strconv.Atoi(os.Getenv("BUILD_TIMEOUT"))
	if err != nil {
		log.Fatalf("error while retrieving env var BUILD_TIMEOUT : %v", err)
	}

	port := os.Getenv("WEBHOOK_PORT")
	if port == "" {
		log.Fatalf("Unable to retrieve WEBHOOK_PORT env var ")
	}

	worker := queue.NewWorker(buildTimeout, 50, runBuildProcess)
	go worker.Start()

	http.HandleFunc("/webhook/{name}", func(w http.ResponseWriter, r *http.Request) { handleWebhook(w, r, worker) })
	log.Printf("Listening webhooks on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))

}

func handleWebhook(w http.ResponseWriter, r *http.Request, worker *queue.Worker) {
	log.Printf("Webhook triggered")
	name := r.PathValue("name")

	if name == "" {
		http.Error(w, "missing project name", http.StatusBadRequest)
		return
	}

	cfg, err := config.GetConfig(name)

	if err != nil {
		http.Error(w, "project not configured", http.StatusNotFound)
		return
	}

	job := queue.Job{
		ID:         uuid.New().String(),
		Config:     cfg,
		ReceivedAt: time.Now(),
	}

	select {
	case worker.Queue <- job:
		log.Printf("[Queue] Job added: %s (Project: %s)", job.ID, cfg.Name)
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(fmt.Sprintf(`{"status": "queued", "job_id": "%s"}`, job.ID)))
	default:
		log.Printf("[Queue] Full! Rejected job for %s", cfg.Name)
		http.Error(w, "Server busy, queue full", http.StatusServiceUnavailable)
	}
}

func runBuildProcess(ctx context.Context, j queue.Job) {
	cfg := j.Config

	log.Printf("[%s] Build started", cfg.Name)

	start := time.Now()
	workdir := fmt.Sprintf("./tmp/%s", cfg.Name)

	if err := os.RemoveAll(workdir); err != nil {
		log.Printf("[%s] Error cleaning workspace: %v", cfg.Name, err)
		return
	}

	if err := os.MkdirAll(workdir, 0755); err != nil {
		log.Printf("[%s] Error creating workspace: %v", cfg.Name, err)
		return
	}

	if err := git.Clone(ctx, cfg.Repo.Branch, cfg.Repo.URL, workdir); err != nil {
		log.Printf("[%s] Git Clone failed: %v", cfg.Name, err)
		return
	}

	versionTag, err := git.LatestTag(ctx, workdir)
	if err != nil {
		log.Printf("[%s] No git tag found, using 'dev'", cfg.Name)
		versionTag = "dev"
	}

	imageName := cfg.Registry.ImageName
	baseTag := fmt.Sprintf("%s:latest", imageName)

	if err := docker.Build(ctx, baseTag, workdir); err != nil {
		log.Printf("[%s] Docker Build failed: %v", cfg.Name, err)
		return
	}

	if err := docker.Login(ctx, cfg.Registry.URL, cfg.Registry.Username, cfg.Registry.PasswordEnv); err != nil {
		log.Printf("[%s] Docker Login failed: %v", cfg.Name, err)
		return
	}

	tagsToPush := []string{"latest", versionTag}
	for _, tag := range tagsToPush {
		fullTag := fmt.Sprintf("%s:%s", imageName, tag)

		if tag != "latest" {
			log.Printf("Tagging image as %s...", fullTag)
			err := docker.Tag(ctx, baseTag, fullTag)
			if err != nil {
				log.Printf("[%s] Error tagging %s: %v", cfg.Name, tag, err)
				continue
			}
		}

		log.Printf("Pushing image : %s", fullTag)
		err := docker.Push(ctx, fullTag)
		if err != nil {
			log.Printf("[%s] Failed to push %s: %v", cfg.Name, fullTag, err)
			return
		} else {
			log.Printf("[%s] Pushed %s", cfg.Name, fullTag)
		}
	}

	duration := time.Since(start)
	log.Printf("[%s] Build finished successfully in %v", cfg.Name, duration)
}
