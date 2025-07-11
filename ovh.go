package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// OVH DNS Record structure
type OVHDNSRecord struct {
	ID        int64  `json:"id,omitempty"`
	Zone      string `json:"zone,omitempty"`
	SubDomain string `json:"subDomain"`
	FieldType string `json:"fieldType"`
	Target    string `json:"target"`
	TTL       int    `json:"ttl,omitempty"`
}

// OVH API Client structure
type OVHClient struct {
	Endpoint          string
	ApplicationKey    string
	ApplicationSecret string
	ConsumerKey       string
	httpClient        *http.Client
}

func ovhHandler(w http.ResponseWriter, r *http.Request) {
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

	// URL format: /ovh/endpoint/domain.com/appkey/?ip=x.x.x.x&hostname=host.domain.com
	pathComponents := strings.Split(r.URL.Path, "/")
	if len(pathComponents) < 5 || pathComponents[2] == "" || pathComponents[3] == "" || pathComponents[4] == "" {
		http.Error(w, "invalid path format - expected /ovh/endpoint/domain.com/appkey/", http.StatusBadRequest)
		logger("client: "+client+" invalid path format", "err")
		return
	}

	endpoint := pathComponents[2]
	domain := pathComponents[3]
	appKey := pathComponents[4]

	if endpoint == "" || domain == "" || appKey == "" {
		http.Error(w, "endpoint, domain and application key are required in URL path", http.StatusBadRequest)
		logger("client: "+client+" missing endpoint, domain or application key in URL path", "err")
		return
	}

	// Setup OVH API Client
	ovhClient, err := ovhSetup(appKey, user, pass, endpoint)
	if err != nil {
		logger("client: "+client+" "+err.Error(), "err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Add debug logging
	logger("client: "+client+" OVH request - endpoint: "+endpoint+" domain: "+domain+" appkey: "+appKey+" user: "+user, "info")

	// Update DNS record
	err = ovhUpdateDNS(ovhClient, domain, hostname, ip)
	if err != nil {
		logger("client: "+client+" "+err.Error(), "err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte("OK\n"))
	logger("client: "+client+" successfully updated OVH DNS hostname: "+hostname+" ip: "+ip, "info")
}

func ovhSetup(applicationKey, applicationSecret, consumerKey, endpoint string) (*OVHClient, error) {
	if applicationKey == "" || applicationSecret == "" || consumerKey == "" {
		return nil, errors.New("OVH credentials incomplete - need application key, secret, and consumer key")
	}

	// Map endpoint names to full URLs
	var apiEndpoint string
	switch strings.ToLower(endpoint) {
	case "eu", "europe":
		apiEndpoint = "https://eu.api.ovh.com/1.0"
	case "ca", "canada":
		apiEndpoint = "https://ca.api.ovh.com/1.0"
	case "us", "usa", "united-states":
		apiEndpoint = "https://api.ovh.com/1.0"
	case "au", "australia":
		apiEndpoint = "https://au.api.ovh.com/1.0"
	default:
		return nil, errors.New("unsupported OVH endpoint - supported: eu, ca, us, au")
	}

	client := &OVHClient{
		Endpoint:          apiEndpoint,
		ApplicationKey:    applicationKey,
		ApplicationSecret: applicationSecret,
		ConsumerKey:       consumerKey,
		httpClient:        &http.Client{Timeout: 30 * time.Second},
	}

	return client, nil
}

func ovhUpdateDNS(client *OVHClient, domain, hostname, ip string) error {
	// Extract subdomain from hostname
	subdomain := ""
	if hostname != domain {
		if strings.HasSuffix(hostname, "."+domain) {
			subdomain = strings.TrimSuffix(hostname, "."+domain)
		} else {
			return errors.New("hostname does not match domain")
		}
	}

	// List existing A records for this subdomain
	records, err := client.listDNSRecords(domain, subdomain, "A")
	if err != nil {
		return errors.New("failed to list DNS records: " + err.Error())
	}

	// If record exists, update it; otherwise create new one
	if len(records) > 0 {
		// Update existing record
		recordID := records[0]
		err = client.updateDNSRecord(domain, recordID, ip)
		if err != nil {
			return errors.New("failed to update DNS record: " + err.Error())
		}
	} else {
		// Create new record
		record := OVHDNSRecord{
			SubDomain: subdomain,
			FieldType: "A",
			Target:    ip,
			TTL:       300,
		}
		_, err = client.createDNSRecord(domain, record)
		if err != nil {
			return errors.New("failed to create DNS record: " + err.Error())
		}
	}

	// Refresh the zone to apply changes
	err = client.refreshZone(domain)
	if err != nil {
		return errors.New("failed to refresh DNS zone: " + err.Error())
	}

	return nil
}

// OVH API helper methods
func (c *OVHClient) listDNSRecords(domain, subdomain, fieldType string) ([]int64, error) {
	path := fmt.Sprintf("/domain/zone/%s/record", domain)

	// Build query parameters
	params := make(map[string]string)
	if subdomain != "" {
		params["subDomain"] = subdomain
	}
	if fieldType != "" {
		params["fieldType"] = fieldType
	}

	body, err := c.makeRequest("GET", path, params, nil)
	if err != nil {
		return nil, err
	}

	var records []int64
	err = json.Unmarshal(body, &records)
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (c *OVHClient) createDNSRecord(domain string, record OVHDNSRecord) (int64, error) {
	path := fmt.Sprintf("/domain/zone/%s/record", domain)

	jsonData, err := json.Marshal(record)
	if err != nil {
		return 0, err
	}

	body, err := c.makeRequest("POST", path, nil, jsonData)
	if err != nil {
		return 0, err
	}

	var result struct {
		ID int64 `json:"id"`
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, err
	}

	return result.ID, nil
}

func (c *OVHClient) updateDNSRecord(domain string, recordID int64, target string) error {
	path := fmt.Sprintf("/domain/zone/%s/record/%d", domain, recordID)

	updateData := map[string]interface{}{
		"target": target,
		"ttl":    300,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		return err
	}

	_, err = c.makeRequest("PUT", path, nil, jsonData)
	return err
}

func (c *OVHClient) refreshZone(domain string) error {
	path := fmt.Sprintf("/domain/zone/%s/refresh", domain)
	_, err := c.makeRequest("POST", path, nil, nil)
	return err
}

func (c *OVHClient) makeRequest(method, path string, params map[string]string, body []byte) ([]byte, error) {
	// Build URL with query parameters
	url := c.Endpoint + path
	if len(params) > 0 {
		url += "?"
		for k, v := range params {
			url += k + "=" + v + "&"
		}
		url = strings.TrimSuffix(url, "&")
	}

	// Create request
	var req *http.Request
	var err error
	if body != nil {
		req, err = http.NewRequest(method, url, bytes.NewBuffer(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return nil, err
		}
	}

	// Add OVH API authentication headers
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)

	// Create signature
	bodyStr := ""
	if body != nil {
		bodyStr = string(body)
	}

	signature := c.createSignature(method, url, bodyStr, timestamp)

	req.Header.Set("X-Ovh-Application", c.ApplicationKey)
	req.Header.Set("X-Ovh-Consumer", c.ConsumerKey)
	req.Header.Set("X-Ovh-Timestamp", timestamp)
	req.Header.Set("X-Ovh-Signature", signature)

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, errors.New(fmt.Sprintf("OVH API error %d: %s", resp.StatusCode, string(respBody)))
	}

	return respBody, nil
}

func (c *OVHClient) createSignature(method, url, body, timestamp string) string {
	// OVH signature format: $1$<sha1_hex>
	toSign := c.ApplicationSecret + "+" + c.ConsumerKey + "+" + method + "+" + url + "+" + body + "+" + timestamp
	hash := sha1.Sum([]byte(toSign))
	return "$1$" + fmt.Sprintf("%x", hash)
}
