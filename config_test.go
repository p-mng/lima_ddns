package main

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	want := &Config{
		Interval: 300,
		Token:    "cYcdrwXvxGpqmec3R464VEHaAL5F7Vt2",
		Domains: []DomainConfig{
			{
				ID: 27272,
				Records: []RecordConfig{
					{
						Type:     "A",
						Content:  "*.example.com",
						Priority: 10,
						TTL:      3600,
					},
					{
						Type:     "A",
						Content:  "example.com",
						Priority: 0,
						TTL:      300,
					},
					{
						Type:     "AAAA",
						Content:  "example.com",
						Priority: 0,
						TTL:      300,
					},
				},
			},
		},
	}

	if err := os.Setenv("CONFIG_FILE", "config.sample.json"); err != nil {
		t.Errorf("failed to set environment variable: %s", err.Error())
	}

	gotJSON, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() failed with error: %s", err.Error())
	}
	if !reflect.DeepEqual(gotJSON, want) {
		t.Errorf("LoadConfig() is %+v, want %+v", gotJSON, want)
	}

	if err := os.Setenv("CONFIG_FILE", "config.sample.yaml"); err != nil {
		t.Errorf("failed to set environment variable: %s", err.Error())
	}

	gotYaml, err := LoadConfig()
	if err != nil {
		t.Errorf("LoadConfig() failed with error: %s", err.Error())
	}
	if !reflect.DeepEqual(gotYaml, want) {
		t.Errorf("LoadConfig() is %+v, want %+v", gotYaml, want)
	}
}

func TestLoadConfigInvalidPath(t *testing.T) {
	if err := os.Setenv("CONFIG_FILE", "config.sample.toml"); err != nil {
		t.Errorf("failed to set environment variable: %s", err.Error())
	}

	if _, err := LoadConfig(); err == nil {
		t.Errorf("LoadConfig() should fail")
	}
}
