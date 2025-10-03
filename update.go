package main

import (
	"fmt"
	"log"
)

func Update(config *Config) error {
	var fetchv4, fetchv6 bool
	for _, domain := range config.Domains {
		for _, record := range domain.Records {
			switch record.Type {
			case "A":
				fetchv4 = true
			case "AAAA":
				fetchv6 = true
			}
		}
	}

	var ip4, ip6 string
	var err error
	if fetchv4 {
		ip4, err = IPv4()
		if err != nil {
			return err
		}
	}
	if fetchv6 {
		ip6, err = IPv6()
		if err != nil {
			return err
		}
	}

	for _, domain := range config.Domains {
		log.Printf("Updating domain %d", domain.ID)
		err := UpdateDomain(&domain, ip4, ip6, config.Token)
		if err != nil {
			return err
		}
	}
	return nil
}

func UpdateDomain(config *DomainConfig, ip4 string, ip6 string, token string) error {
	records, err := GetRecords(config.ID, token)
	if err != nil {
		return err
	}

	for _, record := range config.Records {
		var ip string
		if record.Type == "A" {
			ip = ip4
		} else {
			ip = ip6
		}
		if ip == "" {
			log.Printf("No valid IP address for record %s of type %s", record.Content, record.Type)
		}
		r := RecordExists(&record, records)
		if r == nil {
			log.Printf("Record %s of type %s does not exist, creating", record.Content, record.Type)
			_, err := CreateRecord(
				record.Content,
				record.Type,
				ip,
				fmt.Sprintf("%d", record.Priority),
				record.TTL,
				config.ID,
				token,
			)
			if err != nil {
				return err
			}
			continue
		}

		update := false

		if r.Content != ip {
			update = true
		}
		if r.Priority != record.Priority {
			update = true
		}
		if r.TTL != record.TTL {
			update = true
		}

		if !update {
			log.Printf("Record %s of type %s is up to date", record.Content, record.Type)
			continue
		}

		log.Printf("Record %s of type %s needs updating (old: %s, new: %s)", record.Content, record.Type, r.Content, ip)
		_, err := UpdateRecord(
			record.Content,
			record.Type,
			ip,
			fmt.Sprintf("%d", record.Priority),
			record.TTL,
			config.ID,
			r.ID,
			token,
		)
		if err != nil {
			return err
		}
		log.Printf("Record %s of type %s updated", record.Content, record.Type)
	}
	return nil
}

func RecordExists(config *RecordConfig, records *DNSRecordResponse) *Record {
	for _, record := range records.Records {
		if record.Name == config.Content && record.Type == config.Type {
			return &record
		}
	}
	return nil
}
