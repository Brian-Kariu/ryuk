package db

import (
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type Config struct {
	Key   []byte
	Value []byte
}

func (c Config) ToBytes() (key []byte, value []byte) {
	return []byte(c.Key), []byte(c.Value)
}

type configs struct {
	config Config
}

type client struct {
	name         string
	globalBucket string
	db           *bolt.DB
}

func openDB(name string) (*bolt.DB, error) {
	db, err := bolt.Open(name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("Error opening DB: %v", err)
	}
	return db, nil
}

func (c client) String() string {
	return fmt.Sprintf("{name:%s, globalBucket:%s}", c.name, c.globalBucket)
}

func (c client) CreateBucket(name string) {
	dbError := c.db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return fmt.Errorf("Error creating bucket: %s", err)
		}
		log.Printf("New bucket %s created in db %s\n", name, c.name)
		return nil
	})
	if dbError != nil {
		fmt.Printf("createBucket: %s", dbError)
	}
	defer c.db.Close()
}

func (c client) AddKey(bucket string, data Config) error {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s does not exist", bucket)
		}

		return b.Put([]byte(data.Key), []byte(data.Value))
	})
	if err != nil {
		log.Printf("Error adding key to bucket %s: %v\n", bucket, err)
		return err
	}

	log.Printf("Added config: %s, to bucket: %s\n", data.Key, bucket)
	return nil
}

func (c client) GetKey(bucket string, config string) {
	v := []byte("")
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}
		v = b.Get([]byte(config))

		// HACK: This should change when a config can be json
		fmt.Printf("%s", string(v))
		return nil
	})
	if err != nil {
		log.Printf("Error retrieving key %s\n", config)
	}
}

func (c client) ListVars(bucket string) (map[string]string, error) {
	envVars := make(map[string]string)
	err := c.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		return b.ForEach(func(k, v []byte) error {
			envVars[string(k)] = string(v)
			return nil
		})
	})
	defer c.db.Close()
	return envVars, err
}

func (c client) DeleteKey(bucket, config string) {
	err := c.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return fmt.Errorf("bucket %s not found", bucket)
		}

		err := b.Delete([]byte(config))
		return err
	})
	if err != nil {
		log.Fatal("Delete operation failed: %s", err)
	}

	log.Printf("Config %s has been deleted", config)
}

func NewClient(name, globalBucket string) (*client, error) {
	if globalBucket == "" {
		globalBucket = "global_configs"
	}

	if name == "" {
		name = "ryuk"
	}
	db, err := openDB(name)
	if err != nil {
		return nil, err
	}

	clientInstance := &client{
		name:         name,
		globalBucket: globalBucket,
		db:           db,
	}
	return clientInstance, nil
}
