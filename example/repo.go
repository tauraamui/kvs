package main

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/kvs"
)

const (
	exampleTableName = "example"
)

type ExampleRepo struct {
	DB  kvs.DB
	seq *badger.Sequence
}

func (r ExampleRepo) Save(owner kvs.UUID, val Value) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValue(r.DB, r.tableName(), owner, rowID, val)
}

func (r ExampleRepo) FetchByOwner(owner kvs.UUID) ([]ExampleData, error) {
	return fetchByOwner[ExampleData](r.DB, r.tableName(), owner)
}

func (r ExampleRepo) tableName() string {
	return exampleTableName
}

func (r ExampleRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(exampleTableName), 1)
		if err != nil {
			return 0, err
		}
		r.seq = seq
	}

	s, err := r.seq.Next()
	if err != nil {
		return 0, err
	}
	return uint32(s), nil
}

func (r ExampleRepo) Close() error {
	if r.seq == nil {
		return nil
	}
	r.seq.Release()
	return nil
}

type Value interface {
	SetID(id uint32)
	Ref() interface{}
}

func saveValue(db kvs.DB, tableName string, ownerID kvs.UUID, rowID uint32, v Value) error {
	if v == nil {
		return nil
	}
	entries := kvs.ConvertToEntriesWithUUID(tableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := kvs.Store(db, e); err != nil {
			return err
		}
	}

	v.SetID(rowID)

	return nil
}

func fetchByOwner[E any](db kvs.DB, tableName string, owner kvs.UUID) ([]E, error) {
	dest := []E{}

	typeRef := new(E)

	blankEntries := kvs.ConvertToBlankEntriesWithUUID(tableName, owner, 0, typeRef)
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var rows uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				if len(dest) == 0 || rows >= uint32(len(dest)) {
					dest = append(dest, *new(E))
				}
				item := it.Item()
				ent.RowID = rows
				if err := item.Value(func(val []byte) error {
					ent.Data = val
					return nil
				}); err != nil {
					return err
				}

				if err := kvs.LoadEntry(&dest[rows], ent); err != nil {
					return err
				}
				rows++
			}
			return nil
		})
	}
	return dest, nil
}
