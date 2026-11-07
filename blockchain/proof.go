// Copyright (C) 2026 Lenar Khabibulin <khabibulin.lenar@gmail.com>
// Licensed under the GNU Affero General Public License v3 (AGPL-3.0)
package blockchain

// This file handles the cryptographic mining puzzle metrics and block validation routines.

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Difficulty sets the leading zero-bit requirement for our mining puzzle.
// A setting of 10-14 runs rapidly on test machines for fast development cycles.
const Difficulty = 12

type ProofOfWork struct {
	Block  *Block
	Target *big.Int
}

// NewProofOfWork calculates the maximum integer upper bound target based on system Difficulty rules.
func NewProofOfWork(block *Block) *ProofOfWork {
	targetVal := big.NewInt(1)
	
	// Left shift a 1 by (256 - Difficulty) bits.
	// As difficulty rises, the target cap falls, making it harder to find a matching hash.
	targetVal.Lsh(targetVal, uint(256-Difficulty))
	
	return &ProofOfWork{block, targetVal}
}

// ComputeData aggregates all mutable block fields and every single transaction into a single byte stream.
func (pow *ProofOfWork) ComputeData(nonce int) []byte {
	// Pre-allocate binary safe buffers for our custom integer configurations
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, uint64(nonce))

	diffBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(diffBytes, uint64(Difficulty))

	// Dynamically size our allocation slice to avoid runtime growth copies:
	// 4 Core fields (PrevHash, Data, Nonce, Difficulty) + total transaction footprint
	components := make([][]byte, 0, 4+len(pow.Block.Transactions))
	components = append(components, []byte(pow.Block.PrevHash))
	components = append(components, []byte(pow.Block.Data))
	components = append(components, nonceBytes)
	components = append(components, diffBytes)

	// CRITICAL SYNCHRONIZATION POINT: We iterate through every transaction inside the block
	// using the exact same format signature string as block.go's ComputeHash() method.
	for _, tx := range pow.Block.Transactions {
		txString := fmt.Sprintf("%s->%s:%f(%t)", tx.Sender, tx.Receiver, tx.Amount, tx.Coinbase)
		components = append(components, []byte(txString))
	}

	return bytes.Join(components, []byte{})
}

// MineBlock executes the brute-force search routine, iterating nonces until it hits a valid SHA-256 target footprint.
func (pow *ProofOfWork) MineBlock() (int, string) {
	var intHash big.Int
	var computedHash [32]byte // Upgraded to 32 bytes to cleanly house SHA-256 calculations
	nonce := 0

	fmt.Printf("⛏️  Mining block containing %d transaction(s)...\n", len(pow.Block.Transactions))

	for {
		computedData := pow.ComputeData(nonce)
		computedHash = sha256.Sum256(computedData)
		
		// Print real-time rolling hex sequences cleanly over a single line (\r)
		fmt.Printf("\rCurrent Mining Hash Candidate: %x", computedHash)

		intHash.SetBytes(computedHash[:])

		// Check if the big integer hash representation falls beneath our strict target barrier
		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		
		nonce++
	}
	
	fmt.Printf("\n🎉 Block Successfully Mined! Discovered Nonce: %d\n\n", nonce)
	
	// Return the final successful values to be recorded inside the block context
	return nonce, hex.EncodeToString(computedHash[:])
}

// Validate provides our lightning-fast consensus verification check step.
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	computedData := pow.ComputeData(pow.Block.Nonce)
	computedHash := sha256.Sum256(computedData)
	
	intHash.SetBytes(computedHash[:])

	return intHash.Cmp(pow.Target) == -1
}
