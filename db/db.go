package database

import (
	"context"
	"errors"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	bolt "go.etcd.io/bbolt"
)

const (
	bucketName = "zammad_vk"
)

type DB struct {
	DB  *bolt.DB
	ctx context.Context
}

func NewDB(ctx context.Context) (s *DB, err error) {
	s = &DB{ctx: ctx}
	err = s.connect()
	return
}

func (s *DB) connect() (err error) {
	s.DB, err = bolt.Open("users.db", 0600, &bolt.Options{Timeout: 2 * time.Second})
	if err != nil {
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
			return errors.New("bucket not found")
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
			return errors.New("bucket not found")
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
			return errors.New("bucket not found")
		}

		if v := b.Get([]byte(strconv.Itoa(vk))); v != nil {
			vk, _ = strconv.Atoi(string(v))
		}

		return nil
	})

	if err != nil && !errors.Is(err, bolt.ErrBucketNotFound) {
		log.Error(err)
	} else {
		err = nil
	}

	return
}

func (s *DB) SelectVK(zammad int) (vk int, err error) {
	err = s.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b == nil {
			return errors.New("bucket not found")
		}

		c := b.Cursor()
		for key, value := c.First(); key != nil; key, value = c.Next() {
			if string(value) == strconv.Itoa(zammad) {
				vk, _ = strconv.Atoi(string(key))
				return nil
			}
		}

		return nil
	})

	if err != nil && !errors.Is(err, bolt.ErrBucketNotFound) {
		log.Error(err)
	} else {
		err = nil
	}

	return
}
