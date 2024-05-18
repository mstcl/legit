package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Repo struct {
		ScanPath   string   `yaml:"scanPath"`
		Readme     []string `yaml:"readme"`
		MainBranch []string `yaml:"mainBranch"`
		Ignore     []string `yaml:"ignore,omitempty"`
	} `yaml:"repo"`
	Meta struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
	} `yaml:"meta"`
	Server struct {
		Name string `yaml:"name,omitempty"`
		Host string `yaml:"host"`
		Port int    `yaml:"port"`
	} `yaml:"server"`
}

func Read(f string) (*Config, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	c := Config{}
	if err := yaml.Unmarshal(b, &c); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	if c.Repo.ScanPath, err = filepath.Abs(c.Repo.ScanPath); err != nil {
		return nil, err
	}

	return &c, nil
}
