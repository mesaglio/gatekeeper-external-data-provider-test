package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"mesaglio/gatekeeper-external-data-provider-test/pkg/handler"
	"mesaglio/gatekeeper-external-data-provider-test/pkg/utils"
)

const (
	timeout     = 5 * time.Second
	defaultPort = 8090

	certName = "tls.crt"
	keyName  = "tls.key"
)

var (
	certDir      string
	clientCAFile string
	port         int
)

func init() {
	flag.StringVar(&certDir, "cert-dir", "", "path to directory containing TLS certificates")
	flag.StringVar(&clientCAFile, "client-ca-file", "", "path to client CA certificate")
	flag.IntVar(&port, "port", defaultPort, "Port for the server to listen on")
	flag.Parse()
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", processTimeout(handler.Handler, timeout))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(5) * time.Second,
	}

	config := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
	if clientCAFile != "" {
		fmt.Printf("loading Gatekeeper's CA certificate: %s\n", clientCAFile)
		caCert, err := os.ReadFile(clientCAFile)
		if err != nil {
			fmt.Printf("ERROR: unable to load Gatekeeper's CA certificate: %s\nError: %s\n", clientCAFile, err.Error())
			os.Exit(1)
		}

		clientCAs := x509.NewCertPool()
		clientCAs.AppendCertsFromPEM(caCert)

		config.ClientCAs = clientCAs
		config.ClientAuth = tls.RequireAndVerifyClientCert
		server.TLSConfig = config
	}

	if certDir != "" {
		certFile := filepath.Join(certDir, certName)
		keyFile := filepath.Join(certDir, keyName)

		fmt.Printf("starting external data provider server on port: %d, certFile: %s, keyFile: %s\n", port, certFile, keyFile)
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			fmt.Printf("unable to start external data provider server: %s\n", err.Error())
			os.Exit(1)
		}
	} else {
		fmt.Println("TLS certificates are not provided, the server will not be started")
		os.Exit(1)
	}
}

func processTimeout(h http.HandlerFunc, duration time.Duration) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), duration)
		defer cancel()

		r = r.WithContext(ctx)

		processDone := make(chan bool)
		go func() {
			h(w, r)
			processDone <- true
		}()

		select {
		case <-ctx.Done():
			utils.SendResponse(nil, "operation timed out", w)
		case <-processDone:
		}
	}
}
