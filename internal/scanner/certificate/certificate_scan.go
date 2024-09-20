package certificate

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"strings"
)

func GetCertificateWebsite(url string, port int) ([]*x509.Certificate, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	url = strings.TrimSuffix(url, "/")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	conn, err := tls.Dial("tcp", url+":"+fmt.Sprint(port), conf)
	if err != nil {
		log.Println("Error in Dial", err)
		return nil, err
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates

	return certs, nil
}
