package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

func main() {
	fmt.Println("Generating the FinOps webhook certificate...")

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Println("Private key error:", err)
		os.Exit(1)
	}

	serialNumber, err := rand.Int(
		rand.Reader,
		new(big.Int).Lsh(big.NewInt(1), 128),
	)
	if err != nil {
		fmt.Println("Serial number error:", err)
		os.Exit(1)
	}

	certificate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"FinOps Admission Controller"},
			CommonName:   "finops-webhook.default.svc",
		},
		NotBefore: time.Now().Add(-5 * time.Minute),
		NotAfter:  time.Now().AddDate(1, 0, 0),

		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage: x509.KeyUsageDigitalSignature |
			x509.KeyUsageKeyEncipherment |
			x509.KeyUsageCertSign,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		DNSNames: []string{
			"localhost",
			"finops-webhook",
			"finops-webhook.default",
			"finops-webhook.default.svc",
			"finops-webhook.default.svc.cluster.local",
		},
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader,
		&certificate,
		&certificate,
		&key.PublicKey,
		key,
	)
	if err != nil {
		fmt.Println("Certificate error:", err)
		os.Exit(1)
	}

	if err := os.MkdirAll("certs", 0755); err != nil {
		fmt.Println("Folder error:", err)
		os.Exit(1)
	}

	certFile, err := os.Create("certs/tls.crt")
	if err != nil {
		fmt.Println("Certificate file error:", err)
		os.Exit(1)
	}

	err = pem.Encode(certFile, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	certFile.Close()
	if err != nil {
		fmt.Println("Certificate writing error:", err)
		os.Exit(1)
	}

	keyFile, err := os.OpenFile(
		"certs/tls.key",
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0600,
	)
	if err != nil {
		fmt.Println("Key file error:", err)
		os.Exit(1)
	}

	err = pem.Encode(keyFile, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	keyFile.Close()
	if err != nil {
		fmt.Println("Key writing error:", err)
		os.Exit(1)
	}

	fmt.Println("Certificate created: certs/tls.crt")
	fmt.Println("Private key created: certs/tls.key")
}
