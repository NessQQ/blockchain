// Copyright (C) 2026 Lenar Khabibulin <khabibulin.lenar@gmail.com>
// Licensed under the GNU Affero General Public License v3 (AGPL-3.0)
package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

// Transaction represents a transfer of token value within the blockchain ledger.
// Defining this here eliminates duplicate type errors across other files.
type Transaction struct {
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
	Coinbase bool    `json:"coinbase"` // Marked true exclusively for system block rewards
}

// Block maps the fundamental container unit of the historical ledger array.
type Block struct {
	Hash         string         `json:"hash"`         // Cryptographic hex signature of this block
	Data         string         `json:"data"`         // Optional block metadata or descriptions
	PrevHash     string         `json:"prevHash"`     // Reference hash link to the historical parent block
	Nonce        int            `json:"nonce"`        // Valid puzzle scalar variable derived by the miner
	Transactions []*Transaction `json:"transactions"` // Complete bundled list of state modifications
}

// ComputeHash flattens all internal block properties, including structural transaction fields,
// and maps them directly through an industrial-grade SHA-256 cryptographic gate.
func (b *Block) ComputeHash() {
	// Convert the numerical nonce into a byte layout representation
	nonceBytes := []byte(fmt.Sprintf("%d", b.Nonce))

	// Pre-allocate the slicing matrix to accurately fit:
	// 3 Header fields (Data, PrevHash, Nonce) + the block's transaction counts
	components := make([][]byte, 0, 3+len(b.Transactions))

	// Append raw core structural block header bounds
	components = append(components, []byte(b.Data))
	components = append(components, []byte(b.PrevHash))
	components = append(components, nonceBytes)

	// IMMUTABILITY LOOP: Serialize every transaction field into a deterministic string format.
	// This binds user records directly to the hash; changing even a single byte or value
	// later on will result in a radically different hash, breaking the whole chain.
	for _, tx := range b.Transactions {
		txString := fmt.Sprintf("%s->%s:%f(%t)", tx.Sender, tx.Receiver, tx.Amount, tx.Coinbase)
		components = append(components, []byte(txString))
	}

	// Compress the structured matrices down into a solitary byte slice pipeline
	concatenatedData := bytes.Join(components, []byte{})

	// Fire final flattened stream through the SHA-256 refinery
	computedHash := sha256.Sum256(concatenatedData)

	// FIX: Safely translate binary slice into a clean, un-corrupted hexadecimal string
	b.Hash = hex.EncodeToString(computedHash[:])
}

// CreateBlock instantiates a fresh block container, provisions a random initial nonce,
// executes the ProofOfWork mine loop, and records the validated cryptographic outputs.
func CreateBlock(data string, prevHash string, transactions []*Transaction) *Block {
	// Seed our localized random sequence parameters to safely vary puzzle starts
	seed := time.Now().UnixNano()
	source := rand.NewSource(seed)
	rng := rand.New(source)
	initialNonce := rng.Intn(10000)

	// Instantiate our raw un-mined block shell pointer reference
	block := &Block{
		Hash:         "",
		Data:         data,
		PrevHash:     prevHash,
		Nonce:        initialNonce,
		Transactions: transactions,
	}

	// Instantiate the Proof of Work engine node and initiate brute-force search routines
	newPow := NewProofOfWork(block)
	nonce, hash := newPow.MineBlock()

	// Capture the verified metrics discovered by the mining sequence pass
	block.Nonce = nonce
	block.Hash = hash

	return block
}

// Genesis manufactures the zero-index foundation node needed to safely bootstrap the blockchain network.
func Genesis() *Block {
	// Structure our initial Coinbase mint transaction to instantiate ledger history
	coinbaseTransaction := &Transaction{
		Sender:   "Coinbase",
		Receiver: "Genesis",
		Amount:   0.0,
		Coinbase: true,
	}

	// Initialize our system root block (with no parent PrevHash parameters possible)
	return CreateBlock("Genesis", "", []*Transaction{coinbaseTransaction})
}
