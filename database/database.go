package database

import (
	"github.com/google/logger"
	"github.com/syndtr/goleveldb/leveldb"
)

var db *leveldb.DB

// Close closes database
func Close()  {
	err := db.Close()
	if err != nil {
		logger.Error(err)
	}
}

// Get retrieves value from database, returning an error if not found or an
// error is encountered
func Get(key []byte) ([]byte, error) {
	return db.Get(key, nil)
}

// Open opens an existing database, or creates one if it doesn't
// already exist
func Open() error {
	var err error
	db, err = leveldb.OpenFile("./database/data", nil)
	return err
}

// Put enters value into database, returning an error if encountered
func Put(key []byte, value []byte) error {
	return db.Put(key, value, nil)
}
