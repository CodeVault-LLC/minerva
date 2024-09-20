package security

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

var tlsVersions = map[uint16]string{
	tls.VersionSSL30: "SSL 3.0 (Insecure)",
	tls.VersionTLS10: "TLS 1.0 (Insecure)",
	tls.VersionTLS11: "TLS 1.1 (Insecure)",
	tls.VersionTLS12: "TLS 1.2 (Secure)",
	tls.VersionTLS13: "TLS 1.3 (Most Secure)",
}

// ScanProtocolSupport attempts to establish a TLS connection with a specific TLS version
func ScanProtocolSupport(addr string) ([]string, error) {
	var tlsVersionsChecked = []string{}

	for tlsVersion, versionName := range tlsVersions {
		tlsConfig := &tls.Config{
			MinVersion: tlsVersion,
			MaxVersion: tlsVersion,
		}

		// Set a timeout for the connection
		dialer := &net.Dialer{
			Timeout: 5 * time.Second,
		}

		// Establish a TCP connection
		conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
		if err != nil {
			fmt.Printf("Failed to connect using %s: %v\n", versionName, err)
			continue
		}
		defer conn.Close()

		// Print the successful connection details
		state := conn.ConnectionState()
		fmt.Printf("Connected using %s\n", tlsVersions[state.Version])

		tlsVersionsChecked = append(tlsVersionsChecked, tlsVersions[state.Version])
	}

	if len(tlsVersionsChecked) == 0 {
		return nil, fmt.Errorf("no supported TLS versions found for %s", addr)
	}

	return tlsVersionsChecked, nil
}
