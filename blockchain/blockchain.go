// Copyright (C) 2026 Lenar Khabibulin <khabibulin.lenar@gmail.com>
// Licensed under the GNU Affero General Public License v3 (AGPL-3.0)
package blockchain

// This file manages the blockchain chronological state sequence and block addition loops.

type BlockChain struct {
	// Blocks represents the persistent, sequential in-memory array ledger of mined blocks.
	Blocks []*Block
}

// InitBlockChain boots the immutable ledger system by processing the hardcoded Genesis block.
func InitBlockChain() *BlockChain {
	// Bootstrap the chain array context with only the initial Genesis node block
	return &BlockChain{[]*Block{Genesis()}}
}

// AddBlock handles transaction grouping, automatically prepends the network coinbase mining block reward,
// triggers the proof-of-work puzzle loop, and links the completed container into the chain.
func (chain *BlockChain) AddBlock(data string, coinbaseRcpt string, transactions []*Transaction) {
	// Fetch the absolute tip of the current chain to inherit its parent cryptographic hash marker
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	
	// Mint the standard block subsidy reward transaction directly to the block miner's recipient address.
	// This mirrors real-world mining mechanics where the network rewards the miner for solving the PoW puzzle.
	coinbaseTransaction := &Transaction{
		Sender:   "Coinbase",
		Receiver: coinbaseRcpt,
		Amount:   10.0, // Fixed 10 unit reward token payout
		Coinbase: true,
	}
	
	// Append any user-submitted ledger transactions directly behind the system coinbase reward transaction
	fullTransactionSet := append([]*Transaction{coinbaseTransaction}, transactions...)
	
	// Create the block, which kicks off our SHA-256 MineBlock loop internally
	newBlock := CreateBlock(data, prevBlock.Hash, fullTransactionSet)
	
	// Append the validated, mined block securely onto the sequential ledger tracking array
	chain.Blocks = append(chain.Blocks, newBlock)
}
