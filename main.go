package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Inline self-signed certificate
	const selfSignedCert = `
-----BEGIN CERTIFICATE-----
<YOUR_SELF_SIGNED_CERTIFICATE>
-----END CERTIFICATE-----
`

	// Create a certificate pool and add the self-signed certificate
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM([]byte(selfSignedCert)) {
		log.Fatalf("Failed to append self-signed certificate")
	}

	// Create a custom transport with the self-signed certificate
	tlsConfig := &tls.Config{
		RootCAs: certPool,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: transport,
	}

	// URL of the HTTPS service
	url := "https://example.com"

	// Send a GET request
	resp, err := client.Get(url)
	if err != nil {
		log.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check and print the status code
	fmt.Printf("Response status: %s\n", resp.Status)
}
