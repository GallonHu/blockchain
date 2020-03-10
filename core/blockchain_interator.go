package core

import (
	"github.com/etcd-io/bbolt"
	"log"
)

// BlockchainInterator interators over blockchain blocks
type BlockchainInterator struct {
	currentHash []byte
	db          *bbolt.DB
}

func (bc *Blockchain) Interator() *BlockchainInterator {
	bci := &BlockchainInterator{bc.tip, bc.Db}

	return bci
}

// Next returns next block starting from the tip
func (i *BlockchainInterator) Next() *Block {
	var block *Block

	err := i.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock := b.Get(i.currentHash)
		block = DeserializeBlock(encodedBlock)

		return nil
	})
	if err != nil {
		log.Panic(err)
	}

	i.currentHash = block.PrevBlockHash

	return block
}
