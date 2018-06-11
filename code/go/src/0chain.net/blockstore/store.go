package blockstore

import (
	"0chain.net/block"
)

/*BlockStore - an interface to read and write blocks to some storage */
type BlockStore interface {
	Write(b *block.Block) error
	Read(hash string, round int64) (*block.Block, error)
}

var Store BlockStore

/*GetStore - get the block store that's is setup */
func GetStore() BlockStore {
	return Store
}