package database

import (
	"context"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

const (
	bucketName = "users"
)

type DB struct {
	DB  *bolt.DB
	ctx context.Context
}

func NewDB(ctx context.Context) (*DB, error) {
	s := &DB{ctx: ctx}
	return s, s.connect()
}

func (s *DB) connect() (err error) {
	if s.DB, err = bolt.Open("./db_data/users.db", 0600, &bolt.Options{Timeout: 2 * time.Second}); err != nil {
		return
	}

	return s.DB.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(bucketName))
		return err
	})
}

func (s *DB) InsertUser(vk, zammad int) (err error) {
	err = s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		if err = b.Delete([]byte(strconv.Itoa(vk))); err != nil {
			return err
		}

		return b.Put([]byte(strconv.Itoa(vk)), []byte(strconv.Itoa(zammad)))
	})

	if err != nil {
		log.Error(err)
	}
	return
}

func (s *DB) DeleteUser(vk int) (err error) {
	err = s.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		return b.Delete([]byte(strconv.Itoa(vk)))
	})

	if err != nil {
		log.Error(err)
	}
	return
}

func (s *DB) SelectZammad(vk int) (zammad int, err error) {
	err = s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		if value := b.Get([]byte(strconv.Itoa(vk))); value != nil {
			zammad, err = strconv.Atoi(string(value))
			if err != nil {
				err = nil
				_ = b.Delete([]byte(strconv.Itoa(vk)))
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
	}
	return
}

func (s *DB) SelectVK(zammad int) (vk int, err error) {
	err = s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		c := b.Cursor()
		for key, value := c.First(); key != nil; key, value = c.Next() {
			if string(value) == strconv.Itoa(zammad) {
				vk, err = strconv.Atoi(string(key))
				if err != nil {
					err = nil
					_ = b.Delete(key)
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Error(err)
	}
	return
}
