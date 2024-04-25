package blockchain

// This file will contain the logic for managing user wallets.

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
)

// Wallet is a digital container that stores the cryptographic keys used to access and manage a user’s assets.
type Wallet struct {
	// Creates digital signatures, signs transactions, and decrypts data encrypted with the corresponding public key.
	PrivateKey *rsa.PrivateKey
	// Mainly used to receive funds and verify the authenticity of digital signatures.
	PublicKey *rsa.PublicKey
}

// It generates a pair of RSA public and private keys.
func GenerateRSAKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		return nil, nil, err
	}

	publicKey := &privateKey.PublicKey

	return privateKey, publicKey, nil
}

// Function that creates a new wallet instance.
func NewWallet() (*Wallet, error) {
	// It calls the GenerateRSAKeys function to obtain the private and public keys.
	privateKey, publicKey, err := GenerateRSAKeys()

	if err != nil {
		return nil, err
	}

	// It creates and returns a new Wallet struct and assigns the generated private and public keys to the respective attributes.
	return &Wallet{
		PrivateKey: privateKey,
		PublicKey:  publicKey,
	}, nil
}

// Function that takes a transaction pointer as an argument
func (wallet *Wallet) SignTransaction(transaction *Transaction) (string, error) {
	// It generates a data string by concatenating the sender, receiver, amount, and transaction’s coinbase fields.
	dataString := fmt.Sprintf("%s%s%f%t", transaction.Sender, transaction.Receiver, transaction.Amount, transaction.Coinbase)

	// It hashes the generated data string.
	hashedData := sha256.Sum256([]byte(dataString))

	// It signs the hashed data using the current wallet’s private key and obtains the signature.
	signature, err := rsa.SignPKCS1v15(rand.Reader, wallet.PrivateKey, crypto.SHA256, hashedData[:])

	if err != nil {
		return "", err
	}

	// It encodes the signature as a base64 string and returns the encoded signature.
	return base64.StdEncoding.EncodeToString(signature), nil
}

// Function that takes a transaction, a public key, and a signature string as arguments.
func VerifyTransaction(transaction *Transaction, publicKey *rsa.PublicKey, signature string) error {
	// It generates a data string by concatenating the sender, receiver, amount, and transaction’s coinbase fields.
	dataString := fmt.Sprintf("%s%s%f%t", transaction.Sender, transaction.Receiver, transaction.Amount, transaction.Coinbase)

	// It hashes the generated data string.
	hashedData := sha256.Sum256([]byte(dataString))

	// It decodes the signature string to obtain the signature bytes.
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)

	if err != nil {
		return err
	}

	// It verifies the signature using the provided public key.
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hashedData[:], signatureBytes)

	// It returns nil if the signature is valid and an error if the signature is invalid.
	if err != nil {
		return errors.New("Transaction Signature not valid")
	}
	return nil
}
