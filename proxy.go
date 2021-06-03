package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var host string
var port int
var insecure bool

func init() {
	flag.StringVar(&host, "host", "", "Host to start the proxy (default all interfaces)")
	flag.IntVar(&port, "port", 8080, "Port to start the proxy")
	flag.BoolVar(&insecure, "insecure", false, "Allow insecure server connections when using SSL")
}

func usage() {
	log.Print("Usage: proxy [options] http://example.org\n\nOptions:\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	target := flag.Arg(0)
	targetUri, err := url.ParseRequestURI(target)

	if err != nil {
		log.Print("Error parsing target URI")
		log.Print(err)
		usage()
		os.Exit(1)
	}

	address := fmt.Sprintf("%s:%d", host, port)

	log.Printf("Proxying requests from %s to %s", address, target)
	log.Fatal(http.ListenAndServe(address, &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = targetUri
			r.Host = targetUri.Host
		},
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
		},
	}))
}
