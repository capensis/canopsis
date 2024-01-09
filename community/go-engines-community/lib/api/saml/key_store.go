package saml

import (
	"crypto/rsa"

	dsig "github.com/russellhaering/goxmldsig"
)

func NewCanopsisX509KeyStore(key *rsa.PrivateKey, cert []byte) dsig.X509KeyStore {
	return &canopsisX509KeyStore{
		key:  key,
		cert: cert,
	}
}

type canopsisX509KeyStore struct {
	key  *rsa.PrivateKey
	cert []byte
}

func (ks *canopsisX509KeyStore) GetKeyPair() (privateKey *rsa.PrivateKey, cert []byte, err error) {
	return ks.key, ks.cert, nil
}
