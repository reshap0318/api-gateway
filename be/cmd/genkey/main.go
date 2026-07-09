package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/reshap0318/api-gateway/internal/helpers"
)

func main() {
	// Parse flags
	force := flag.Bool("f", false, "Force overwrite existing keys without confirmation")
	flag.Parse()

	// Find .env file
	envPath, err := findEnvFile()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	// Check if keys already exist
	privateKeyPath := "storage/keys/private.pem"
	publicKeyPath := "storage/keys/public.pem"

	if _, err := os.Stat(privateKeyPath); err == nil && !*force {
		fmt.Println("⚠️  WARNING: Existing keys found!")
		fmt.Println("⚠️  This will overwrite existing keys!")
		fmt.Println("⚠️  All existing JWT tokens will be invalidated!")
		fmt.Print("❓ Continue? (y/n): ")
		
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" {
			fmt.Println("❌ Operation cancelled")
			return
		}
		fmt.Println()
	}

	// Generate RSA key pair
	privateKey, publicKey, err := generateRSAKeys()
	if err != nil {
		log.Fatalf("Error generating keys: %v", err)
	}

	// Generate secure passphrase
	passphrase, err := helpers.GenerateRandomString(32)
	if err != nil {
		log.Fatalf("Error generating passphrase: %v", err)
	}

	// Create keys directory
	if err := os.MkdirAll("storage/keys", 0700); err != nil {
		log.Fatalf("Error creating storage/keys directory: %v", err)
	}

	// Save private key (encrypted)
	if err := saveEncryptedPrivateKey(privateKeyPath, privateKey, passphrase); err != nil {
		log.Fatalf("Error saving private key: %v", err)
	}

	// Save public key
	if err := savePublicKey(publicKeyPath, publicKey); err != nil {
		log.Fatalf("Error saving public key: %v", err)
	}

	// Save passphrase to file
	passphraseFilePath := "storage/keys/passphrase"
	if err := savePassphraseFile(passphraseFilePath, passphrase); err != nil {
		log.Fatalf("Error saving passphrase file: %v", err)
	}

	// Update .env file (paths only, no passphrase)
	if err := updateEnvFile(envPath, privateKeyPath, publicKeyPath, passphraseFilePath); err != nil {
		log.Fatalf("Error updating .env file: %v", err)
	}

	fmt.Println("✅ Keys generated successfully!")
	fmt.Printf("📁 Private key: %s (encrypted)\n", privateKeyPath)
	fmt.Printf("📁 Public key: %s\n", publicKeyPath)
	fmt.Printf("📁 Passphrase: %s\n", passphraseFilePath)
	fmt.Printf("📁 Updated: %s\n", envPath)
	fmt.Println("⚠️  Keep storage/keys/ out of version control!")
	fmt.Println("⚠️  If you regenerate keys, all existing JWT tokens will be invalid!")
}

// generateRSAKeys generates RSA key pair (2048 bits)
func generateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate RSA private key: %w", err)
	}

	return privateKey, &privateKey.PublicKey, nil
}

// saveEncryptedPrivateKey saves private key in PEM format with passphrase encryption
func saveEncryptedPrivateKey(path string, privateKey *rsa.PrivateKey, passphrase string) error {
	// Marshal private key to DER format
	derBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// Encrypt with passphrase using PEM
	block, err := x509.EncryptPEMBlock(
		rand.Reader,
		"RSA PRIVATE KEY",
		derBytes,
		[]byte(passphrase),
		x509.PEMCipherAES256,
	)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Write to file
	privateKeyPEM := pem.EncodeToMemory(block)
	if err := os.WriteFile(path, privateKeyPEM, 0600); err != nil {
		return fmt.Errorf("failed to write private key file: %w", err)
	}

	return nil
}

// savePublicKey saves public key in PEM format
func savePublicKey(path string, publicKey *rsa.PublicKey) error {
	// Marshal public key to PKIX ASN.1 DER
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}

	// Encode to PEM
	block := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	}
	publicKeyPEM := pem.EncodeToMemory(block)

	// Write to file
	if err := os.WriteFile(path, publicKeyPEM, 0644); err != nil {
		return fmt.Errorf("failed to write public key file: %w", err)
	}

	return nil
}

// savePassphraseFile writes the passphrase to a restricted file.
func savePassphraseFile(path, passphrase string) error {
	if err := os.WriteFile(path, []byte(passphrase), 0600); err != nil {
		return fmt.Errorf("failed to write passphrase file: %w", err)
	}
	return nil
}

// updateEnvFile updates .env with new JWT key paths (passphrase is stored in a file, not env).
func updateEnvFile(envPath, privateKeyPath, publicKeyPath, passphraseFilePath string) error {
	content, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("failed to read .env file: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	lines = updateOrAddEnv(lines, "JWT_PRIVATE_KEY_PATH", privateKeyPath)
	lines = updateOrAddEnv(lines, "JWT_PUBLIC_KEY_PATH", publicKeyPath)
	lines = updateOrAddEnv(lines, "JWT_PASSPHRASE_PATH", passphraseFilePath)

	// Remove old JWT_PASSPHRASE entry if present
	filtered := lines[:0]
	for _, line := range lines {
		if !strings.HasPrefix(line, "JWT_PASSPHRASE=") {
			filtered = append(filtered, line)
		}
	}

	updated := strings.Join(filtered, "\n")
	if err := os.WriteFile(envPath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("failed to write .env file: %w", err)
	}

	return nil
}

// updateOrAddEnv updates existing env var or adds new one
func updateOrAddEnv(lines []string, key, value string) []string {
	prefix := key + "="
	found := false

	for i, line := range lines {
		if strings.HasPrefix(line, prefix) {
			lines[i] = prefix + value
			found = true
			break
		}
	}

	if !found {
		lines = append(lines, prefix+value)
	}

	return lines
}

// findEnvFile searches for .env file in current directory and parent directories
func findEnvFile() (string, error) {
	// Check current directory first
	if _, err := os.Stat(".env"); err == nil {
		return ".env", nil
	}

	// Get absolute path of current directory
	absPath, err := filepath.Abs(".")
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Search in parent directories (max 5 levels)
	current := absPath
	for i := 0; i < 4; i++ {
		envPath := filepath.Join(current, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}
		current = filepath.Dir(current)
	}

	return "", fmt.Errorf(".env file not found (create one or run from project root)")
}
