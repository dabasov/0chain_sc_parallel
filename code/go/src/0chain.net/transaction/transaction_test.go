package transaction

import (
	"context"
	"fmt"
	"testing"
	"time"

	"0chain.net/client"
	"0chain.net/common"
	"0chain.net/datastore"
	"0chain.net/encryption"
)

var keyPairs = make(map[string]string)
var publicKeys = make([]string, 0, 1000)

func TestTransactionWrite(t *testing.T) {
	fmt.Printf("time : %v\n", time.Now().UnixNano()/int64(time.Millisecond))
	numClients := 1000
	createClients(numClients)
	start := time.Now()
	numTxns := 100000
	done := make(chan bool, numTxns)
	txnchannel := make(chan *Transaction, 10000)
	for i := 1; i <= 100; i++ {
		go processWorker(txnchannel, done)
	}
	for i := 1; i <= numTxns; i++ {
		publicKey := publicKeys[i%1000]
		pvtKey := keyPairs[publicKey]
		txnData := fmt.Sprintf("Txn(%v) Pay %v from %s\n", i, i%100, publicKey)
		go postTransaction(pvtKey, publicKey, txnData, txnchannel, done)
	}
	for count := 0; true; {
		<-done
		count++
		if count == numTxns {
			close(txnchannel)
			break
		}
	}
	fmt.Printf("Elapsed time for txns: %v\n", time.Since(start))
	time.Sleep(60 * time.Second)
}

func createClients(numClients int) {
	start := time.Now()
	fmt.Printf("Testing at %v\n", start)
	done := make(chan bool, numClients)
	for i := 1; i <= numClients; i++ {
		// This should be replaced with Key generation
		privateKey, publicKey := encryption.GenerateKeys()
		keyPairs[publicKey] = privateKey
		publicKeys = append(publicKeys, publicKey)
		go postClient(privateKey, publicKey, done)
	}
	for count := 0; true; {
		<-done
		count++
		if count == numClients {
			break
		}
	}
	fmt.Printf("Elapsed time for clients: %v\n", time.Since(start))
	time.Sleep(time.Second)
}

func postClient(privateKey string, publicKey string, done chan<- bool) {
	entity := client.ClientProvider()
	c, ok := entity.(*client.Client)
	if !ok {
		fmt.Printf("it's not ok!\n")
	}
	c.PublicKey = publicKey
	c.ID = encryption.Hash(publicKey)
	ctx := datastore.WithAsyncChannel(context.Background(), client.ClientEntityChannel)
	//ctx := datastore.WithConnection(context.Background())
	_, err := client.PutClient(ctx, entity)
	//datastore.GetCon(ctx).Close()
	if err != nil {
		fmt.Printf("error for %v : %v\n", publicKey, err)
	}
	done <- true
}

func postTransaction(privateKey string, publicKey string, txnData string, txnChannel chan<- *Transaction, done chan<- bool) {
	entity := TransactionProvider()
	t, ok := entity.(*Transaction)
	if !ok {
		fmt.Printf("it's not ok!\n")
	}
	t.ClientID = encryption.Hash(publicKey)
	t.TransactionData = txnData
	t.CreationDate = common.Now()
	t.Signature = encryption.Hash(txnData) //TODO: This should eventually be encryption.Sign(privateKey,txnData)
	t.Hash = encryption.Hash(t.Signature)
	txnChannel <- t
}

func processWorker(txnChannel <-chan *Transaction, done chan<- bool) {
	ctx := datastore.WithConnection(context.Background())
	defer datastore.GetCon(ctx).Close()

	for entity := range txnChannel {
		ctx = datastore.WithAsyncChannel(ctx, TransactionEntityChannel)
		_, err := PutTransaction(ctx, entity)
		if err != nil {
			fmt.Printf("error for %v : %v\n", entity, err)
		}
		done <- true
	}
}