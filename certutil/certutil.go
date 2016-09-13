package certutil

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"time"
)

// Store holds an x509 certificate, its corresponding private key, and all certs in its trust chain.
type Store struct {
	key   interface{}
	cert  *x509.Certificate
	chain []*x509.Certificate // Index 0 is the root of the trust chain.
}

// WriteCert writes the PEM-encoded ASN.1 DER representation of the certificate.
func (s Store) WriteCert(w io.Writer) error {
	return writePem(w, s.cert.Raw, "CERTIFICATE")
}

// WriteIntermediates writes the PEM-encoded ASN.1 DER representation of all of the intermediate certs in the chain of trust.
// This is a no-op if there are no intermediate certificates.
func (s Store) WriteIntermediates(w io.Writer) error {
	for i := len(s.chain) - 1; i >= 1; i-- {
		err := writePem(w, s.chain[i].Raw, "CERTIFICATE")
		if err != nil {
			return err
		}
	}

	return nil
}

// WriteRoot writes the PEM-encoded ASN.1 DER representation of the root certificate in the chain of trust.
// This is a no-op if this is a self-signed certificate.
func (s Store) WriteRoot(w io.Writer) error {
	if len(s.chain) == 0 {
		return nil
	}

	return writePem(w, s.chain[0].Raw, "CERTIFICATE")
}

// WriteKey writes the PEM-encoded ASN.1 DER representation of the private key.
func (s Store) WriteKey(w io.Writer) error {
	switch key := s.key.(type) {
	case *rsa.PrivateKey:
		return writePem(w, x509.MarshalPKCS1PrivateKey(key), "RSA PRIVATE KEY")

	case *ecdsa.PrivateKey:
		buf, err := x509.MarshalECPrivateKey(key)

		if err != nil {
			return err
		}

		return writePem(w, buf, "EC PRIVATE KEY")
	}

	return fmt.Errorf("Unknown private key type %v", s.key)
}

// Load creates a certificate from a PEM-encoded certificate and private key.
// It is assumed that this is a self-signed certificate.
func Load(certPem []byte, keyPem []byte) (Store, error) {
	certBlock, _ := pem.Decode(certPem)

	if certBlock == nil {
		return Store{}, fmt.Errorf("Certificate is not in PEM format")
	}

	if certBlock.Type != "CERTIFICATE" {
		return Store{}, fmt.Errorf("Provided certificate PEM is not a 'CERTIFICATE'")
	}

	cert, err := x509.ParseCertificate(certBlock.Bytes)

	if err != nil {
		return Store{}, err
	}

	keyBlock, _ := pem.Decode(keyPem)

	if keyBlock == nil {
		return Store{}, fmt.Errorf("Key is not in PEM format")
	}

	var key interface{}

	switch keyBlock.Type {
	case "PRIVATE KEY":
		key, err = x509.ParsePKCS8PrivateKey(keyBlock.Bytes)

	case "RSA PRIVATE KEY":
		key, err = x509.ParsePKCS1PrivateKey(keyBlock.Bytes)

	case "EC PRIVATE KEY":
		key, err = x509.ParseECPrivateKey(keyBlock.Bytes)

	default:
		return Store{}, fmt.Errorf("Provided key PEM is not a private key")
	}

	if err != nil {
		log.Fatal(err)
	}

	result := Store{
		key:   key,
		cert:  cert,
		chain: nil,
	}

	return result, nil
}

// New creates a new non-CA certificate signed by the current certificate.
func (s Store) New(commonName string, dnsNames []string, ipAddresses []string) (Store, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return Store{}, err
	}

	serial, err := rand.Int(rand.Reader, big.NewInt(2<<60))

	if err != nil {
		return Store{}, err
	}

	ips := make([]net.IP, len(ipAddresses))

	for i, v := range ipAddresses {
		ips[i] = net.ParseIP(v)
	}

	now := time.Now().UTC()

	template := x509.Certificate{
		SerialNumber: serial,
		Subject: pkix.Name{
			CommonName:         commonName,
			Organization:       []string{"jmn.link"},
			Country:            []string{"CA"},
			Locality:           []string{"Kingston"},
			OrganizationalUnit: []string{"IT"},
			Province:           []string{"Ontario"},
		},
		DNSNames:    dnsNames,
		IPAddresses: ips,
		NotAfter:    now.AddDate(100, 0, 0),
		NotBefore:   now,
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, s.cert, key.Public(), s.key)
	if err != nil {
		return Store{}, err
	}

	cert, err := x509.ParseCertificate(certBytes)

	if err != nil {
		return Store{}, err
	}

	chain := make([]*x509.Certificate, len(s.chain)+1)
	copy(chain, s.chain)
	chain[len(chain)-1] = s.cert

	result := Store{
		key,
		cert,
		chain,
	}

	return result, nil
}

func writePem(w io.Writer, b []byte, header string) error {
	err := pem.Encode(w, &pem.Block{Bytes: b, Type: header})

	if err != nil {
		return err
	}

	return nil
}
