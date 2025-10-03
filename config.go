package main

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Domains  []DomainConfig `json:"domains"`
	Interval int            `json:"interval"`
	Token    string         `json:"token"`
}

type DomainConfig struct {
	ID      int            `json:"id"`
	Records []RecordConfig `json:"records"`
}

type RecordConfig struct {
	Content  string `json:"content"`
	Priority int    `json:"priority,omitempty"`
	TTL      int    `json:"ttl"`
	Type     string `json:"type"`
}

func LoadConfig() (*Config, error) {
	filename := "/config.yaml"
	if os.Getenv("CONFIG_FILE") != "" {
		filename = os.Getenv("CONFIG_FILE")
	}

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			log.Printf("Warning: Error closing config file: %s", err.Error())
		}
	}(file)

	conf := &Config{}
	err = yaml.NewDecoder(file).Decode(conf)
	if err != nil {
		return nil, err
	}

	for i, domain := range conf.Domains {
		for k, record := range conf.Domains[i].Records {
			if record.Type == "" || (record.Type != "A" && record.Type != "AAAA") {
				log.Printf("Warning: Record %+v of domain %d has no or an invalid type. Defaulting to A.", record, domain.ID)
				conf.Domains[i].Records[k].Type = "A"
			}
			if record.Content == "" {
				log.Printf("Warning: Record %+v of domain %d has no content. Defaulting to @.", record, domain.ID)
				conf.Domains[i].Records[k].Content = "@"
			}
			if record.TTL == 0 {
				log.Printf("Warning: TTL for record %+v of domain %d is 0, defaulting to 3600.", record, domain.ID)
				conf.Domains[i].Records[k].TTL = 3600
			}
		}
	}

	return conf, nil
}
