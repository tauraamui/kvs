package query_test

import (
	"testing"

	"github.com/tauraamui/kvs/v2"
	"github.com/tauraamui/kvs/v2/query"
	"github.com/tauraamui/kvs/v2/storage"
)

func BenchmarkQueryWithSingleFilterWithTwoRecords(b *testing.B) {
	db, err := kvs.NewMemKVDB()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	store.Save(kvs.RootOwner{}, &Balloon{Color: "WHITE", Size: 366})
	store.Save(kvs.RootOwner{}, &Balloon{Color: "RED", Size: 695})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("WHITE"))
	}
}

func BenchmarkQueryWithMultiFilterWithTwoRecords(b *testing.B) {
	db, err := kvs.NewMemKVDB()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	store.Save(kvs.RootOwner{}, &Balloon{Color: "WHITE", Size: 366})
	store.Save(kvs.RootOwner{}, &Balloon{Color: "RED", Size: 695})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("WHITE").Filter("size").Eq(306))
	}
}

func BenchmarkQueryWithMultiFilterWithFiveHunderedRecordsWithMatchingFilter(b *testing.B) {
	db, err := kvs.NewMemKVDB()
	if err != nil {
		b.Fatal(err)
	}
	defer db.Close()

	store := storage.New(db)
	defer store.Close()

	for i := 0; i < 500; i++ {
		color := "RED"
		if i%2 == 0 {
			color = "BLUE"
		}
		store.Save(kvs.RootOwner{}, &Balloon{Color: color, Size: i})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		query.Run[Balloon](store, kvs.RootOwner{}, query.New().Filter("color").Eq("RED").Filter("size").Eq(306, 422, 211))
	}
}
