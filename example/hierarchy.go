package main

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/tauraamui/kvs"
	"github.com/tauraamui/kvs/storage"
)

type SmallChild struct {
	UUID         kvs.UUID
	HungryMetric uint32
	Norished     bool
}

type Cake struct {
	ID       uint32 `mdb:"ignore"`
	Type     string
	Calories int
}

func (b Cake) TableName() string { return "cakes" }

func hierarchy() {
	db, err := kvs.NewMemKVDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	child := SmallChild{
		UUID:         uuid.New(),
		HungryMetric: 100,
	}

	disguistingVeganCake := Cake{Type: "INEDIBLE", Calories: -38}
	healthyishCarrotCake := Cake{Type: "CARROT", Calories: 280}
	redVelvetCake := Cake{Type: "RED_VELVET", Calories: 410}

	store.Save(child.UUID, &disguistingVeganCake)
	store.Save(child.UUID, &healthyishCarrotCake)
	store.Save(child.UUID, &redVelvetCake)

	bs, err := storage.LoadAll(store, Cake{}, child.UUID)
	for _, cake := range bs {
		fmt.Printf("ROWID: %d, %+v\n", cake.ID, cake)
	}
}
