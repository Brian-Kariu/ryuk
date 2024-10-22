package db

import (
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

func boltClient() *bolt.DB {
	// Open the ryuk.db data file in your current directory.
	// It will be created if it doesn't exist.
	dbName := "ryuk.db"
	db, err := bolt.Open(dbName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func bucketClient(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("global_configs"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func storeClient(db *bolt.DB) {
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("global_configs"))
		err := b.Put([]byte("log_level"), []byte("debug"))
		return err
	})
}
