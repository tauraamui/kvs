package example

import (
	"github.com/dgraph-io/badger/v3"
	"github.com/tauraamui/kvs"
)

type Value interface {
	SetID(id uint32)
	Ref() interface{}
}

type GenericRepo struct {
	TableName string
	DB        kvs.DB
	seq       *badger.Sequence
}

func saveValueWithUUID(db kvs.DB, tableName string, ownerID kvs.UUID, rowID uint32, v Value) error {
	entries := kvs.ConvertToEntriesWithUUID(tableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := kvs.Store(db, e); err != nil {
			return err
		}
	}

	v.SetID(rowID)
	return nil
}

func saveValue(db kvs.DB, tableName string, rowID, ownerID uint32, v Value) error {
	entries := kvs.ConvertToEntries(tableName, ownerID, rowID, v)
	for _, e := range entries {
		if err := kvs.Store(db, e); err != nil {
			return err
		}
	}

	v.SetID(rowID)
	return nil
}

func (r *GenericRepo) Save(ownerID uint32, v Value) error {
	rowID, err := r.nextRowID()
	if err != nil {
		return err
	}

	return saveValue(r.DB, r.TableName, rowID, ownerID, v)
}

func (r *GenericRepo) nextRowID() (uint32, error) {
	if r.seq == nil {
		seq, err := r.DB.GetSeq([]byte(r.TableName), 100)
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

func (r *GenericRepo) Close() {
	if r.seq == nil {
		return
	}
	r.seq.Release()
}
