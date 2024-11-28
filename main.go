package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Embed the CA certificate directly into the code as a string
const caCertPEM = `
-----BEGIN CERTIFICATE-----
MIIDqTCCApGgAwIBAgIUJMkswrVQuxpk0mrvVPznZ9wEYNswDQYJKoZIhvcNAQEL
BQAwajELMAkGA1UEBhMCTEsxEDAOBgNVBAgMB2NvbG9tYm8xFjAUBgNVBAcMDVNh
biBGcmFuY2lzY28xDTALBgNVBAoMBHdzbzIxIjAgBgNVBAMMGXZhamlyYS0yLnRh
aWwwYmIyNC50cy5uZXQwHhcNMjQxMTI4MDgyNzAzWhcNMjUxMTI4MDgyNzAzWjBq
MQswCQYDVQQGEwJMSzEQMA4GA1UECAwHY29sb21ibzEWMBQGA1UEBwwNU2FuIEZy
YW5jaXNjbzENMAsGA1UECgwEd3NvMjEiMCAGA1UEAwwZdmFqaXJhLTIudGFpbDBi
YjI0LnRzLm5ldDCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALbmr3/y
TgeCQ9gNIY3T12UBWslhVMeHOsrU3XGS91hzFESfh8tCmDEhWFrfqMKOQim/5Cpx
hBlH2cqt7UK46/IhnNmUJeJYmHcbLswUwFrQJhuOqhiuT9B3H39EgZSZhctSY+uE
BUV2xvbOhbhh4qgni2e7Thr63e2BHag1e1Gr12yoByNbUCVFTFuWhmbFTrYscdYY
6OVsmr7LYRXYgQWEoSs6yTraUbEN2fN0M90StWPB29SzQ/iXERkMMMdf7mjouK2g
qhAJlouQdnSJAojcPDt6nm3/fY0lYzWozEICY9ShnMYIS1/1wFjJuQ7Ok7UZiwMx
vQ33FUdEWUWPT/0CAwEAAaNHMEUwJAYDVR0RBB0wG4IZdmFqaXJhLTIudGFpbDBi
YjI0LnRzLm5ldDAdBgNVHQ4EFgQU0/2+0mjLdZS+/tHq3Fi8b18M4ecwDQYJKoZI
hvcNAQELBQADggEBALZj5ttz5jtzTYTnyMhCenKRsrVCl9fOc4UYpaGx84IEfHtD
Wk76dxDfCtAr9oSNab1Z4L1pZdSORqXKzCSBCjahuG3bgwWmFXkyCK2k42sndQDq
Zz4P4wiolmOvxowv0Aw+HZGGBnvunqp9pFoC9Pe4GTyX/mvmuOzQweNHEbFzbLe+
z93biIBJmIceB6L4gqBoTQK4xmOSRMgxMHiYXxG2vZ5HuXa1FtF63i4l3HdZgiOH
HjR+fnnKLqCAdD8Uj2L9iDR2+9h3uS0WhQ4LrUvjdZ/kQqiqAl/qH6Nh6j8Xx8vk
8wZFN9OxBB3NLzQVBrMfU3nhN4KT/IwLYDwJInw=
-----END CERTIFICATE-----
`

func main() {
	// Start the proxy HTTP server
	http.HandleFunc("/", proxyHandler)

	port := ":8080"
	fmt.Printf("Proxy service is running on http://localhost%s/proxy\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}

// Proxy handler function
func proxyHandler(w http.ResponseWriter, r *http.Request) {
	// Create a CA certificate pool and add the embedded CA certificate
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM([]byte(caCertPEM)); !ok {
		http.Error(w, "Failed to append CA certificate to pool", http.StatusInternalServerError)
		return
	}

	// Create a TLS configuration with the CA certificate and SNI enabled
	tlsConfig := &tls.Config{
		RootCAs:    caCertPool,
		ServerName: "vajira-2.tail0bb24.ts.net", // SNI
	}

	// Use a custom Transport with the TLS configuration
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	// Create an HTTP client with the custom Transport
	client := &http.Client{
		Transport: transport,
	}

	// Make the request to the target server

	//http://ts-proxy-2-3414549712.dp-development-tailscaleproxyexpo-7190-1897193973.svc.cluster.local:8080
	url := "https://ts-byoc-1708985667.dp-development-testproject-8874-3254367504.svc.cluster.local:8080"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
		return
	}

	// Add the resolve header equivalent
	req.Host = "vajira-2.tail0bb24.ts.net" // Override DNS resolution

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to make the request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Copy the response back to the caller
	w.WriteHeader(resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read response body: %v", err), http.StatusInternalServerError)
		return
	}
	w.Write(body)
}
