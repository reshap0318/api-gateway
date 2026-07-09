package services

import (
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"

	"github.com/reshap0318/api-gateway/internal/helpers"
)

// JWKS represents a JSON Web Key Set response
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWK represents a JSON Web Key
type JWK struct {
	KeyType     string `json:"kty"`
	PublicKeyUse string `json:"use"`
	KeyID       string `json:"kid"`
	Algorithm   string `json:"alg"`
	Modulus     string `json:"n"`
	Exponent    string `json:"e"`
}

// JWKSManager manages RSA keys for JWT operations
type JWKSManager struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	keyID      string
}

// Initialize loads and validates RSA keys
func (j *JWKSManager) Initialize(privateKeyPath, publicKeyPath, passphrase string) error {
	// Load private key
	privateKey, err := helpers.LoadPrivateKey(privateKeyPath, passphrase)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}
	j.privateKey = privateKey

	// Load public key
	publicKey, err := helpers.LoadPublicKey(publicKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load public key: %w", err)
	}
	j.publicKey = publicKey

	// Generate key ID from public key fingerprint
	j.keyID = helpers.GenerateKeyID(publicKey)

	return nil
}

// GetPrivateKey returns the active private key
func (j *JWKSManager) GetPrivateKey() *rsa.PrivateKey {
	return j.privateKey
}

// GetPublicKey returns the active public key
func (j *JWKSManager) GetPublicKey() *rsa.PublicKey {
	return j.publicKey
}

// GetKeyID returns the active key ID
func (j *JWKSManager) GetKeyID() string {
	return j.keyID
}

// GetJWKS returns the JWKS response with active public key
func (j *JWKSManager) GetJWKS() JWKS {
	// Convert RSA public key to JWK format
	nBytes := j.publicKey.N.Bytes()
	eBytes := big.NewInt(int64(j.publicKey.E)).Bytes()

	return JWKS{
		Keys: []JWK{
			{
				KeyType:     "RSA",
				PublicKeyUse: "sig",
				KeyID:       j.keyID,
				Algorithm:   "RS256",
				Modulus:     base64.RawURLEncoding.EncodeToString(nBytes),
				Exponent:    base64.RawURLEncoding.EncodeToString(eBytes),
			},
		},
	}
}
