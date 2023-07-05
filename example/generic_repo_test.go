package example_test

import (
	"testing"

	"github.com/matryer/is"
	"github.com/tauraamui/kvs"
	repo "github.com/tauraamui/kvs/example"
)

type exampleModel struct {
	ID   uint32 `mdb:"ignore"`
	UID  uint32
	Name string
}

func (m *exampleModel) SetID(id uint32)  { m.ID = id }
func (m *exampleModel) Ref() interface{} { return m }

func resolveGenericRepo() (repo.GenericRepo, error) {
	db, err := kvs.NewMemDB()
	if err != nil {
		return repo.GenericRepo{}, err
	}

	return repo.GenericRepo{TableName: "mailboxes", DB: db}, nil
}

func TestSaveGeneric(t *testing.T) {
	t.Skip("pending migration to UUID")
	is := is.New(t)

	r, err := resolveGenericRepo()
	is.NoErr(err)
	defer r.Close()

	example := exampleModel{
		UID:  83,
		Name: "Fake",
	}

	is.NoErr(r.Save(0, &example))
}
