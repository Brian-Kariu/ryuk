package db

import (
	"fmt"
	"log"
	"time"

	bolt "go.etcd.io/bbolt"
)

type config struct {
	key   []byte
	value []byte
}

func (c config) ToBytes() (key []byte, value []byte) {
	return []byte(c.key), []byte(c.value)
}

type configs struct {
	config config
}

type client struct {
	name         string
	globalBucket string
}

func (c client) String() string {
	return fmt.Sprintf("{name:%s, globalBucket:%s}", c.name, c.globalBucket)
}

func (c *client) openDB() (*bolt.DB, error) {
	db, err := bolt.Open(c.name, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("Error opening DB: %v", err)
	}
	return db, nil
}

func (c client) init() {
	db, err := c.openDB()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	log.Printf("Ryuk %s initialized.\n", c.name)
}

func (c client) CreateBucket(name string) {
	db, err := c.openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	dbError := db.Update(func(tx *bolt.Tx) error {
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
}

func (c client) addKey(bucket string, data config) {
	db, err := c.openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		key, value := data.ToBytes()
		err := b.Put(key, value)
		return fmt.Errorf("Error adding key to bucket %s: %s", bucket, err)
	})
	log.Printf("Added config: %s\n", data.key)
}

func (c client) getKey(bucket string, config string) string {
	db, err := c.openDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	v := []byte("")

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		v = b.Get([]byte(config))
		if v == nil {
			log.Printf("Config %s does not exist or is nested.\n", config)
		} else {
			log.Printf("Found config %s.\n", config)
		}
		return nil
	})

	// HACK: This should change when a config can be json
	return string(v)
}

func NewClient(name, globalBucket string) *client {
	if globalBucket == "" {
		globalBucket = "global_configs"
	}

	if name == "" {
		name = "ryuk"
	}
	clientInstance := &client{
		name:         name + ".db",
		globalBucket: globalBucket,
	}
	clientInstance.init()
	return clientInstance
}
