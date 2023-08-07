# KVS (Key Value Store lib)

[![Maintainability](https://api.codeclimate.com/v1/badges/f3947361002e02193fdc/maintainability)](https://codeclimate.com/github/tauraamui/kvs/maintainability)
[![codecov](https://codecov.io/gh/tauraamui/kvs/branch/master/graph/badge.svg?token=UXP68F5SVG)](https://codecov.io/gh/tauraamui/kvs)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftauraamui%2Fkvs.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftauraamui%2Fkvs?ref=badge_shield)

## Usage

```go
type Balloon struct {
	ID    uint32 `mdb:"ignore"`
	Color string
	Size  int
}

func (b Balloon) TableName() string { return "balloons" }

func main() {
	conn, err := badger.Open(badger.DefaultOptions("example.db"))
	// or
	db, err := kvs.NewMemKVDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	store := storage.New(db) // or storage.New(conn)
	defer store.Close()

	store.Save(kvs.RootOwner{}, &Balloon{Color: "RED", Size: 695})
	store.Save(kvs.RootOwner{}, &Balloon{Color: "WHITE", Size: 366})

	bs, err := storage.LoadAll[Balloon](store, kvs.RootOwner{})
	for rowID, balloon := range bs {
		fmt.Printf("ROWID: %d, %+v\n", rowID, balloon)
	}
}
```

## Why?

Sometimes SQL (yes even SQLite) has too much overhead or added complexity, especially if you want to just save/persist structures
to storage and the most common relationship you have is one of "ownership". An alternative option is to use a key value DB,
but most of the time this is done by "pickling/encoding" the structure using some language native facility to store the structure,
and so you lose the ability to extract individual field entries if you want to.

Also, it's fast compared to similar libraries and in some cases SQLite: https://github.com/tauraamui/kvs-bench

## How?

When a structure is "saved", each struct field and it's value become a new entry. The library uses this key structure internally
"TABLE_NAME.COLUMN_NAME.OWNERUUID.ROWID" to store each of the struct's fields. The code in the example for example (hmm) results in:

```
key=balloons, value=
key=balloons.color.root.0, value=RED
key=balloons.color.root.1, value=WHITE
key=balloons.size.root.0, value=695
key=balloons.size.root.1, value=366
```

being stored. Using the power of key prefix iteration, you can then extract all structures which have a specific owner,
which the library does internally by using a prefix key which is the full key except for the row number element which is
a wildcard.



## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Ftauraamui%2Fkvs.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Ftauraamui%2Fkvs?ref=badge_large)
