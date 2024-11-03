package modules

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"strings"

	"github.com/codevault-llc/humblebrag-api/internal/models/entities"
)

type CertificateModule struct{}

func (m *CertificateModule) Run(job entities.JobModel) (interface{}, error) {
	certifiate, err := GetCertificateWebsite(job.URL, 443)
	if err != nil {
		return nil, err
	} else {
		return certifiate, nil
	}
}

func (m *CertificateModule) Name() string {
	return "Certificate"
}

// GetCertificateWebsite returns the certificate of a website
func GetCertificateWebsite(url string, port int) ([]*x509.Certificate, error) {
	conf := &tls.Config{
		// file deepcode ignore TooPermissiveTrustManager: a scanning module to verify third party certificates
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
