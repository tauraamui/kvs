package kvs

import (
	"fmt"
	"io"
	"os"

	"github.com/dgraph-io/badger/v3"
)

type KVDB struct {
	conn *badger.DB
}

func NewKVDB(db *badger.DB) (KVDB, error) {
	return newKVDB(db)
}

func NewMemKVDB() (KVDB, error) {
	return newKVDB(nil)
}

func newKVDB(db *badger.DB) (KVDB, error) {
	if db == nil {
		db, err := badger.Open(badger.DefaultOptions("").WithLogger(nil).WithInMemory(true))
		if err != nil {
			return KVDB{}, err
		}
		return KVDB{conn: db}, nil
	}

	return KVDB{conn: db}, nil
}

func (db KVDB) GetSeq(key []byte, bandwidth uint64) (*badger.Sequence, error) {
	return db.conn.GetSequence(key, bandwidth)
}

func (db KVDB) View(f func(txn *badger.Txn) error) error {
	return db.conn.View(f)
}

func (db KVDB) Update(f func(txn *badger.Txn) error) error {
	return db.conn.Update(f)
}

func (db KVDB) DumpTo(w io.Writer) error {
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

func (db KVDB) DumpToStdout() error {
	return db.DumpTo(os.Stdout)
}

func (db KVDB) Close() error {
	return db.conn.Close()
}
