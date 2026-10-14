package tls

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
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

func VerifyConnectionClientSide(roots *x509.CertPool) func(tls.ConnectionState) error {
	return func(cs tls.ConnectionState) error {
		opts := x509.VerifyOptions{
			Roots:         roots,
			DNSName:       cs.ServerName,
			Intermediates: x509.NewCertPool(),
		}
		for _, cert := range cs.PeerCertificates[1:] {
			opts.Intermediates.AddCert(cert)
		}
		_, err := cs.PeerCertificates[0].Verify(opts)
		if err != nil && roots == nil {
			// If the verification failed and no custom roots were provided, it may be because the certificate is self-signed.
			if _, ok := errors.AsType[x509.UnknownAuthorityError](err); ok {
				return nil
			}
		}
		return err
	}
}

// CreateTLSConfigFromEnv creates a TLS configuration based on environment variables.
// It allows for either using a custom CA certificate or skipping verification (not recommended for production).
func CreateTLSConfigFromEnv(envVarCaCertFile, envVarInsecureSkipVerify string, logger zerolog.Logger) *tls.Config {
	cf := os.Getenv(envVarCaCertFile)
	if cf != "" {
		tc := &tls.Config{}
		var roots *x509.CertPool
		caCert, err := os.ReadFile(cf)
		if err != nil {
			logger.Warn().Err(err).Msgf("failed to read certificate file %q, skip it", cf)
		} else {
			roots = x509.NewCertPool()
			if !roots.AppendCertsFromPEM(caCert) {
				logger.Warn().Msgf("failed to append certificate from file %q, skip it", cf)
			} else {
				logger.Debug().Msgf("certificate from file %q added to CA", cf)
			}
		}
		tc.RootCAs = roots
		return tc
	}

	hasInsecureSkipVerify := strings.ToLower(os.Getenv(envVarInsecureSkipVerify)) == "true"
	if hasInsecureSkipVerify {
		tc := &tls.Config{
			InsecureSkipVerify: true, // #nosec G402
		}
		logger.Warn().Msgf("environment variable %q empty, skip adding certificate file to CA", envVarCaCertFile)
		tc.VerifyPeerCertificate = VerifySelfSignedCertificate(tc)
		tc.VerifyConnection = VerifyConnectionClientSide(nil)
		return tc
	}
	return nil
}
