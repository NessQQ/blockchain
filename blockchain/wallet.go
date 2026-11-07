// Copyright (C) 2026 Lenar Khabibulin <khabibulin.lenar@gmail.com>
// Licensed under the GNU Affero General Public License v3 (AGPL-3.0)
package blockchain

// This file handles the logic for user cryptographic wallets, asymmetric key generation, and identity signatures.

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

// Wallet stores the cryptographic keys used to access balances and sign outgoing transactions.
type Wallet struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

// GenerateRSAKeys provisions a fresh asymmetric identity pair.
func GenerateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	// A 2048-bit key size is mathematically secure and highly performant for local simulation loops
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey
	return privateKey, publicKey, nil
}

// NewWallet creates an isolated user identity container.
func NewWallet() (*Wallet, error) {
	privateKey, publicKey, err := GenerateRSAKeys()
	if err != nil {
		return nil, err
	}

	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// SignTransaction accepts an outbound payment order, computes its payload signature,
// and saves the resulting base64 string directly onto the transaction container.
func (wallet *Wallet) SignTransaction(transaction *Transaction) (string, error) {
	// FIX: Synchronized to use the exact same layout signature format applied inside block.go and proof.go
	dataString := fmt.Sprintf("%s->%s:%f(%t)", transaction.Sender, transaction.Receiver, transaction.Amount, transaction.Coinbase)

	hashedData := sha256.Sum256([]byte(dataString))

	// Sign the SHA-256 hash using the user's private key
	signature, err := rsa.SignPKCS1v15(rand.Reader, wallet.PrivateKey, crypto.SHA256, hashedData[:])
	if err != nil {
		return "", err
	}

	base64Signature := base64.StdEncoding.EncodeToString(signature)
	
	// Store the generated signature directly inside the transaction so it travels through the blocks natively
	transaction.Signature = base64Signature

	return base64Signature, nil
}

// VerifyTransaction checks if a provided signature matches the immutable transaction properties.
func VerifyTransaction(transaction *Transaction, publicKey *rsa.PublicKey, signature string) error {
	// FIX: Must match the signature format exactly to generate matching hash evaluations
	dataString := fmt.Sprintf("%s->%s:%f(%t)", transaction.Sender, transaction.Receiver, transaction.Amount, transaction.Coinbase)

	hashedData := sha256.Sum256([]byte(dataString))

	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("failed to decode base64 signature string: %w", err)
	}

	// Validate the cryptographic proof against the owner's public key
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedData[:], signatureBytes)
	if err != nil {
		return errors.New("transaction signature verification failed: invalid cryptographic proof")
	}
	
	return nil
}
