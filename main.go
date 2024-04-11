package main

import (
	"log/syslog"
	"net"
	"net/http"
	"os"
	"strconv"
)

var listenIP net.IP
var port int

const appName = "cloud-ddns"

func main() {
	// set default values to for ip/port to bind too, change this if you really know what you're doing
	// other wise run a reverse proxy with ssl to provide external access to this questionable app
	listenIP = net.ParseIP("127.0.0.1")
	port = 8080
	parseArgs()
	connectionString := listenIP.String() + ":" + strconv.Itoa(port)
	http.HandleFunc("/aws/", BasicAuth(awsHandler))
	http.HandleFunc("/cloudfare/", BasicAuth(cfHandler))
	http.ListenAndServe(connectionString, nil)
}

func parseArgs() {
	// Check if command line arguments are provided and if they are in the right format
	// this app only accepts 1 arg ip.ad.dr.es port ie ./cloud-ddns 10.0.0.1 8080
	if len(os.Args) == 3 {
		if net.ParseIP(os.Args[1]) != nil {
			listenIP = net.ParseIP(os.Args[1])
			portNum, err := strconv.Atoi(os.Args[2])
			if err == nil && portNum > 1 && portNum < 65534 {
				port = portNum
				logger("application listening on "+listenIP.String()+":"+strconv.Itoa(port), "info")
			} else {
				logger("failed to start invalid port specified", "err")
				panic("invalid port specified")
			}
		} else {
			logger("failed to start invalid ip specified", "err")
			panic("invalid ip specified")
		}
	} else if len(os.Args) == 1 {
		// fmt.Print(len(os.Args))
		// nothign to do here really
	} else {
		// fmt.Print(len(os.Args))
		logger("failed to start invalid command line args", "err")
		panic("too many or too few arguments this accepts either no arguments or 2 arguments in the format of:\n cloud-dns ipaddress port")
	}
}

func logger(message, loglevel string) {
	switch loglevel {
	case "err":
		errLog, _ := syslog.New(syslog.LOG_ERR, appName)
		errLog.Err(message)
	case "notice":
		noticeLog, _ := syslog.New(syslog.LOG_NOTICE, appName)
		noticeLog.Notice(message)
	case "info":
		infoLog, _ := syslog.New(syslog.LOG_INFO, appName)
		infoLog.Info(message)
	}
}
