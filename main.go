// Copyright (C) 2026 Lenar Khabibulin <khabibulin.lenar@gmail.com>
// Licensed under the GNU Affero General Public License v3 (AGPL-3.0)
package main

import (
	"crypto/rsa"
	"fmt"
	"log"

	"github.com/NessQQ/blockchain/blockchain"
	"github.com/NessQQ/blockchain/exchange"
)

func main() {
	fmt.Println("================== BLOCKCHAIN TESTING ENGINE ==================")

	// 1. Initialize our Blockchain public key registry context
	networkPublicKeys := make(map[string]*rsa.PublicKey)

	// 2. Generate secure identities (Wallets) for our market participants
	fmt.Println("\n🔑 Generating cryptographic wallet keypairs for Alice and Bob...")

	aliceWallet, err := blockchain.NewWallet()
	if err != nil {
		log.Fatalf("Fatal: Failed to provision Alice's identity: %v", err)
	}
	aliceAddress := fmt.Sprintf("%p", aliceWallet.PublicKey)
	networkPublicKeys[aliceAddress] = aliceWallet.PublicKey
	fmt.Printf("[Wallet] Alice Initialized. Network Address: %s\n", aliceAddress)

	bobWallet, err := blockchain.NewWallet()
	if err != nil {
		log.Fatalf("Fatal: Failed to provision Bob's identity: %v", err)
	}
	bobAddress := fmt.Sprintf("%p", bobWallet.PublicKey)
	networkPublicKeys[bobAddress] = bobWallet.PublicKey
	fmt.Printf("[Wallet] Bob Initialized.   Network Address: %s\n", bobAddress)

	// 3. Bootstrap our settlement structures and state engines
	fmt.Println("\n⚙️  Initializing In-Memory Order Book and Blockchain Ledger...")
	book := exchange.NewOrderBook()
	chain := blockchain.InitBlockChain()

	// 4. Populate the matching engine book with passive liquidity (Maker Orders)
	fmt.Println("\n📥 Bob places a passive MAKER Sell Order on the book ladder...")
	makerOrder := &exchange.Order{
		ID:        "ORDER-001",
		WalletID:  bobAddress,
		Side:      exchange.Sell,
		Price:     90000.0,
		Amount:    2.5,
		FilledQty: 0.0,
	}
	// Processing this order queues it safely inside Bob's Ask PriceLevel list since no bids match yet
	book.ProcessOrder(makerOrder)
	fmt.Printf("[Order Book] Ask PriceLevel active at $%0.2f (Volume: %0.2f tokens)\n", makerOrder.Price, makerOrder.Remaining())

	// 5. Inject an aggressive cross-market order (Taker Order) to trigger immediate matching
	fmt.Println("\n⚡ Alice fires an aggressive TAKER Buy Order hitting Bob's price tier...")
	takerOrder := &exchange.Order{
		ID:        "ORDER-002",
		WalletID:  aliceAddress,
		Side:      exchange.Buy,
		Price:     90000.0,
		Amount:    1.5, // Alice buys 1.5 out of Bob's 2.5 tokens available
		FilledQty: 0.0,
	}

	// EXECUTE IN-MEMORY MATCHING: The FIFO pipeline crosses the order parameters instantly
	tradeLogs := book.ProcessOrder(takerOrder)

	if len(tradeLogs) == 0 {
		log.Fatalf("Fatal Error: Trade matching engine failed to calculate crossing parameters.")
	}

	fmt.Printf("🎯 Match Found! Matched Volume: %0.2f tokens at $%0.2f\n", tradeLogs[0].Amount, tradeLogs[0].Price)

	// 6. BRIDGE LAYER: Convert memory trade logs into cryptographically signed ledger transactions
	fmt.Println("\n🔗 Converting trade execution reports into secure blockchain transactions...")
	var signedMempool []*blockchain.Transaction

	for _, trade := range tradeLogs {
		// Map our internal memory matching entities directly to the blockchain.Transaction schema structures
		tx := &blockchain.Transaction{
			Sender:   aliceAddress,               // Taker Buyer transfers token capital value
			Receiver: bobAddress,                 // Maker Seller receives funds
			Amount:   trade.Amount * trade.Price, // Total transaction value size in fiat/USDT units
			Coinbase: false,
		}

		// Secure the movement explicitly using Alice's private key wallet signature routines
		fmt.Println("✍️  Signing settlement ledger frame using Alice's private wallet keys...")
		_, err := aliceWallet.SignTransaction(tx)
		if err != nil {
			log.Fatalf("Fatal Error: Transaction signing loop failed: %v", err)
		}

		signedMempool = append(signedMempool, tx)
	}

	// 7. SETTLEMENT LAYER: Clear transactions through the consensus gate and mine the block
	fmt.Println("\n📥 Dispatching signed transactions context into BlockChain.AddBlock clearinghouse...")

	// Bob acts as the network block miner and intercepts the standard Coinbase block reward
	minerRecipient := bobAddress
	err = chain.AddBlock("Settlement Block: Engine Matches Clear", minerRecipient, signedMempool, networkPublicKeys)
	if err != nil {
		log.Fatalf("❌ Block settlement failed validation tests: %v", err)
	}

	// 8. TOTAL LEDGER METRIC AUDIT LOG
	minedBlock := chain.Blocks[len(chain.Blocks)-1]
	fmt.Println("\n=================== SYSTEM STATE SETTLEMENT AUDIT ===================")
	fmt.Printf("📦 Confirmed Block Index: %d\n", len(chain.Blocks)-1)
	fmt.Printf("🔒 Block Cryptographic ID: %s\n", minedBlock.Hash)
	fmt.Printf("🔗 Hex Link to Parent:    %s\n", minedBlock.PrevHash)
	fmt.Printf("🔢 Derived Miner Nonce:   %d\n", minedBlock.Nonce)
	fmt.Printf("🧾 Total Transactions:    %d\n", len(minedBlock.Transactions))

	cb := minedBlock.Transactions[0]
	settledTx := minedBlock.Transactions[1]
	fmt.Printf("   -> Tx 0 [Mint]:   %s awarded %0.1f coins to Miner (%s...)\n", cb.Sender, cb.Amount, cb.Receiver[:10])
	fmt.Printf("   -> Tx 1 [Settle]: %s paid $%0.2f USDT to %s...\n", settledTx.Sender[:10], settledTx.Amount, settledTx.Receiver[:10])
	fmt.Printf("      Cryptographic Proof: [Signature Verified - Ledger Sealed Immutably]\n")
	fmt.Println("=====================================================================")

	fmt.Println("\n✨ Comprehensive pipeline execution complete. Low-latency logic loop pristine.")
}
