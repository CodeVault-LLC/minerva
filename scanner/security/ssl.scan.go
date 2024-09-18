package security

import (
	"crypto/tls"
	"fmt"
	"net"
	"time"
)

// TLS versions with human-readable names
var tlsVersions = map[uint16]string{
	tls.VersionSSL30: "SSLv3 (Insecure and Deprecated)",
	tls.VersionTLS10: "TLS 1.0 (Insecure)",
	tls.VersionTLS11: "TLS 1.1 (Deprecated)",
	tls.VersionTLS12: "TLS 1.2 (Secure)",
	tls.VersionTLS13: "TLS 1.3 (Most Secure)",
}

// CheckProtocolSupport attempts to establish a TLS connection with a specific TLS version
func ScanProtocolSupport(host string, tlsVersion uint16) error {
	// Create a custom TLS configuration to enforce the specific TLS version
	tlsConfig := &tls.Config{
		MinVersion: tlsVersion,
		MaxVersion: tlsVersion,
	}

	// Set a timeout for the connection
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}

	// Establish a TCP connection
	conn, err := tls.DialWithDialer(dialer, "tcp", host, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// Print the successful connection details
	state := conn.ConnectionState()
	fmt.Printf("Connected using %s\n", tlsVersions[state.Version])
	return nil
}
