package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strconv"
	"time"
)

const targetBits = 24 // Difficulty of the Proof-of-Work algorithm

// Block represents a block in the blockchain
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Hash          []byte
	Nonce         int
}

// SetHash calculates and sets the hash of the block using Proof of Work
func (b *Block) SetHash() {
	var (
		hashInt big.Int
		hash    [32]byte
		nonce   int = 0
	)
	target := big.NewInt(1)
	target.Lsh(target, uint(256-targetBits))

	for {
		data := bytes.Join(
			[][]byte{
				b.PrevBlockHash,
				b.Data,
				[]byte(strconv.FormatInt(b.Timestamp, 10)),
				[]byte(strconv.Itoa(nonce)),
			},
			[]byte{},
		)
		hash = sha256.Sum256(data)
		hashInt.SetBytes(hash[:])
		if hashInt.Cmp(target) == -1 {
			break
		} else {
			nonce++
		}
	}
	b.Hash = hash[:]
	b.Nonce = nonce
}

// NewBlock creates a new block with Proof of Work
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, []byte{}, 0}
	block.SetHash()
	return block
}

// Blockchain maintains the chain of blocks
type Blockchain struct {
	blocks []*Block
}

// AddBlock adds a new block to the blockchain after verifying Proof of Work
func (bc *Blockchain) AddBlock(data string) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newBlock := NewBlock(data, prevBlock.Hash)
	bc.blocks = append(bc.blocks, newBlock)
}

// NewGenesisBlock creates the genesis block
func NewGenesisBlock() *Block {
	return NewBlock("Genesis Block", []byte{})
}

// NewBlockchain initializes a new Blockchain with a genesis block
func NewBlockchain() *Blockchain {
	return &Blockchain{[]*Block{NewGenesisBlock()}}
}

func main() {
	bc := NewBlockchain()
	bc.AddBlock("Send 12 doge to Aryan")
	bc.AddBlock("Send 3 more doge to Sai")

	for _, block := range bc.blocks {
		fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
		fmt.Printf("Data: %s\n", block.Data)
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Println()
	}
}
