package kvs

import (
	"testing"

	"github.com/matryer/is"
)

func TestConvertToBytes(t *testing.T) {
	is := is.New(t)

	// Test input and expected output.
	tests := []struct {
		input    interface{}
		expected []byte
	}{
		{[]byte{1, 2, 3}, []byte{1, 2, 3}},
		{"hello", []byte("hello")},
		{struct{ A int }{5}, []byte("{\"A\":5}")},
	}

	// Iterate over the tests and compare the output of convertToBytes to the expected output.
	for _, test := range tests {
		result, err := convertToBytes(test.input)
		is.NoErr(err)
		is.Equal(result, test.expected)
	}
}

func TestConvertBytesFromBytes(t *testing.T) {
	is := is.New(t)

	input := []byte{1, 2, 3}
	var destination []byte
	err := convertFromBytes(input, &destination)
	is.NoErr(err)
	is.Equal(destination, input)

}

func TestConvertStringFromBytes(t *testing.T) {
	is := is.New(t)

	input := []byte("hello")
	var destination string
	err := convertFromBytes(input, &destination)
	is.NoErr(err)
	is.Equal(destination, string(input))
}

func TestConvertStructFromBytes(t *testing.T) {
	is := is.New(t)

	type TestStruct struct {
		A int
		B string
	}
	input := []byte("{\"A\":5,\"B\":\"hello\"}")
	var destination TestStruct
	err := convertFromBytes(input, &destination)
	is.NoErr(err)
	is.Equal(destination, TestStruct{A: 5, B: "hello"})
}
