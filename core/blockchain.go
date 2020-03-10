package core

import (
	// "bytes"
	// "crypto/ecdsa"
	// "encoding/hex"
	// "errors"
	// "fmt"
	"log"
	// "os"

	"github.com/etcd-io/bbolt"
)

const dbFile = "database/blockchain.db"
const blocksBucket = "blocks"
const genesisCoinbaseData = "The Times 03/Jan/2009 Chancellor on brink of second bailout for banks"

// Blockchain implements interactions with a DB
type Blockchain struct {
	tip []byte
	Db  *bbolt.DB
}

// AddBlock add the block into the blockchain
func (bc *Blockchain) AddBlock(data string) {
	var lashHash []byte

	err := bc.Db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		lashHash = b.Get([]byte("1"))

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	newBlock := NewBlock(data, lashHash)

	err = bc.Db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		err := b.Put(newBlock.Hash, newBlock.Serialize())
		if err != nil {
			log.Panic(err)
		}

		err = b.Put([]byte("1"), newBlock.Hash)
		if err != nil {
			log.Panic(err)
		}
		bc.tip = newBlock.Hash

		return nil
	})
	if err != nil {
		log.Panic(err)
	}
}

func NewBlockchain() *Blockchain {
	var tip []byte
	db, err := bbolt.Open(dbFile, 0600, nil)
	if err != nil {
		log.Panic(err)
	}

	err = db.Update(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		if b == nil {
			genisis := NewGenesisBlock()
			b, err := tx.CreateBucket([]byte(blocksBucket))
			if err != nil {
				log.Panic(err)
			}

			err = b.Put(genisis.Hash, genisis.Serialize())
			if err != nil {
				log.Panic(err)
			}

			err = b.Put([]byte("1"), genisis.Hash)
			if err != nil {
				log.Panic(err)
			}
			tip = genisis.Hash
		} else {
			tip = b.Get([]byte("1"))
		}

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	bc := Blockchain{tip, db}

	return &bc
}
