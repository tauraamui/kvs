package main

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/kvs"
)

func resolveGenericRepo() (ExampleRepo, error) {
	db, err := kvs.NewMemDB()
	if err != nil {
		return ExampleRepo{}, err
	}

	return ExampleRepo{DB: db}, nil
}

func TestSaveGeneric(t *testing.T) {
	is := is.New(t)

	r, err := resolveGenericRepo()
	is.NoErr(err)
	defer r.Close()

	fakeDataOne := ExampleData{
		Title: "Fake",
	}

	fakeDataTwo := ExampleData{
		Title: "Fakefake",
	}

	is.NoErr(r.Save(kvs.RootOwner{}, &fakeDataOne))
	is.NoErr(r.Save(kvs.RootOwner{}, &fakeDataTwo))
	is.Equal(fakeDataOne.ID, uint32(0))
	is.Equal(fakeDataTwo.ID, uint32(1))
}
