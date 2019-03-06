package sharder

import (
	"context"
	"sync"

	"0chain.net/core/cache"
	"0chain.net/core/ememorystore"

	"0chain.net/chaincore/block"
	"0chain.net/chaincore/chain"
	"0chain.net/chaincore/round"
	"0chain.net/core/datastore"
	"0chain.net/sharder/blockstore"
)

const (
	Normal  = 0
	Syncing = 1
	Accept  = 2
)

var sharderChain = &Chain{}

/*SetupSharderChain - setup the sharder's chain */
func SetupSharderChain(c *chain.Chain) {
	sharderChain.Chain = c
	sharderChain.BlockChannel = make(chan *block.Block, 128)
	sharderChain.RoundChannel = make(chan *round.Round, 128)
	blockCacheSize := 100
	sharderChain.BlockCache = cache.NewLRUCache(blockCacheSize)
	transactionCacheSize := int(c.BlockSize) * blockCacheSize
	sharderChain.BlockTxnCache = cache.NewLRUCache(transactionCacheSize)
	c.SetFetchedNotarizedBlockHandler(sharderChain)
	sharderChain.BSync = &BlockSync{
		//TODO configure acceptance tolerance value
		AcceptanceTolerance	: 	65,
		syncStatus			:   Normal,
		mutex				:   sync.RWMutex{},
	}
}

/*GetSharderChain - get the sharder's chain */
func GetSharderChain() *Chain {
	return sharderChain
}

/*Chain - A chain structure to manage the sharder activities */
type Chain struct {
	*chain.Chain
	BlockChannel  chan *block.Block
	RoundChannel  chan *round.Round
	BlockCache    cache.Cache
	BlockTxnCache cache.Cache
	SharderStats  Stats
	BSync         *BlockSync
}

/*BlockSync - A struct to track the block sync */
type BlockSync struct {
	AcceptanceTolerance int64
	acceptRound         int64
	syncRound           int64
	finalizeRound       int64
	syncStatus          int
	mutex               sync.RWMutex
}

/*GetStatus - get block sync status */
func (bs *BlockSync) GetStatus() int {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	return bs.syncStatus
}

/*SetStatus - set block sync status */
func (bs *BlockSync) SetStatus(status int) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	bs.syncStatus = status
}

/*GetFinalizationRound - get round to be finalized during block sync */
func (bs *BlockSync) GetFinalizationRound() int64 {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	return bs.finalizeRound
}

/*SetFinalizationRound - set round to be finalized during block sync */
func (bs *BlockSync) SetFinalizationRound(r int64) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	bs.finalizeRound = r
}

/*GetSyncingRound - get current syncing round */
func (bs *BlockSync) GetSyncingRound() int64 {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	return bs.syncRound
}

/*SetSyncingRund - set current syncing round */
func (bs *BlockSync) SetSyncingRound(r int64) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	bs.syncRound = r
}

/*GetAcceptanceRound - get current syncing round */
func (bs *BlockSync) GetAcceptanceRound() int64 {
	bs.mutex.RLock()
	defer bs.mutex.RUnlock()
	return bs.acceptRound
}

/*SetAcceptanceRound - set acceptance round during block sync */
func (bs *BlockSync) SetAcceptanceRound(r int64) {
	bs.mutex.Lock()
	defer bs.mutex.Unlock()
	bs.acceptRound = r
}

/*GetBlockChannel - get the block channel where the incoming blocks from the network are put into for further processing */
func (sc *Chain) GetBlockChannel() chan *block.Block {
	return sc.BlockChannel
}

/*GetRoundChannel - get the round channel where the finalized rounds are put into for further processing */
func (sc *Chain) GetRoundChannel() chan *round.Round {
	return sc.RoundChannel
}

/*SetupGenesisBlock - setup the genesis block for this chain */
func (sc *Chain) SetupGenesisBlock(hash string) *block.Block {
	gr, gb := sc.GenerateGenesisBlock(hash)
	if gr == nil || gb == nil {
		panic("Genesis round/block can not be null")
	}
	//sc.AddRound(gr)
	sc.AddGenesisBlock(gb)
	return gb
}

/*GetBlockFromStore - get the block from the store */
func (sc *Chain) GetBlockFromStore(blockHash string, round int64) (*block.Block, error) {
	bs := block.BlockSummary{Hash: blockHash, Round: round}
	return sc.GetBlockFromStoreBySummary(&bs)
}

/*GetBlockFromStoreBySummary - get the block from the store */
func (sc *Chain) GetBlockFromStoreBySummary(bs *block.BlockSummary) (*block.Block, error) {
	return blockstore.GetStore().ReadWithBlockSummary(bs)
}

/*GetRoundFromStore - get the round from a store*/
func (sc *Chain) GetRoundFromStore(ctx context.Context, roundNum int64) (*round.Round, error) {
	r := datastore.GetEntity("round").(*round.Round)
	r.Number = roundNum
	roundEntityMetadata := r.GetEntityMetadata()
	rctx := ememorystore.WithEntityConnection(ctx, roundEntityMetadata)
	defer ememorystore.Close(rctx)
	err := r.Read(rctx, r.GetKey())
	return r, err
}

/*GetBlockHash - get the block hash for a given round */
func (sc *Chain) GetBlockHash(ctx context.Context, roundNumber int64) (string, error) {
	r := sc.GetSharderRound(roundNumber)
	if r == nil {
		sr, err := sc.GetRoundFromStore(ctx, roundNumber)
		if err != nil {
			return "", err
		}
		r = sr
	}
	return r.BlockHash, nil
}

//GetSharderRound - get the sharder's version of the round
func (sc *Chain) GetSharderRound(roundNumber int64) *round.Round {
	r := sc.GetRound(roundNumber)
	if r == nil {
		return nil
	}
	sr, ok := r.(*round.Round)
	if !ok {
		return nil
	}
	return sr
}
