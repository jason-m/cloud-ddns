package main

// This file contains the generic DDNS related functions and http handlers for each cloud service
// for cloud service specific functions stored in there own CLOUDNAME.go files in the same directory

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
)

var user, pass string
var ok bool

func awsBasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok = r.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
			w.WriteHeader(401)
			w.Write([]byte("You are unauthorized to access the application.\n"))
			return
		}

		// fmt.Printf("Username: %s, Password: %s\n", user, pass)
		handler(w, r)
	}
}

func awsHandler(w http.ResponseWriter, r *http.Request) {
	// first check for required form entries to satisfy ddns standard
	// then check for aws specific values (ie /aws/ZONEID)
	ip, hostname, err := checkForms(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// gets 3rd entry since url should be hostname/aws/zoneid
	getZoneid := r.URL.Path
	getZoneid = strings.Split(getZoneid, "/")[2]

	// Setup AWS Session
	awsSession, err := awsSetup(user, pass)

	if awsSession != nil {
		fmt.Println("Sessionn setup")
		err = awsRoute53(awsSession, getZoneid, hostname, ip)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	} else {
		// fmt.Printf(" %v \n\n", awsSession)
	}

}

func checkForms(r *http.Request) (ip string, hostname string, err error) {
	// since dyndns proto always requires these 2 form values generic function for checking them
	err = r.ParseForm()
	if err != nil {
		err = errors.New("failed to parse form")
		return "", "", err
	}
	var ipCheck []string
	var check bool
	ipCheck, check = r.Form["ip"]
	if !check {
		err = errors.New("required form value \"ip\"")
		return "", "", err
	} else if net.ParseIP(ipCheck[0]) == nil {
		err = errors.New("ip address invalid")
		return "", "", err
	} else {
		ip = ipCheck[0]
	}
	var nameCheck []string
	nameCheck, check = r.Form["hostname"]
	if !check {
		err = errors.New("required form value \"hostname\"")
		return "", "", err
	} else {
		hostname = nameCheck[0]
	}
	return ip, hostname, nil
}
