package kvs

import (
	"fmt"
	"io"
	"os"

	"github.com/dgraph-io/badger/v3"
)

type DB struct {
	conn *badger.DB
}

func NewDB(db *badger.DB) (DB, error) {
	return newDB(false)
}

func NewMemDB() (DB, error) {
	return newDB(true)
}

func newDB(inMemory bool) (DB, error) {
	db, err := badger.Open(badger.DefaultOptions("").WithLogger(nil).WithInMemory(inMemory))
	if err != nil {
		return DB{}, err
	}

	return DB{conn: db}, nil
}

func (db DB) GetSeq(key []byte, bandwidth uint64) (*badger.Sequence, error) {
	return db.conn.GetSequence(key, bandwidth)
}

func (db DB) View(f func(txn *badger.Txn) error) error {
	return db.conn.View(f)
}

func (db DB) Update(f func(txn *badger.Txn) error) error {
	return db.conn.Update(f)
}

func (db DB) DumpTo(w io.Writer) error {
	return db.conn.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Fprintf(w, "key=%s, value=%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (db DB) DumpToStdout() error {
	return db.DumpTo(os.Stdout)
}

func (db DB) Close() error {
	return db.conn.Close()
}
