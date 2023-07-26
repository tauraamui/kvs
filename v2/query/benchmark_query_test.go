package query_test

import (
	"testing"

	"github.com/tauraamui/kvs/v2"
	"github.com/tauraamui/kvs/v2/query"
	"github.com/tauraamui/kvs/v2/storage"
)

func BenchmarkQueryWithSingleFilter(b *testing.B) {
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
		query.Run[Balloon](store, query.New().Filter("color").Eq("WHITE"))
	}
}
