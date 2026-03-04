package ca

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// StorageDir returns the directory where CA and leaf cert files are stored.
// Uses $XDG_CONFIG_HOME/ennyn/ca/ if set, otherwise ~/.config/ennyn/ca/.
func StorageDir() (string, error) {
	base := os.Getenv("XDG_CONFIG_HOME")
	if base == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("cannot determine home directory: %w", err)
		}
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "ennyn", "ca"), nil
}

// CA holds a loaded or generated certificate authority.
type CA struct {
	Cert    *x509.Certificate
	Key     *ecdsa.PrivateKey
	CertPEM []byte
}

// LoadOrGenerate loads an existing CA from disk, or generates a new one if none exists.
// Returns the CA and whether it was freshly generated.
func LoadOrGenerate() (*CA, bool, error) {
	dir, err := StorageDir()
	if err != nil {
		return nil, false, err
	}

	certPath := filepath.Join(dir, "ca.crt")
	keyPath := filepath.Join(dir, "ca.key")

	// Try loading existing CA
	if ca, err := load(certPath, keyPath); err == nil {
		return ca, false, nil
	}

	// Generate new CA
	ca, err := generate()
	if err != nil {
		return nil, false, fmt.Errorf("generating CA: %w", err)
	}

	// Persist to disk
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, false, fmt.Errorf("creating CA directory: %w", err)
	}

	if err := writePEM(certPath, "CERTIFICATE", ca.CertPEM, 0644); err != nil {
		return nil, false, fmt.Errorf("writing CA certificate: %w", err)
	}

	keyDER, err := x509.MarshalECPrivateKey(ca.Key)
	if err != nil {
		return nil, false, fmt.Errorf("marshaling CA key: %w", err)
	}
	if err := writePEM(keyPath, "EC PRIVATE KEY", keyDER, 0600); err != nil {
		return nil, false, fmt.Errorf("writing CA key: %w", err)
	}

	return ca, true, nil
}

func load(certPath, keyPath string) (*CA, error) {
	certPEMBytes, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	keyPEMBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}

	certBlock, _ := pem.Decode(certPEMBytes)
	if certBlock == nil {
		return nil, fmt.Errorf("failed to decode CA certificate PEM")
	}
	cert, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing CA certificate: %w", err)
	}

	keyBlock, _ := pem.Decode(keyPEMBytes)
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to decode CA key PEM")
	}
	key, err := x509.ParseECPrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("parsing CA key: %w", err)
	}

	return &CA{
		Cert:    cert,
		Key:     key,
		CertPEM: certBlock.Bytes,
	}, nil
}

func generate() (*CA, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("generating serial number: %w", err)
	}

	pubKeyBytes, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("marshaling public key: %w", err)
	}
	subjectKeyID := sha1.Sum(pubKeyBytes)

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Ennyn Local CA",
			Organization: []string{"Ennyn"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(10 * 365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		SubjectKeyId:          subjectKeyID[:],
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, template, &key.PublicKey, key)
	if err != nil {
		return nil, fmt.Errorf("creating certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certDER)
	if err != nil {
		return nil, fmt.Errorf("parsing created certificate: %w", err)
	}

	return &CA{
		Cert:    cert,
		Key:     key,
		CertPEM: certDER,
	}, nil
}

func writePEM(path, blockType string, der []byte, perm os.FileMode) error {
	data := pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: der})
	return os.WriteFile(path, data, perm)
}

// Leaf holds a leaf certificate and its private key.
type Leaf struct {
	CertPEM []byte
	KeyPEM  []byte
}

// LoadOrGenerateLeaf loads an existing leaf cert or generates a new one signed by the CA.
// If forceRegen is true, always regenerates (used when CA was just created).
func LoadOrGenerateLeaf(authority *CA, forceRegen bool) (*Leaf, error) {
	dir, err := StorageDir()
	if err != nil {
		return nil, err
	}

	certPath := filepath.Join(dir, "localhost.crt")
	keyPath := filepath.Join(dir, "localhost.key")

	if !forceRegen {
		if leaf, err := loadLeaf(certPath, keyPath); err == nil {
			return leaf, nil
		}
	}

	leaf, err := generateLeaf(authority)
	if err != nil {
		return nil, fmt.Errorf("generating leaf certificate: %w", err)
	}

	if err := os.WriteFile(certPath, leaf.CertPEM, 0644); err != nil {
		return nil, fmt.Errorf("writing leaf certificate: %w", err)
	}
	if err := os.WriteFile(keyPath, leaf.KeyPEM, 0600); err != nil {
		return nil, fmt.Errorf("writing leaf key: %w", err)
	}

	return leaf, nil
}

func loadLeaf(certPath, keyPath string) (*Leaf, error) {
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, err
	}
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	return &Leaf{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}

func generateLeaf(authority *CA) (*Leaf, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating leaf key: %w", err)
	}

	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return nil, fmt.Errorf("generating serial number: %w", err)
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Ennyn localhost",
			Organization: []string{"Ennyn"},
		},
		DNSNames:    []string{"*.localhost", "localhost"},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(825 * 24 * time.Hour),
		KeyUsage:    x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	}

	certDER, err := x509.CreateCertificate(rand.Reader, template, authority.Cert, &key.PublicKey, authority.Key)
	if err != nil {
		return nil, fmt.Errorf("creating leaf certificate: %w", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDER})

	keyDER, err := x509.MarshalECPrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("marshaling leaf key: %w", err)
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: keyDER})

	return &Leaf{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}
