// Copyright (C) 2026 Lenar Khabibulin <khabibulin.lenar@gmail.com>
// Licensed under the GNU Affero General Public License v3 (AGPL-3.0)
package main

import (
	"crypto/rsa"
	"fmt"
	"log"
	"github.com/NessQQ/blockchain/blockchain"
)

func main() {
	fmt.Println("============ BLOCKCHAIN TESTING ENGINE ============")

	// 1. Initialize a clean, trusted Network Public Key Registry.
	// In a real distributed system, nodes exchange public keys to verify inbound sender claims.
	networkPublicKeys := make(map[string]*rsa.PublicKey)

	// 2. Instantiate individual user cryptographic wallets (Alice and Bob)
	fmt.Println("\n🔑 Generating secure 2048-bit RSA keys for network participants...")
	
	aliceWallet, err := blockchain.NewWallet()
	if err != nil {
		log.Fatalf("Fatal: Failed to generate Alice's wallet keys: %v", err)
	}
	// Alice's wallet address string is represented by her unique public key structure pointer
	aliceAddress := fmt.Sprintf("%p", aliceWallet.PublicKey)
	networkPublicKeys[aliceAddress] = aliceWallet.PublicKey
	fmt.Printf("[Registry] Alice Wallet Activated. Address: %s\n", aliceAddress)

	bobWallet, err := blockchain.NewWallet()
	if err != nil {
		log.Fatalf("Fatal: Failed to generate Bob's wallet keys: %v", err)
	}
	bobAddress := fmt.Sprintf("%p", bobWallet.PublicKey)
	networkPublicKeys[bobAddress] = bobWallet.PublicKey
	fmt.Printf("[Registry] Bob Wallet Activated.   Address: %s\n", bobAddress)

	// 3. Initialize the sequential ledger tracking chain (Processes the Genesis block instantly)
	fmt.Println("\n⚙️  Bootstrapping Blockchain Ledger Architecture...")
	chain := blockchain.InitBlockChain()
	
	// Print Genesis block characteristics to confirm hex string encoding formatting success
	genesisBlock := chain.Blocks[0]
	fmt.Printf("✔ Genesis Block Created! Hash: %s\n", genesisBlock.Hash)
	fmt.Printf("  Initial Transaction: %s -> %s (%f coins, Coinbase: %t)\n", 
		genesisBlock.Transactions[0].Sender, 
		genesisBlock.Transactions[0].Receiver, 
		genesisBlock.Transactions[0].Amount, 
		genesisBlock.Transactions[0].Coinbase,
	)

	// 4. Construct a standard Peer-to-Peer payment transaction order (Alice transfers funds to Bob)
	fmt.Println("\n💸 Constructing outbound transaction: Alice sending 4.5 tokens to Bob...")
	tx := &blockchain.Transaction{
		Sender:   aliceAddress,
		Receiver: bobAddress,
		Amount:   4.5,
		Coinbase: false,
	}

	// 5. SIGN PHASE: Lock transaction data payload using Alice's private key
	fmt.Println("✍️  Signing transaction payload with Alice's private key context...")
	signature, err := aliceWallet.SignTransaction(tx)
	if err != nil {
		log.Fatalf("Fatal: Transaction signing routine failed: %v", err)
	}
	fmt.Printf("✔ Digital Signature Generated! Base64 Segment: %s...\n", signature[:30])

	// 6. CHAIN INJECTION: Group the transaction and attempt to mine a new block
	// We pass the transaction along with our public key registry map so the chain can run its security gate checks
	fmt.Println("\n📥 Submitting transaction bundle to BlockChain.AddBlock gateway...")
	
	// Miner address for block rewards is assigned to Bob's address string marker
	minerAddress := bobAddress
	err = chain.AddBlock("Block 1 Data: Payment Processing", minerAddress, []*blockchain.Transaction{tx}, networkPublicKeys)
	if err != nil {
		log.Fatalf("❌ Block rejected at the validation gate: %v", err)
	}

	// 7. VERIFICATION MATRIX: Audit the newly appended block state values
	minedBlock := chain.Blocks[len(chain.Blocks)-1]
	fmt.Println("===============================================================")
	fmt.Printf("📦 Mined Block Index:  %d\n", len(chain.Blocks)-1)
	fmt.Printf("🔒 Block Identifier:   %s\n", minedBlock.Hash)
	fmt.Printf("🔗 Link to Parent:     %s\n", minedBlock.PrevHash)
	fmt.Printf("🔢 Winning Nonce:      %d\n", minedBlock.Nonce)
	fmt.Printf("🧾 Total Transactions: %d\n", len(minedBlock.Transactions))
	
	// Double check that the coinbase transaction block reward printed correctly as element 0
	cbTx := minedBlock.Transactions[0]
	userTx := minedBlock.Transactions[1]
	fmt.Printf("   -> Tx 0 (Reward):  %s minted %0.2f tokens to Miner (%s)\n", cbTx.Sender, cbTx.Amount, cbTx.Receiver[:12])
	fmt.Printf("   -> Tx 1 (Payment): %s transferred %0.2f tokens to %s\n", userTx.Sender[:12], userTx.Amount, userTx.Receiver[:12])
	fmt.Printf("      Cryptographic Signature Proof: [Valid Asset Pass]\n")
	fmt.Println("===============================================================")
	
	fmt.Println("\n✨ System execution pipeline completed successfully. Code compilation pristine.")
}
