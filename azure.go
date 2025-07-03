package main

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dns/armdns"
)

func azureHandler(w http.ResponseWriter, r *http.Request) {

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

	// gets path components since url should be hostname/azure/tenantid/subscriptionid/resource-group/zone-name
	// makes sure the tenantid, subscriptionid, resource group, and zone name are included in the url
	pathComponents := strings.Split(r.URL.Path, "/")
	// Expected: ["", "azure", "tenantid", "subscriptionid", "resource-group", "zone-name", ""]
	if len(pathComponents) != 7 {
		http.Error(w, "invalid path format - expected /azure/tenantid/subscriptionid/resource-group/zone-name/", http.StatusBadRequest)
		logger("client: "+client+" invalid path format", "err")
		return
	}

	getTenantId := pathComponents[2]
	getSubscriptionId := pathComponents[3]
	getResourceGroup := pathComponents[4]
	getZoneName := pathComponents[5]

	if getTenantId == "" || getSubscriptionId == "" || getResourceGroup == "" || getZoneName == "" {
		http.Error(w, "tenantid, subscriptionid, resource-group, and zone-name are required", http.StatusBadRequest)
		logger("client: "+client+" missing required path components", "err")
		return
	}

	// Setup Azure DNS Client
	azureClient, err := azureSetup(user, pass, getTenantId, getSubscriptionId)

	if azureClient != nil {
		// if client is created then update dns
		err = azureDNS(azureClient, getResourceGroup, getZoneName, hostname, ip)
	}
	if err != nil {
		logger("client:"+client+" "+err.Error(), "err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
		logger("client: "+client+" succesfully updated Azure DNS hostname: "+hostname+" ip: "+ip, "info")
	}

}

func azureSetup(clientId, clientSecret, tenantId, subscriptionId string) (*armdns.RecordSetsClient, error) {
	cred, err := azidentity.NewClientSecretCredential(tenantId, clientId, clientSecret, nil)
	if err != nil {
		return nil, errors.New("failed to create azure credentials")
	}

	client, err := armdns.NewRecordSetsClient(subscriptionId, cred, nil)
	if err != nil {
		return nil, errors.New("failed to create azure dns client")
	}

	return client, nil
}

func azureDNS(client *armdns.RecordSetsClient, resourceGroupName string, zoneName string, hostname string, ip string) error {
	ctx := context.Background()

	// Extract record name from hostname and zone name
	// For example: host.example.com with zone example.com -> host
	// If hostname equals zone name, it's the root record (@)
	var recordName string
	if hostname == zoneName {
		recordName = "@"
	} else if strings.HasSuffix(hostname, "."+zoneName) {
		recordName = strings.TrimSuffix(hostname, "."+zoneName)
	} else {
		return errors.New("hostname " + hostname + " does not belong to zone " + zoneName)
	}

	// Create the A record data
	aRecords := []*armdns.ARecord{
		{
			IPv4Address: to.Ptr(ip),
		},
	}

	recordSetParams := armdns.RecordSet{
		Properties: &armdns.RecordSetProperties{
			TTL:      to.Ptr[int64](300),
			ARecords: aRecords,
		},
	}

	// Try to create or update the record (UPSERT operation)
	_, err := client.CreateOrUpdate(ctx, resourceGroupName, zoneName, recordName, armdns.RecordTypeA, recordSetParams, nil)
	if err != nil {
		return errors.New("failed to create/update dns record: " + err.Error())
	}

	return nil
}
