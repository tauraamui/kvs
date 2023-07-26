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

func TestQueryFilterWithSinglePredicateSuccess(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	store.Save(kvs.RootOwner{}, &Balloon{Color: "WHITE", Size: 366})
	store.Save(kvs.RootOwner{}, &Balloon{Color: "RED", Size: 695})

	bs, err := query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("WHITE"))
	is.NoErr(err)
	is.Equal(len(bs), 1)
	is = is.NewRelaxed(t)
	is.Equal(bs[0].Color, "WHITE")
	is.Equal(bs[0].Size, 366)

	is = is.New(t)

	bs, err = query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("RED"))
	is.NoErr(err)
	is.Equal(len(bs), 1)
	is = is.NewRelaxed(t)
	is.Equal(bs[0].Color, "RED")
	is.Equal(bs[0].Size, 695)
}

func TestQueryFilterWithSinglePredicateMultipleValuesSuccess(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	store.Save(kvs.RootOwner{}, &Balloon{Color: "WHITE", Size: 366})
	store.Save(kvs.RootOwner{}, &Balloon{Color: "RED", Size: 695})

	bs, err := query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("GREEN", "WHITE", "CYAN", "PURPLE", "RED", "GOLD"))
	is.NoErr(err)
	is.Equal(len(bs), 2)
	is = is.NewRelaxed(t)
	is.Equal(bs[0].Color, "WHITE")
	is.Equal(bs[0].Size, 366)

	is.Equal(bs[1].Color, "RED")
	is.Equal(bs[1].Size, 695)
}

func TestQueryFilterWithSinglePredicateFailure(t *testing.T) {
	is := is.New(t)

	db, err := kvs.NewMemKVDB()
	if err != nil {
		is.NoErr(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	store.Save(kvs.RootOwner{}, &Balloon{Color: "WHITE", Size: 366})
	store.Save(kvs.RootOwner{}, &Balloon{Color: "RED", Size: 695})

	bs, err := query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("deef"))
	is.NoErr(err)
	is.Equal(len(bs), 0)

	bs, err = query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("rgrr"))
	is.NoErr(err)
	is.Equal(len(bs), 0)
}

func TestQueryFilterWithMultiplePredicateSuccess(t *testing.T) {
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

	bs, err := query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("WHITE").Filter("size").Eq(366))
	is.NoErr(err)
	is.Equal(len(bs), 1)
	is = is.NewRelaxed(t)
	is.Equal(bs[0].Color, "WHITE")
	is.Equal(bs[0].Size, 366)

	is = is.New(t)

	bs, err = query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("RED").Filter("size").Eq(695))
	is.NoErr(err)
	is.Equal(len(bs), 1)
	is = is.NewRelaxed(t)
	is.Equal(bs[0].Color, "RED")
	is.Equal(bs[0].Size, 695)
}

func TestQueryFilterWithMultiplePredicateMultipleValuesSuccess(t *testing.T) {
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

	bs, err := query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("ZIMA_BLUE", "BLACK", "WHITE", "RED").Filter("size").Eq(222, 366, 948))
	is.NoErr(err)
	is.Equal(len(bs), 1)
	is = is.NewRelaxed(t)
	is.Equal(bs[0].Color, "WHITE")
	is.Equal(bs[0].Size, 366)
}

func TestQueryFilterWithMultiplePredicateFailure(t *testing.T) {
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

	bs, err := query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("WHITE").Filter("size").Eq(110))
	is.NoErr(err)
	is.Equal(len(bs), 0)

	bs, err = query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("RED").Filter("size").Eq(548))
	is.NoErr(err)
	is.Equal(len(bs), 0)
}

func TestQueryFilterWithMultiplePredicateMultipleValuesFailure(t *testing.T) {
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

	bs, err := query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("GREY", "PINK").Filter("size").Eq(110, 366))
	is.NoErr(err)
	is.Equal(len(bs), 0)
}
