package main

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func awsSetup(accessKey, secretKey string) (*session.Session, error) {
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String("global"),
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
	})
	if err != nil {
		return nil, err
	}
	return session, nil
}

func awsRoute53(session *session.Session, zoneid string, hostname string, ip string) error {
	r53 := route53.New(session)
	// query list of zones
	zones, err := r53.ListHostedZones(nil)
	if err != nil {
		return err
	}
	// ensure zoneid is exists
	var foundZone bool
	for z := range zones.HostedZones {
		foundZone = false
		if strings.Contains(*zones.HostedZones[z].Id, zoneid) {
			if strings.Contains(hostname+".", *zones.HostedZones[z].Name) {
				// fmt.Println("Hostname matches zone")
				foundZone = true
				break
			} else {
				return errors.New("hostname does not match zone")
			}

		}
	}
	if !foundZone {
		return errors.New("zone not found")
	}
	// time for the ugly aws route53 update

	_, err = r53.ChangeResourceRecordSets(&route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(hostname),
						Type: aws.String("A"),
						TTL:  aws.Int64(300),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String(ip),
							},
						},
					},
				},
			},
		},
		HostedZoneId: aws.String(zoneid),
	})
	if err != nil {
		return err
	}

	return nil
}
