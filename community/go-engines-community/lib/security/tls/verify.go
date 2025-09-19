package tls

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"time"
)

func VerifySelfSignedCertificate(cfg *tls.Config) func([][]byte, [][]*x509.Certificate) error {
	return func(rawCerts [][]byte, _ [][]*x509.Certificate) error {
		isSelfSigned := true
		certs := make([]*x509.Certificate, len(rawCerts))
		for i, asn1Data := range rawCerts {
			cert, err := x509.ParseCertificate(asn1Data)
			if err != nil {
				return fmt.Errorf("tls: failed to parse certificate from server: %w", err)
			}

			certs[i] = cert
			if !bytes.Equal(cert.RawIssuer, cert.RawSubject) {
				isSelfSigned = false
			}
		}

		if isSelfSigned {
			return nil
		}

		opts := x509.VerifyOptions{
			Roots:         cfg.RootCAs,
			CurrentTime:   time.Now(),
			DNSName:       cfg.ServerName,
			Intermediates: x509.NewCertPool(),
		}

		for _, cert := range certs[1:] {
			opts.Intermediates.AddCert(cert)
		}

		_, err := certs[0].Verify(opts)
		if err != nil {
			return &tls.CertificateVerificationError{UnverifiedCertificates: certs, Err: err}
		}

		return nil
	}
}
