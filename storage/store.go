package storage

import (
	"strconv"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/kvs"
)

type Value interface {
	TableName() string
}

type Store struct {
	db  kvs.KVDB
	pks map[string]*badger.Sequence
}

func New(db kvs.KVDB) Store {
	return Store{db: db, pks: map[string]*badger.Sequence{}}
}

func (s Store) Save(owner kvs.UUID, value Value) error {
	rowID, err := nextRowID(s.db, value.TableName(), s.pks)
	if err != nil {
		return err
	}

	return saveValue(s.db, value.TableName(), owner, rowID, value)
}

func (s Store) Update(owner kvs.UUID, value Value, rowID uint32) error {
	return saveValue(s.db, value.TableName(), owner, rowID, value)
}

func saveValue(db kvs.KVDB, tableName string, ownerID kvs.UUID, rowID uint32, v Value) error {
	if v == nil {
		return nil
	}
	entries := kvs.ConvertToEntries(tableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := kvs.Store(db, e); err != nil {
			return err
		}
	}

	return kvs.LoadID(v, rowID)
}

func (s Store) Delete(owner kvs.UUID, value Value, rowID uint32) error {
	db := s.db

	blankEntries := kvs.ConvertToBlankEntries(value.TableName(), owner, rowID, value)
	for _, ent := range blankEntries {
		db.Update(func(txn *badger.Txn) error {
			return txn.Delete(ent.Key())
		})
	}

	return nil
}

func Load[T Value](s Store, dest T, owner kvs.UUID, rowID uint32) error {
	db := s.db

	blankEntries := kvs.ConvertToBlankEntries(dest.TableName(), owner, rowID, dest)
	for _, ent := range blankEntries {
		db.View(func(txn *badger.Txn) error {
			item, err := txn.Get(ent.Key())
			if err != nil {
				return err
			}

			if err := item.Value(func(val []byte) error {
				ent.Data = val
				return nil
			}); err != nil {
				return err
			}

			return nil
		})

		ent.RowID = rowID

		if err := kvs.LoadEntry(dest, ent); err != nil {
			return err
		}
	}

	return kvs.LoadID(dest, rowID)
}

func extractRowFromKey(k string) (int, error) {
	rowPos := strings.LastIndex(k, ".")
	return strconv.Atoi(k[rowPos+1:])
}

func LoadAll[T Value](s Store, v T, owner kvs.UUID) ([]T, error) {
	db := s.db
	dest := []T{}

	// keep for later reference
	/*
		typeRef := new(E)
	*/

	blankEntries := kvs.ConvertToBlankEntries(v.TableName(), owner, 0, v)
	for _, ent := range blankEntries {
		// iterate over all stored values for this entry
		prefix := ent.PrefixKey()
		if err := db.View(func(txn *badger.Txn) error {
			it := txn.NewIterator(badger.DefaultIteratorOptions)
			defer it.Close()

			var structFieldIndex uint32 = 0
			for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
				item := it.Item()

				if len(dest) == 0 || structFieldIndex >= uint32(len(dest)) {
					key := string(item.Key())
					rowID, err := extractRowFromKey(string(key))
					if err != nil {
						return err
					}

					dest = append(dest, *new(T))

					// for reasons, we have to just keep assigning the current "field we're on" as the full entry's ID
					if err := kvs.LoadID(&dest[structFieldIndex], uint32(rowID)); err != nil {
						return err
					}
				}

				ent.RowID = structFieldIndex
				if err := item.Value(func(val []byte) error {
					ent.Data = val
					return nil
				}); err != nil {
					return err
				}

				if err := kvs.LoadEntry(&dest[structFieldIndex], ent); err != nil {
					return err
				}
				structFieldIndex++
			}
			return nil
		}); err != nil {
			return nil, err
		}
	}
	return dest, nil
}

func (s Store) Close() (err error) {
	if s.pks == nil {
		return
	}
	for _, seq := range s.pks {
		if seq == nil {
			continue
		}
		// TODO:(tauraamui) should collect all errors into group here rather than immediately escape
		if err = seq.Release(); err != nil {
			return
		}
	}

	s.pks = nil

	return
}

func nextRowID(db kvs.KVDB, tableName string, pks map[string]*badger.Sequence) (uint32, error) {
	seq, err := resolveSequence(db, tableName, pks)
	if err != nil {
		return 0, err
	}

	s, err := seq.Next()
	if err != nil {
		return 0, err
	}
	return uint32(s), nil
}

func nextSequence(seq *badger.Sequence) (uint32, error) {
	s, err := seq.Next()
	if err != nil {
		return 0, err
	}
	return uint32(s), nil
}

func resolveSequence(db kvs.KVDB, tableName string, pks map[string]*badger.Sequence) (*badger.Sequence, error) {
	seq, ok := pks[tableName]
	var err error
	if !ok {
		seq, err = db.GetSeq([]byte(tableName), 1)
		if err != nil {
			return nil, err
		}
		pks[tableName] = seq
	}

	return seq, nil
}
