package storage_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/kvs"
	"github.com/tauraamui/kvs/storage"
)

type Balloon struct {
	ID    uint32 `mdb:"ignore"`
	Color string
	Size  int
}

func (b Balloon) TableName() string { return "balloons" }
func (b *Balloon) SetID(id uint32)  { b.ID = id }
func (b *Balloon) Ref() any         { return b }

func TestStoreNewSuccess(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	is.NoErr(err)
	defer db.Close()

	store := storage.New(db)
	is.True(store != nil)
}

func TestStoreBalloonsSuccess(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	is.NoErr(err)
	defer db.Close()

	store := storage.New(db)
	is.True(store != nil)

	bigRedBalloon := Balloon{Color: "RED", Size: 695}
	smallYellowBalloon := Balloon{Color: "YELLOW", Size: 112}
	mediumWhiteBalloon := Balloon{Color: "WHITE", Size: 366}
	is.NoErr(store.Save(kvs.RootOwner{}, &bigRedBalloon))
	is.NoErr(store.Save(kvs.RootOwner{}, &smallYellowBalloon))
	is.NoErr(store.Save(kvs.RootOwner{}, &mediumWhiteBalloon))

	is.Equal(bigRedBalloon.ID, uint32(0))
	is.Equal(smallYellowBalloon.ID, uint32(1))
	is.Equal(mediumWhiteBalloon.ID, uint32(2))
}
