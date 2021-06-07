package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

var host string
var port int

var insecure bool
var caFile, certFile, keyFile string

func lookupEnvOrString(key string, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func init() {
	flag.StringVar(&host, "host", lookupEnvOrString("HOST", ""), "Host to start the proxy (default all interfaces)")
	flag.IntVar(&port, "port", 8080, "Port to start the proxy")

	flag.BoolVar(&insecure, "insecure", false, "Allow insecure server connections when using SSL")
	flag.StringVar(&caFile, "cacert", "", "CA certificate to verify peer against")
	flag.StringVar(&certFile, "cert", "", "Client certificate file")
	flag.StringVar(&keyFile, "key", "", "Private key file name")
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

	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecure,
	}

	if caFile != "" {
		caCert, err := ioutil.ReadFile(caFile)
		if err != nil {
			log.Print("Error loading certificate and/or private key file")
			log.Print(err)
			os.Exit(1)
		}

		tlsConfig.RootCAs = x509.NewCertPool()
		tlsConfig.RootCAs.AppendCertsFromPEM(caCert)
	}

	if certFile != "" || keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Print("Error loading certificate and/or private key file")
			log.Print(err)
			os.Exit(1)
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	address := fmt.Sprintf("%s:%d", host, port)

	log.Printf("Proxying requests from %s to %s", address, target)
	log.Fatal(http.ListenAndServe(address, &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL = targetUri
			r.Host = targetUri.Host
		},
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}))
}
