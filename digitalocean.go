package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/digitalocean/godo"
)

func doHandler(w http.ResponseWriter, r *http.Request) {

	client := r.Header.Get("X-Forwarded-For")
	if client == "" {
		client = r.RemoteAddr
	}

	ip, hostname, err := checkForms(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger("client: "+client+" "+err.Error(), "err")
		return
	}

	// DigitalOcean DNS API is straightforward like Cloudflare
	err = doDoUpdate(user, pass, hostname, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger("client: "+client+" "+err.Error(), "err")
	} else {
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
		logger("client: "+client+" successfully updated DigitalOcean DNS hostname: "+hostname+" ip: "+ip, "info")
	}
}

func doDoUpdate(domainName, apiToken, hostname, ip string) error {
	client := godo.NewFromToken(apiToken)

	// Extract domain from hostname if not provided separately
	// For DO, we expect the domain name to be passed as username
	domain := domainName

	// Get the record name (subdomain part)
	recordName := hostname
	if strings.Contains(hostname, domain) {
		// Remove the domain part to get just the record name
		recordName = strings.TrimSuffix(hostname, "."+domain)
		if recordName == domain {
			recordName = "@" // root domain
		}
	}

	// Check if DNS Record exists
	records, _, err := client.Domains.Records(context.Background(), domain, nil)
	if err != nil {
		return errors.New("failed to list DNS records or domain not found")
	}

	var existingRecord *godo.DomainRecord
	for _, record := range records {
		if record.Name == recordName && record.Type == "A" {
			existingRecord = &record
			break
		}
	}

	if existingRecord != nil {
		// Update existing record
		editRequest := &godo.DomainRecordEditRequest{
			Type: "A",
			Name: recordName,
			Data: ip,
			TTL:  300,
		}
		_, _, err = client.Domains.EditRecord(context.Background(), domain, existingRecord.ID, editRequest)
		if err != nil {
			return err
		}
	} else {
		// Create new record
		createRequest := &godo.DomainRecordEditRequest{
			Type: "A",
			Name: recordName,
			Data: ip,
			TTL:  300,
		}
		_, _, err = client.Domains.CreateRecord(context.Background(), domain, createRequest)
		if err != nil {
			return err
		}
	}

	return nil
}
