package main

type ExampleData struct {
	ID    uint32 `mdb:"ignore"`
	Title string
}

func (d *ExampleData) SetID(id uint32) { d.ID = id }
func (d *ExampleData) Ref() any        { return d }
