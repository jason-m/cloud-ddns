package main

import (
	"fmt"
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
	fmt.Fprintf(w, "Hello, you've entered %s!\n", r.URL.Path[1:])
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}
	for key, values := range r.Form { // range over map
		for _, value := range values { // range over []string
			fmt.Printf("Form data: %s = %s\n", key, value)
		}
	}
}
