package main

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const DefaultConfigurationFilepath = "../configs/control-plane.yaml"

type ScannerConfiguration struct {
	Image string `json:"image"`
}

type ScannersConfigurations struct {
	Trivy ScannerConfiguration `json:"trivy"`
	Grype ScannerConfiguration `json:"grype"`
}

type Configuration struct {
	Scanners ScannersConfigurations `json:"scanners"`
}

func GetConfiguration() (*Configuration, error) {
	filepath, ok := os.LookupEnv("APP_CONFIG_FILE")
	if !ok {
		filepath = DefaultConfigurationFilepath
	}

	f, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("cannot open configuration file: %w", err)
	}

	var cfg Configuration
	if err = yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("cannot parse configuration file: %w", err)
	}

	return &cfg, nil
}
