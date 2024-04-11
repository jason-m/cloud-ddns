package main

import (
	"context"
	"errors"
	"net/http"

	"github.com/cloudflare/cloudflare-go"
)

func cfHandler(w http.ResponseWriter, r *http.Request) {
	// first check for required form entries to satisfy ddns standard
	// then check for aws specific values (ie /aws/ZONEID)
	// fmt.Println(r.RemoteAddr)
	ip, hostname, err := checkForms(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger("client: "+r.RemoteAddr+" "+err.Error(), "err")
	}
	// fmt.Println(ip, hostname)

	// CF Is a much simpler api/package so it will all e done in this one step
	err = cfDoUpdate(user, pass, hostname, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger("client: "+r.RemoteAddr+" "+err.Error(), "err")
	} else {
		w.WriteHeader(400)
		w.Write([]byte("OK\n"))
		logger("client: "+r.RemoteAddr+" succesfully updated Cloudflare DNS hostname: "+hostname+" ip: "+ip, "info")
	}

}

func cfDoUpdate(zoneName, apiToken, hostname, ip string) error {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return errors.New("failed to create cloudfare api session")
	}
	zoneId, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return errors.New("zoneName not found (should be provided as username in cloudflare mode)")
	}

	// Check if DNS Record exists
	records, _, err := api.ListDNSRecords(context.Background(), cloudflare.ZoneIdentifier(zoneId), cloudflare.ListDNSRecordsParams{Name: hostname})
	if err != nil {
		// there was an error checking for the record
		return err
	} else {
		// DNS Record was found probably lets just double check the value is populated
		if len(records) == 0 {
			// record not actually found
			// Ok so there is no record found, lets then create a new DNS entry
			record := cloudflare.CreateDNSRecordParams{
				Type:    "A",
				Name:    hostname,
				Content: ip,
			}
			_, err = api.CreateDNSRecord(context.TODO(), cloudflare.ZoneIdentifier(zoneId), record)
			if err != nil {
				return err
			}
		} else {
			recordId := records[0].ID
			record := cloudflare.UpdateDNSRecordParams{
				Type:    "A",
				Name:    hostname,
				Content: ip,
				ID:      recordId,
			}
			_, err = api.UpdateDNSRecord(context.TODO(), cloudflare.ZoneIdentifier(zoneId), record)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
