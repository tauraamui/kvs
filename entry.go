package kvs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/dgraph-io/badger/v3"
	"github.com/google/uuid"
)

type Entry struct {
	TableName  string
	ColumnName string
	OwnerUUID  UUID
	RowID      uint32
	Data       []byte
}

func (e Entry) PrefixKey() []byte {
	return []byte(fmt.Sprintf("%s.%s.%s", e.TableName, e.ColumnName, e.resolveOwnerID()))
}

func (e Entry) Key() []byte {
	return []byte(fmt.Sprintf("%s.%s.%s.%d", e.TableName, e.ColumnName, e.resolveOwnerID(), e.RowID))
}

func (e Entry) resolveOwnerID() string {
	if e.OwnerUUID == nil {
		e.OwnerUUID = RootOwner{}
	}
	return e.OwnerUUID.String()
}

func Store(db KVDB, e Entry) error {
	return db.conn.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(e.Key()), e.Data)
	})
}

func Get(db KVDB, e *Entry) error {
	return db.conn.View(func(txn *badger.Txn) error {
		lookupKey := e.Key()
		item, err := txn.Get(lookupKey)
		if err != nil {
			return fmt.Errorf("%s: %s", strings.ToLower(err.Error()), lookupKey)
		}

		if err := item.Value(func(val []byte) error {
			e.Data = val
			return nil
		}); err != nil {
			return err
		}

		return nil
	})
}

func ConvertToBlankEntries(tableName string, ownerID UUID, rowID uint32, x any) []Entry {
	v := reflect.ValueOf(x)
	return convertToEntries(tableName, ownerID, rowID, v, false)
}

func ConvertToEntries(tableName string, ownerID UUID, rowID uint32, x interface{}) []Entry {
	v := reflect.ValueOf(x)
	return convertToEntries(tableName, ownerID, rowID, v, true)
}

type UUID interface {
	String() string
}

type RootOwner struct{}

func (o RootOwner) String() string { return "root" }

func LoadEntry(s interface{}, entry Entry) error {
	// convert the interface value to a reflect.Value so we can access its fields
	val := reflect.ValueOf(s).Elem()

	field, err := resolveFieldRef(val, entry.ColumnName)
	if err != nil {
		return err
	}

	// convert the entry's Data field to the type of the target field
	if err := convertFromBytes(entry.Data, field.Addr().Interface()); err != nil {
		return fmt.Errorf("failed to convert entry data to field type: %v", err)
	}

	return nil
}

func resolveFieldRef(v reflect.Value, nameToMatch string) (reflect.Value, error) {
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if strings.EqualFold(field.Name, nameToMatch) {
			return v.Field(i), nil
		}
	}

	return reflect.Zero(reflect.TypeOf(v)), fmt.Errorf("struct does not have a field with name %q", nameToMatch)
}

func LoadEntries(s interface{}, entries []Entry) error {
	for _, entry := range entries {
		if err := LoadEntry(s, entry); err != nil {
			return err
		}
	}

	return nil
}

func convertToEntries(tableName string, ownerUUID UUID, rowID uint32, v reflect.Value, includeData bool) []Entry {
	entries := []Entry{}

	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		vv := reflect.Indirect(v)
		f := vv.Type().Field(i)

		fOpts := resolveFieldOptions(f)
		if fOpts.Ignore {
			continue
		}

		e := Entry{
			TableName:  tableName,
			ColumnName: strings.ToLower(f.Name),
			OwnerUUID:  ownerUUID,
			RowID:      rowID,
		}

		if includeData {
			bd, err := convertToBytes(v.Field(i).Interface())
			if err != nil {
				return entries
			}
			e.Data = bd
		}

		entries = append(entries, e)
	}

	return entries
}

func convertToBytes(i interface{}) ([]byte, error) {
	// Check the type of the interface.
	switch v := i.(type) {
	case []byte:
		// Return the input as a []byte if it is already a []byte.
		return v, nil
	case string:
		// Convert the string to a []byte and return it.
		return []byte(v), nil
	default:
		// Use json.Marshal to convert the interface to a []byte.
		return json.Marshal(v)
	}
}

func convertFromBytes(data []byte, i interface{}) error {
	// Check that the destination argument is a pointer.
	if reflect.TypeOf(i).Kind() != reflect.Ptr {
		return fmt.Errorf("destination must be a pointer")
	}

	// Check the type of the interface.
	switch v := i.(type) {
	case *[]byte:
		// Set the value of the interface to the []byte if it is a pointer to a []byte.
		*v = data
		return nil
	case *string:
		// Convert the []byte to a string and set the value of the interface to the string.
		*v = string(data)
		return nil
	case *UUID:
		// Convert the []byte to a UUID instance and set the value of the interface to it.
		uuidv, err := uuid.ParseBytes(data)
		if err != nil {
			return err
		}
		*v = uuidv
		return nil
	default:
		// Use json.Unmarshal to convert the []byte to the interface.
		return json.Unmarshal(data, v)
	}
}

type mdbFieldOptions struct {
	Ignore bool
}

func resolveFieldOptions(f reflect.StructField) mdbFieldOptions {
	mdbTagValue := f.Tag.Get("mdb")
	return mdbFieldOptions{
		Ignore: strings.Contains(mdbTagValue, "ignore"),
	}
}
