package blockchain

// This file will contain the logic for the proof of work (PoW) algorithm.

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"math/big"
)

// Define a variable called Difficulty that represents the desired difficulty level for mining new blocks.
const Difficulty = 10

type ProofOfWork struct {
	// The block to be validated
	Block *Block
	// The difficulty level for mining
	Target *big.Int
}

// Consensus algorithm to validate transactions and add new blocks. Takes a block as its parameter and returns a new ProofOfWork instance.
func NewProofOfWork(block *Block) *ProofOfWork {
	// Initialize a target variable with a value of 1.
	targetVal := big.NewInt(1)
	// Left shift the target’s value by (256 - Difficulty) bits. This sets the difficulty level.
	targetVal.Lsh(targetVal, uint(256-Difficulty))
	// Return a newly created ProofOfWork instance by assigning the block and target values.
	return &ProofOfWork{block, targetVal}
}

// Function that takes an integer nonce as a parameter and returns the prepared data for hashing the block.
func (pow *ProofOfWork) ComputeData(nonce int) []byte {
	// Concatenate the hash of the previous block and the data contained within the current block.
	data := bytes.Join(
		[][]byte{
			[]byte(pow.Block.PrevHash),
			[]byte(pow.Block.Data),
			make([]byte, 8),
			make([]byte, 8),
		},
		[]byte{},
	)
	// Concatenate the data from above with the nonce value.
	binary.BigEndian.PutUint64(data[len(data)-16:], uint64(nonce))
	// Concatenate the data from above with the difficulty level.
	binary.BigEndian.PutUint64(data[len(data)-8:], uint64(Difficulty))
	// Return the prepared data for hashing the block
	return data
}

// Function that performs the mining process to find a valid block hash that meets the target.
func (pow *ProofOfWork) MineBlock() (int, []byte) {
	// Integer value of the computed hash.
	var intHash big.Int
	// Computed hash byte slice.
	var computedHash [16]byte
	// Initialize a variable nonce with a value of 0.
	nonce := 0
	// Loop that runs until a valid block hash is found.
	for {
		// Call the ComputeData function with the current nonce value to obtain the data for hashing.
		computedData := pow.ComputeData(nonce)
		// Compute the hash of the data using the md5.Sum() function.
		computedHash = md5.Sum(computedData)
		// Print the computed hash.
		fmt.Printf("\r%x", computedHash)
		// Integer value of the computed hash set to an integer representation of the computed hash byte slice.
		intHash.SetBytes(computedHash[:])
		// Check if the computed hash value is less than the target defined for the PoW algorithm and if it is, break the loop because a valid block hash has been found.
		if intHash.Cmp(pow.Target) == -1 {
			break
		}
		// Increment the nonce to continue searching for a valid hash.
		nonce++
	}
	fmt.Println()
	// Return the nonce and the valid computed hash.
	return nonce, computedHash[:]
}

// Function that validates the derived hash and essentially validates the work done by the PoW algorithm. It takes no parameters and returns a boolean value indicating the validity of the derived hash.
func (pow *ProofOfWork) Validate() bool {
	// Integer value of the computed hash.
	var intHash big.Int
	// Call the ComputeData function with the block’s nonce value to obtain the data for hashing.
	computedData := pow.ComputeData(pow.Block.Nonce)
	// Compute the hash of the data using the md5.Sum() function.
	computedHash := md5.Sum(computedData)
	// Integer value of the computed hash set to an integer representation of the computed hash byte slice.
	intHash.SetBytes(computedHash[:])
	// Return true if the derived hash is less than the target, a valid hash, and false otherwise.
	if intHash.Cmp(pow.Target) == -1 {
		return true
	} else {
		return false
	}
}
