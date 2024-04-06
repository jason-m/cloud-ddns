package main

import (
	"net"
	"net/http"
	"os"
	"strconv"
)

var ip net.IP
var port int

func main() {
	// set default values to for ip/port to bind too, change this if you really know what you're doing
	// other wise run a reverse proxy with ssl to provide external access to this questionable app
	ip = net.ParseIP("127.0.0.1")
	port = 8080
	parseArgs()
	connectionString := ip.String() + ":" + strconv.Itoa(port)
	http.HandleFunc("/aws/", awsBasicAuth(awsHandler))
	http.ListenAndServe(connectionString, nil)
}

func parseArgs() {
	// Check if command line arguments are provided and if they are in the right format
	// this app only accepts 1 arg ip.ad.dr.es port ie ./cloud-ddns 10.0.0.1 8080
	if len(os.Args) == 3 {
		if net.ParseIP(os.Args[1]) != nil {
			ip = net.ParseIP(os.Args[1])
			portNum, err := strconv.Atoi(os.Args[2])
			if err == nil && portNum > 1 && portNum < 65534 {
				port = portNum
			} else {
				panic("invalid port specified")
			}
		} else {
			panic("invalid ip specified")
		}
	} else if len(os.Args) == 1 {
		// fmt.Print(len(os.Args))
		// nothign to do here really
	} else {
		// fmt.Print(len(os.Args))
		panic("too many or too few arguments this accepts either no arguments or 2 arguments in the format of:\n cloud-dns ipaddress port")
	}
}
