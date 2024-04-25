package blockchain

// This file will contain the logic for creating and managing blocks of a blockchain.

import (
	"bytes"
	"crypto/md5"
	"math/rand"
	"time"
)

type Block struct {
	// The hash of the current block.
	Hash string
	// The data stored in the block.
	Data string
	// The hash of the previous block in the chain.
	PrevHash string
	// the nonce of the block
	Nonce int
	// Transactions
	Transactions []*Transaction
}

// Function to compute the hash of a block by using the md5.Sum() function to apply the MD5 hashing algorithm for calculating the hash of the concatenated data.
func (b *Block) ComputeHash() {
	// Concatenate the block’s data and previous hash.
	concatenatedData := bytes.Join([][]byte{[]byte(b.Data), []byte(b.PrevHash)}, []byte{})
	// Calculate the hash of the concatenated data.
	computedHash := md5.Sum(concatenatedData)
	// Assign the calculated hash value to the Hash attribute of the block.
	b.Hash = string(computedHash[:])
}

// Function to generate new blocks in the blockchain, each connected to the previous block. It takes the block data the previous hash, and transactions as arguments.
func CreateBlock(data string, prevHash string, transactions []*Transaction) *Block {
	// Seed the random number generator and create an initial nonce random value
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)
	initialNonce := rng.Intn(10000)
	// Create a new block by initializing a Block object with the provided data, previous hash, initial nonce, and transactions.
	block := &Block{"", data, prevHash, initialNonce, transactions}
	// Create a new ProofOfWork instance for the block.
	newPow := NewProofOfWork(block)
	// Obtain the nonce and derived hash values by calling the MineBlock function on the ProofOfWork instance.
	nonce, hash := newPow.MineBlock()
	// Set the derived hash value.
	block.Hash = string(hash[:])
	// Set the derived nonce value.
	block.Nonce = nonce
	// Return the created block.
	return block
}

// Function to generate Genesis block in the blockchain, not connected to any previous blocks. It takes no arguments.
func Genesis() *Block {
	// Initial coinbase transaction
	coinbaseTransaction := &Transaction{
		Sender:   "Coinbase",
		Receiver: "Genesis",
		Amount:   0.0,
		Coinbase: true,
	}
	// Create and return a pointer to a Block object, which contains hardcoded data and an empty previous hash.
	return CreateBlock("Genesis", "", []*Transaction{coinbaseTransaction})
}
