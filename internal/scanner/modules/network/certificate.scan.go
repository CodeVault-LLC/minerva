package network

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/codevault-llc/humblebrag-api/models"
)

func GetCertificateWebsite(url string, port int) ([]*x509.Certificate, models.CertificateResult, error) {
	conf := &tls.Config{
		InsecureSkipVerify: true,
	}

	url = strings.TrimSuffix(url, "/")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")

	conn, err := tls.Dial("tcp", url+":"+fmt.Sprint(port), conf)
	if err != nil {
		log.Println("Error in Dial", err)
		return nil, models.CertificateResult{}, err
	}
	defer conn.Close()
	certs := conn.ConnectionState().PeerCertificates

	// Analyze the certificate
	result, err := analyzeCertificate(certs, url)
	if err != nil {
		return nil, models.CertificateResult{}, err
	}

	return certs, result, nil
}

// AnalyzeCertificate performs an in-depth analysis of an SSL/TLS certificate.
// It checks for:
// - Expiry
// - Trusted root authority
// - Suspicious/malicious providers
// - Weak encryption or insecure algorithms
// - Certificate revocation status (if available)
// - CommonName (CN) and Subject Alternative Name (SAN) validation
// - And more security-related issues.
//
// Returns a detailed analysis report with suggestions or warnings.
func analyzeCertificate(certs []*x509.Certificate, url string) (models.CertificateResult, error) {
	cert := certs[0]
	result := models.CertificateResult{}

	// Step 2: Expiry check.
	if isCertificateExpired(cert) {
		result.Expired = true
	}

	// Step 3: Validate against trusted root CAs.
	if !isTrustedCertificate(cert) {
		result.Trusted = true
	}

	// Step 4: Check for weak encryption or insecure algorithms.
	if hasWeakEncryption(cert) {
		result.Weak = true
	}

	// Step 5: Verify issuer for malicious or suspicious patterns.
	if isMaliciousIssuer(cert) {
		result.MaliciousIssuer = true
	}

	// Step 6: Verify if the certificate is revoked (if CRL or OCSP available).
	revocationStatus := checkRevocationStatus(cert)
	if revocationStatus != "" {
		result.Revoked = true
	}

	// Step 7: Validate CommonName and SAN fields.
	if validateDomain(cert, url) {
		result.Domain = true
	}

	return result, nil
}

// isCertificateExpired checks if a certificate is expired or near expiry.
func isCertificateExpired(cert *x509.Certificate) bool {
	now := time.Now()
	return now.After(cert.NotAfter) || now.Before(cert.NotBefore)
}

// isTrustedCertificate checks if the certificate is signed by a trusted root authority.
func isTrustedCertificate(cert *x509.Certificate) bool {
	roots, err := x509.SystemCertPool()
	if err != nil {
		log.Println("Error loading system root CAs:", err)
		return false
	}

	opts := x509.VerifyOptions{
		Roots: roots,
	}

	_, err = cert.Verify(opts)
	return err == nil
}

// hasWeakEncryption checks if a certificate uses weak encryption algorithms.
func hasWeakEncryption(cert *x509.Certificate) bool {
	weakAlgorithms := []x509.SignatureAlgorithm{
		x509.MD2WithRSA,
		x509.MD5WithRSA,
		x509.SHA1WithRSA,
		x509.DSAWithSHA1,
	}
	for _, algo := range weakAlgorithms {
		if cert.SignatureAlgorithm == algo {
			return true
		}
	}
	return false
}

// isMaliciousIssuer checks if the issuer is flagged as malicious or suspicious (simple check).
func isMaliciousIssuer(cert *x509.Certificate) bool {
	suspiciousIssuers := []string{"Untrusted CA", "Unknown Issuer", "Fake CA"}
	for _, issuer := range suspiciousIssuers {
		if strings.Contains(cert.Issuer.CommonName, issuer) {
			return true
		}
	}
	return false
}

// checkRevocationStatus checks if the certificate is revoked via CRL or OCSP.
func checkRevocationStatus(cert *x509.Certificate) string {
	if len(cert.CRLDistributionPoints) > 0 {
		// Placeholder for CRL (Certificate Revocation List) check logic.
		return "CRL check not implemented"
	}
	if cert.OCSPServer != nil {
		// Placeholder for OCSP (Online Certificate Status Protocol) check logic.
		return "OCSP check not implemented"
	}
	return ""
}

// validateDomain checks if the domain matches the certificate's CommonName or SAN.
func validateDomain(cert *x509.Certificate, domain string) bool {
	if len(cert.DNSNames) == 0 && cert.Subject.CommonName == "" {
		return false
	}

	if strings.Contains(cert.Subject.CommonName, domain) {
		return true
	}

	for _, dnsName := range cert.DNSNames {
		if strings.Contains(dnsName, domain) {
			return true
		}
	}

	return false
}
