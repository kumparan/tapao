package tapao

import (
	"encoding/json"
	"testing"

	"github.com/vmihailenco/msgpack"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Int64   int64   `json:"int64"`
	Float64 float64 `json:"float64"`
	Bool    bool    `json:"bool"`
	String  string  `json:"string"`
}

var testStruct = TestStruct{
	Int64:   222,
	Float64: 333.333,
	Bool:    true,
	String:  "ferdian the best",
}

func TestMarshal(t *testing.T) {
	in := testStruct
	msgpackOut, err := msgpack.Marshal(in)
	assert.NoError(t, err)
	jsonOut, err := json.Marshal(in)
	assert.NoError(t, err)

	t.Run("default", func(t *testing.T) {
		out, err := Marshal(in)
		assert.NoError(t, err)
		assert.Equal(t, msgpackOut, out)
	})

	t.Run("use json", func(t *testing.T) {
		out, err := Marshal(in, Use(JSON))
		assert.NoError(t, err)
		assert.Equal(t, jsonOut, out)
	})

	t.Run("use msgpack", func(t *testing.T) {
		out, err := Marshal(in, Use(MessagePack))
		assert.NoError(t, err)
		assert.Equal(t, jsonOut, out)
	})

	// TODO: test fallback, skip for now
}

func TestUnmarshal(t *testing.T) {
	source := testStruct
	msgpackIn, err := msgpack.Marshal(source)
	assert.NoError(t, err)
	jsonIn, err := json.Marshal(source)
	assert.NoError(t, err)

	t.Run("default", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(msgpackIn, &out)
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("use json", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(jsonIn, &out, Use(JSON))
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("use msgpack", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(msgpackIn, &out, Use(MessagePack))
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("use json, fallback msgpack", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(msgpackIn, &out, Use(JSON), FallbackWith(MessagePack))
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("use msgpack, fallback json", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(jsonIn, &out, Use(MessagePack), FallbackWith(JSON))
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("use msgpack, fallback msgpack, error", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(jsonIn, &out, Use(MessagePack), FallbackWith(MessagePack))
		assert.Error(t, err)
	})

	t.Run("use json, fallback json, error", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(msgpackIn, &out, Use(JSON), FallbackWith(JSON))
		assert.Error(t, err)
	})
}
