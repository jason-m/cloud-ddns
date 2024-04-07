package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

func awsBasicAuth(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()

		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Please enter your username and password"`)
			w.WriteHeader(401)
			w.Write([]byte("You are unauthorized to access the application.\n"))
			return
		}

		fmt.Printf("Username: %s, Password: %s\n", user, pass)
		handler(w, r)
	}
}

func awsHandler(w http.ResponseWriter, r *http.Request) {
	// first check for required form entries to satisfy ddns standard
	// then check for aws specific values (ie /aws/ZONEID)
	err := checkForms(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		w.WriteHeader(400)
		w.Write([]byte("OK.\n"))
	}
}

func checkForms(r *http.Request) (err error) {
	// since dyndns proto always requires these 2 form values generic function for checking them
	err = r.ParseForm()
	if err != nil {
		err = errors.New("failed to parse form")
		return err
	}

	ip, check := r.Form["ip"]
	if !check {
		err = errors.New("required form value \"ip\"")
		return err
	} else if net.ParseIP(ip[0]) == nil {
		err = errors.New("ip address invalid")
		return err
	} else {
		return nil
	}
}
