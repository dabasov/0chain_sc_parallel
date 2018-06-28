package transaction

import (
	"context"

	"0chain.net/datastore"
	"0chain.net/util"
)

/*Confirmation - a data structure that provides the confirmation that a transaction is included into the block chain */
type Confirmation struct {
	Version   string `json:"version"`
	Hash      string `json:"hash"`
	BlockHash string `json:"block_hash"`
	datastore.CreationDateField
	Round           int64             `json:"round"`
	RoundRandomSeed int64             `json:"round_random_seed"`
	MerkleTreeRoot  string            `json:"merkle_tree_root"`
	MerkleTreePath  []util.MTPathNode `json:"merkle_tree_path"`
}

var transactionConfirmationEntityMetadata *datastore.EntityMetadataImpl

/*GetEntityMetadata - implementing the interface */
func (c *Confirmation) GetEntityMetadata() datastore.EntityMetadata {
	return transactionConfirmationEntityMetadata
}

/*SetKey - implement interface */
func (c *Confirmation) SetKey(key datastore.Key) {
	c.Hash = datastore.ToString(key)
}

/*GetKey - implement interface */
func (c *Confirmation) GetKey() datastore.Key {
	return datastore.ToKey(c.Hash)
}

/*ComputeProperties - implement interface */
func (c *Confirmation) ComputeProperties() {

}

func (c *Confirmation) Validate(ctx context.Context) error {
	return nil
}

/*Read - store read */
func (c *Confirmation) Read(ctx context.Context, key datastore.Key) error {
	return c.GetEntityMetadata().GetStore().Read(ctx, key, c)
}

/*Write - store read */
func (c *Confirmation) Write(ctx context.Context) error {
	return c.GetEntityMetadata().GetStore().Write(ctx, c)
}

/*Delete - store read */
func (c *Confirmation) Delete(ctx context.Context) error {
	return c.GetEntityMetadata().GetStore().Delete(ctx, c)
}

/*GetHash - hashable implementation */
func (c *Confirmation) GetHash() string {
	return c.Hash
}

func TransactionConfirmationProvider() datastore.Entity {
	t := &Confirmation{}
	t.Version = "1.0"
	return t
}

func SetupTxnConfirmationEntity(store datastore.Store) {
	transactionConfirmationEntityMetadata = datastore.MetadataProvider()
	transactionConfirmationEntityMetadata.Name = "txn_confirmation"
	transactionConfirmationEntityMetadata.Provider = TransactionConfirmationProvider
	transactionConfirmationEntityMetadata.Store = store
	datastore.RegisterEntityMetadata("txn_confirmation", transactionConfirmationEntityMetadata)
}