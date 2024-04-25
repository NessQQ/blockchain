package main

// This file will contain the code for testing the blockchain.

import (
	"fmt"
	"strconv"

	"github.com/NessQQ/blockchain/blockchain"
)

func main() {
	chain := blockchain.InitBlockChain()

	// Create a wallet for User1.
	userOneWallet, err := blockchain.NewWallet()

	if err != nil {
		fmt.Println("Error creating User1 wallet:", err)
		return
	}
	fmt.Println("User1 wallet created successfully")

	// Create a wallet for User2.
	userTwoWallet, err := blockchain.NewWallet()

	if err != nil {
		fmt.Println("Error creating User2's wallet:", err)
		return
	}
	fmt.Println("User2's wallet created successfully")

	// Create a transaction from User1 to User2.
	tx := &blockchain.Transaction{
		Sender:   userOneWallet.PublicKey.N.String(),
		Receiver: userTwoWallet.PublicKey.N.String(),
		Amount:   5.0,
	}
	fmt.Println("User1 to User2 Transaction created successfully")

	// Sign the transaction with User1’s wallet.
	signature, err := userOneWallet.SignTransaction(tx)
	if err != nil {
		fmt.Println("Error signing the transaction:", err)
		return
	}

	// Verify the transaction using User1’s wallet, public key, and the signature.
	err = blockchain.VerifyTransaction(tx, userOneWallet.PublicKey, signature)
	if err != nil {
		fmt.Println("Transaction verification failed:", err)
		return
	}
	fmt.Println("Transaction Verified Successfully")

	// Add the verified transaction to the blockchain.
	chain.AddBlock("Block 1", "User1", []*blockchain.Transaction{tx})
	fmt.Println()

	// Print the blockchain.
	for _, block := range chain.Blocks {
		fmt.Printf("Previous hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)

		pow := blockchain.NewProofOfWork(block)
		fmt.Printf("IsValidPoW: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		fmt.Println("Transactions:")

		for _, tx := range block.Transactions {
			fmt.Printf("Sender: %s\n", tx.Sender)
			fmt.Printf("Receiver: %s\n", tx.Receiver)
			fmt.Printf("Amount: %f\n", tx.Amount)
			fmt.Printf("Coinbase: %t\n", tx.Coinbase)
			fmt.Println()
		}
		fmt.Println()
	}
}
