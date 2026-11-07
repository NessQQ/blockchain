// Copyright (C) 2026 Lenar Khabibulin <khabibulin.lenar@gmail.com>
// Licensed under the GNU Affero General Public License v3 (AGPL-3.0)
package blockchain

import (
	"crypto/rsa"
	"errors"
	"log"
	"fmt"
)

type BlockChain struct {
	Blocks []*Block
}

func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}

// AddBlock now acts as a strict security clearinghouse. 
// It requires an explicit map of Public Keys to cross-reference and validate transaction signatures 
// before letting them contaminate the immutable ledger array history.
func (chain *BlockChain) AddBlock(data string, coinbaseRcpt string, transactions []*Transaction, publicKeys map[string]*rsa.PublicKey) error {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	
	// 1. SECURITY GATEWAY: Validate every single incoming transaction signature
	for _, tx := range transactions {
		// System-minted coinbase transactions don't have user signatures, skip them safely
		if tx.Coinbase {
			continue
		}

		// Fetch the sender's public key from our trusted network registry map
		pubKey, exists := publicKeys[tx.Sender]
		if !exists {
			return errors.New("security violation: transaction rejected - sender public key not found in network registry")
		}

		// Execute your wallet module's validation algorithms directly at the chain boundary gate
		err := VerifyTransaction(tx, pubKey, tx.Signature)
		if err != nil {
			log.Printf("⚠️  Fraudulent transaction blocked! Sender: %s attempted unauthorized transfer.", tx.Sender)
			return fmt.Errorf("transaction verification failed: %w", err)
		}
	}

	// 2. MINT COINBASE REWARD: If all user signatures clear validation, safely mint the mining subsidy block reward
	coinbaseTransaction := &Transaction{
		Sender:   "Coinbase",
		Receiver: coinbaseRcpt,
		Amount:   10.0,
		Coinbase: true,
	}
	
	fullTransactionSet := append([]*Transaction{coinbaseTransaction}, transactions...)
	
	// 3. MINE AND CHAIN: Pass the validated data bundle down into the SHA-256 MineBlock loop
	newBlock := CreateBlock(data, prevBlock.Hash, fullTransactionSet)
	chain.Blocks = append(chain.Blocks, newBlock)
	
	return nil
}
