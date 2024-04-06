package main

import (
	"net"
	"os"
	"strconv"
)

//Globally define default host and port to listen on
//leaving this default is recommended for most use cases,
//and running it behind a reverse proxy with ssl

var ip net.IP
var port int

func main() {
	ip = net.ParseIP("127.0.0.1")
	port = 8080
	parseArgs()
}

func parseArgs() {
	// Check if command line arguments are provided and if they are in the right format
	// this app only accepts 1 arg ip.ad.dr.es port ie cloud-ddns 10.0.0.1 8080
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
