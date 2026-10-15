package tls_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"testing"
	"time"

	libsectls "git.canopsis.net/canopsis/canopsis-community/community/go-engines-community/lib/security/tls"
	"github.com/rs/zerolog"
)

func TestCreateTLSConfigFromEnv(t *testing.T) {
	const envVarCaCertFile = "TEST_TLS_CA_CERT_FILE"
	const envVarInsecureSkipVerify = "TEST_TLS_INSECURE_SKIP_VERIFY"

	validCertPath := writeValidCertFile(t)
	invalidPemPath := writeInvalidCertFile(t)
	validCertPool := certPoolFromFile(t, validCertPath)
	emptyCertPool := x509.NewCertPool()

	tests := []struct {
		name      string
		caValue   string
		skipValue string
		assert    func(*testing.T, *tls.Config)
	}{
		{
			name: "returns nil when both env vars are empty",
			assert: func(t *testing.T, got *tls.Config) {
				t.Helper()
				if got != nil {
					t.Fatalf("expected nil tls config, got %#v", got)
				}
			},
		},
		{
			name:      "returns insecure config when insecure flag is true",
			skipValue: "TrUe",
			assert: func(t *testing.T, got *tls.Config) {
				t.Helper()
				if got == nil {
					t.Fatalf("expected tls config, got nil")
				}
				if !got.InsecureSkipVerify {
					t.Fatalf("expected InsecureSkipVerify=true")
				}
				if got.VerifyPeerCertificate == nil {
					t.Fatalf("expected VerifyPeerCertificate to be set")
				}
				if got.VerifyConnection == nil {
					t.Fatalf("expected VerifyConnection to be set")
				}
				if got.RootCAs != nil {
					t.Fatalf("expected RootCAs to be nil when CA file env is empty")
				}
			},
		},
		{
			name:      "returns nil when insecure flag is not true",
			skipValue: "1",
			assert: func(t *testing.T, got *tls.Config) {
				t.Helper()
				if got != nil {
					t.Fatalf("expected nil tls config, got %#v", got)
				}
			},
		},
		{
			name:    "returns config with nil RootCAs when CA file cannot be read",
			caValue: "/tmp/file-does-not-exist.pem",
			assert: func(t *testing.T, got *tls.Config) {
				t.Helper()
				if got == nil {
					t.Fatalf("expected tls config, got nil")
				}
				if got.RootCAs != nil {
					t.Fatalf("expected RootCAs=nil when cert file is unreadable")
				}
				if got.InsecureSkipVerify {
					t.Fatalf("expected InsecureSkipVerify=false when CA file env is set")
				}
			},
		},
		{
			name:    "returns config with empty RootCAs when CA file is invalid PEM",
			caValue: invalidPemPath,
			assert: func(t *testing.T, got *tls.Config) {
				t.Helper()
				if got == nil {
					t.Fatalf("expected tls config, got nil")
				}
				if got.RootCAs == nil {
					t.Fatalf("expected RootCAs to be initialized")
				}
				if !got.RootCAs.Equal(emptyCertPool) {
					t.Fatalf("expected RootCAs to equal an empty cert pool")
				}
			},
		},
		{
			name:      "CA file takes precedence over insecure flag",
			caValue:   validCertPath,
			skipValue: "true",
			assert: func(t *testing.T, got *tls.Config) {
				t.Helper()
				if got == nil {
					t.Fatalf("expected tls config, got nil")
				}
				if got.RootCAs == nil {
					t.Fatalf("expected RootCAs to be initialized")
				}
				if !got.RootCAs.Equal(validCertPool) {
					t.Fatalf("expected RootCAs to match the cert pool loaded from CA file")
				}
				if got.InsecureSkipVerify {
					t.Fatalf("expected InsecureSkipVerify=false when CA file env is set")
				}
				if got.VerifyPeerCertificate != nil {
					t.Fatalf("expected VerifyPeerCertificate=nil when CA file path is used")
				}
				if got.VerifyConnection != nil {
					t.Fatalf("expected VerifyConnection=nil when CA file path is used")
				}
			},
		},
	}

	logger := zerolog.Nop()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(envVarCaCertFile, tt.caValue)
			t.Setenv(envVarInsecureSkipVerify, tt.skipValue)

			got := libsectls.CreateTLSConfigFromEnv(envVarCaCertFile, envVarInsecureSkipVerify, logger)
			tt.assert(t, got)
		})
	}
}

func certPoolFromFile(t *testing.T, filePath string) *x509.CertPool {
	t.Helper()

	caCert, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed to read certificate file %q: %v", filePath, err)
	}

	pool := x509.NewCertPool()
	if !pool.AppendCertsFromPEM(caCert) {
		t.Fatalf("failed to append certificate from file %q", filePath)
	}

	return pool
}

func writeValidCertFile(t *testing.T) string {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate private key: %v", err)
	}

	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62))
	if err != nil {
		t.Fatalf("failed to generate serial number: %v", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName: "test-ca",
		},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(time.Hour),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("failed to create certificate: %v", err)
	}

	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if pemBytes == nil {
		t.Fatalf("failed to encode certificate to PEM")
	}

	tmpFile, err := os.CreateTemp(t.TempDir(), "ca-*.pem")
	if err != nil {
		t.Fatalf("failed to create temporary cert file: %v", err)
	}
	tmpPath := tmpFile.Name()
	t.Cleanup(func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
	})

	if _, err = tmpFile.Write(pemBytes); err != nil {
		t.Fatalf("failed to write certificate to temporary file: %v", err)
	}

	return tmpPath
}

func writeInvalidCertFile(t *testing.T) string {
	t.Helper()

	tmpFile, err := os.CreateTemp(t.TempDir(), "invalid-*.pem")
	if err != nil {
		t.Fatalf("failed to create temporary invalid cert file: %v", err)
	}
	tmpPath := tmpFile.Name()
	t.Cleanup(func() {
		_ = tmpFile.Close()
		_ = os.Remove(tmpPath)
	})

	if _, err = tmpFile.WriteString("this-is-not-a-pem"); err != nil {
		t.Fatalf("failed to write invalid PEM content: %v", err)
	}

	return tmpPath
}
