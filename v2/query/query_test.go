package query_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/kvs/v2"
	"github.com/tauraamui/kvs/v2/query"
	"github.com/tauraamui/kvs/v2/storage"
)

type Balloon struct {
	ID    uint32 `mdb:"ignore"`
	Color string
	Size  int
}

func (b Balloon) TableName() string { return "balloons" }

func TestQueryFilters(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	store.Save(kvs.RootOwner{}, &Balloon{Color: "RED", Size: 695})
	store.Save(kvs.RootOwner{}, &Balloon{Color: "WHITE", Size: 366})

	bs, err := query.Run[Balloon](store, query.New().Filter("color").Eq("WHITE").Filter("size").Eq(695))
	is.NoErr(err)

	is.Equal(len(bs), 1)
}
