package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func awsHandler(w http.ResponseWriter, r *http.Request) {
	// first check for required form entries to satisfy ddns standard
	// then check for aws specific values (ie /aws/ZONEID)
	// fmt.Println(r.RemoteAddr)
	ip, hostname, err := checkForms(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger("client: "+r.RemoteAddr+" "+err.Error(), "err")
	}

	// gets 3rd entry since url should be hostname/aws/zoneid
	// makes sure the zoneid is included in the url
	getZoneid := r.URL.Path
	// 4th entry should be the ?ip= blah blah
	// 2nd is the /aws/ so
	if len(strings.Split(getZoneid, "/")) != 4 {
		http.Error(w, "zoneid not detected", http.StatusBadRequest)
		logger("client: "+r.RemoteAddr+" zoneid not detected", "err")
		return
	} else {
		getZoneid = strings.Split(getZoneid, "/")[2]
	}

	// Setup AWS Session
	awsSession, err := awsSetup(user, pass)

	if awsSession != nil {
		// if session is created then updated dns
		err = awsRoute53(awsSession, getZoneid, hostname, ip)
	}
	if err != nil {
		logger("client:"+r.RemoteAddr+" "+err.Error(), "err")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		w.WriteHeader(400)
		w.Write([]byte("OK\n"))
		logger("client: "+r.RemoteAddr+" succesfully updated AWS DNS hostname: "+hostname+" ip: "+ip, "info")
	}

}

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
