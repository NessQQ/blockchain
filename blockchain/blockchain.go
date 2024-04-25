package blockchain

// This file will contain the logic for implementing the blockchain.

type BlockChain struct {
	// This slice of Block pointers representing the blocks in the blockchain.
	Blocks []*Block
}

// Function to create the Genesis block and initialize a new blockchain with it. It takes no arguments and returns a pointer to a BlockChain object containing Genesis block.
func InitBlockChain() *BlockChain {
	// Create the Genesis block using the Genesis function, then create and return a new BlockChain object and initialize the Blocks attribute with a slice containing only the Genesis block.
	return &BlockChain{[]*Block{Genesis()}}
}

// Function to add a block to a blockcahin. It takes the data to be added in a block as an argument.
func (chain *BlockChain) AddBlock(data string, coinbaseRcpt string, transactions []*Transaction) {
	// Get the previous block in the blockchain.
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	coinbaseTransaction := &Transaction{
		Sender:   "Coinbase",
		Receiver: coinbaseRcpt,
		Amount:   10.0,
		Coinbase: true,
	}
	// Create a new block using the data and the previous hash using the CreateBlock function. It passes the data and the previous block’s hash as arguments.
	newBlock := CreateBlock(data, prevBlock.Hash, append([]*Transaction{coinbaseTransaction}, transactions...))
	// Append the newly created block to the Blocks attribute of the blockchain.
	chain.Blocks = append(chain.Blocks, newBlock)
}

// Structure for transactions.
type Transaction struct {
	// The sender’s address.
	Sender string
	// The receiver’s address.
	Receiver string
	// The amount being transferred.
	Amount float64
	// Indicator of whether the transaction is a coinbase transaction or not.
	Coinbase bool
}
