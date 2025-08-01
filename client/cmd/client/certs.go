package main

import (
	"crypto/tls"
	"crypto/x509"
	_ "embed"
	"log"
	"net/http"

	"github.com/Le-BlitzZz/blockchain-auth-app/client/internal/session"
)

//go:embed certs/ca.pem
var cliCAPEM []byte

//go:embed certs/client.crt
var cliCertPEM []byte

//go:embed certs/client.key
var cliKeyPEM []byte

func setup() {
	cert, err := tls.X509KeyPair(cliCertPEM, cliKeyPEM)
	if err != nil {
		log.Fatalf("failed to parse client certificate/key: %v", err)
	}
	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(cliCAPEM) {
		log.Fatal("failed to append CA certificate")
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	session.SetHTTPClient(&http.Client{
		Transport: &http.Transport{TLSClientConfig: tlsConfig},
	})
}
