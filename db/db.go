package db

import (
	"encoding/json"
	"errors"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/boltdb/bolt"
)

var bdb *bolt.DB

const rrBucket = "rr"

// Record store a DNS record with metadata
type Record struct {
	RR      string
	Expires int64
}

//Connect create a bucket
func Connect(dbPath string) error {
	db, err := bolt.Open(dbPath, 0600, &bolt.Options{
		Timeout: 10 * time.Second,
	})
	if err != nil {
		return err
	}

	bdb = db
	// Create dns bucket if doesn't exist
	createBucket(rrBucket)

	return nil
}

//Disconnect create a bucket
func Disconnect() error {
	if bdb == nil {
		return nil
	}
	return bdb.Close()
}

//createBucket create bucket if not exists
func createBucket(bucket string) error {
	return bdb.Update(func(tx *bolt.Tx) error {
		_, err1 := tx.CreateBucketIfNotExists([]byte(bucket))
		if err1 != nil {
			e := errors.New("Create bucket:  " + bucket)
			log.Println(e.Error())

			return e
		}
		return nil
	})
}

//DeleteRecord create a bucket
func DeleteRecord(key string) (err error) {
	err = bdb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(rrBucket))
		err1 := b.Delete([]byte(key))

		if err1 != nil {
			e := errors.New("Delete record failed for domain:  " + key)
			log.Println(e.Error())
			return e
		}

		return nil
	})

	return err
}

//StoreRecord save a new record
func StoreRecord(key string, record Record) error {
	return bdb.Update(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(rrBucket))

		val, err := json.Marshal(record)
		if err != nil {
			return err
		}

		err1 := b.Put([]byte(key), val)
		if err1 != nil {
			return err1
		}

		return nil
	})
}

//GetRecord return a stored record for a domain
func GetRecord(key string) (r Record, err error) {
	err = bdb.View(func(tx *bolt.Tx) error {

		b := tx.Bucket([]byte(rrBucket))
		raw := b.Get([]byte(key))

		if len(raw) == 0 {
			e := errors.New("Record not found, key:  " + key)
			log.Println(e.Error())
			return e
		}

		e := json.Unmarshal(raw, r)
		if e != nil {
			return e
		}

		return nil
	})

	return r, err
}
