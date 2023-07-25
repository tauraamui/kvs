package main

import (
	"fmt"

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

func simple() {
	db, err := kvs.NewMemKVDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	store.Save(kvs.RootOwner{}, &Balloon{Color: "RED", Size: 695})
	store.Save(kvs.RootOwner{}, &Balloon{Color: "WHITE", Size: 366})

	q := query.New().Filter("color").Eq("red")
	bs, err := query.Run[Balloon](store, q)

	//bs, err := storage.LoadAll(store, Balloon{}, kvs.RootOwner{})
	for _, balloon := range bs {
		fmt.Printf("ROWID: %d, %+v\n", balloon.ID, balloon)
	}

}
