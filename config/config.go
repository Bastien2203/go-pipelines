package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Name string `yaml:"name"`
	Repo struct {
		URL    string `yaml:"url"`
		Branch string `yaml:"branch"`
	} `yaml:"repo"`
	Registry struct {
		URL         string `yaml:"url"`
		Username    string `yaml:"username"`
		PasswordEnv string `yaml:"password_env"`
		ImageName   string `yaml:"image_name"`
	} `yaml:"registry"`
}

const CONFIG_DIR = "./configs"

func GetConfig(name string) (*Config, error) {
	data, err := os.ReadFile(fmt.Sprintf("%s/%s.yml", CONFIG_DIR, name))
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
