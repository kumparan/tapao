package tapao

import (
	"encoding/json"
	"github.com/vmihailenco/msgpack/v5"
	"testing"

	"google.golang.org/protobuf/proto"

	"github.com/kumparan/tapao/pb"
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

	pbIn := &pb.Greeting{}
	pbOut, err := proto.Marshal(pbIn)
	assert.NoError(t, err)

	t.Run("default", func(t *testing.T) {
		out, err := Marshal(in)
		assert.NoError(t, err)
		assert.Equal(t, msgpackOut, out)
	})

	t.Run("with json", func(t *testing.T) {
		out, err := Marshal(in, With(JSON))
		assert.NoError(t, err)
		assert.Equal(t, jsonOut, out)
	})

	t.Run("with msgpack", func(t *testing.T) {
		out, err := Marshal(in, With(MessagePack))
		assert.NoError(t, err)
		assert.Equal(t, msgpackOut, out)
	})

	t.Run("with protobuf", func(t *testing.T) {
		out, err := Marshal(pbIn, With(Protobuf))
		assert.NoError(t, err)
		assert.Equal(t, out, pbOut)
	})

	// TODO: test fallback, skip for now
}

func TestUnmarshal(t *testing.T) {
	source := testStruct
	msgpackIn, err := msgpack.Marshal(source)
	assert.NoError(t, err)
	jsonIn, err := json.Marshal(source)
	assert.NoError(t, err)

	pbTest := &pb.Greeting{}
	pbIn, err := proto.Marshal(pbTest)
	assert.NoError(t, err)

	t.Run("default", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(msgpackIn, &out)
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("with json", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(jsonIn, &out, With(JSON))
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("with msgpack", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(msgpackIn, &out, With(MessagePack))
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("with protobuf", func(t *testing.T) {
		var out pb.Greeting
		err := Unmarshal(pbIn, &out, With(Protobuf))
		assert.NoError(t, err)
		assert.Equal(t, pbTest, &out)
	})

	t.Run("with json, fallback msgpack", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(msgpackIn, &out, With(JSON), FallbackWith(MessagePack))
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("with msgpack, fallback json", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(jsonIn, &out, With(MessagePack), FallbackWith(JSON))
		assert.NoError(t, err)
		assert.Equal(t, source, out)
	})

	t.Run("with msgpack, fallback msgpack, error", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(jsonIn, &out, With(MessagePack), FallbackWith(MessagePack))
		assert.Error(t, err)
	})

	t.Run("with json, fallback json, error", func(t *testing.T) {
		var out TestStruct
		err := Unmarshal(msgpackIn, &out, With(JSON), FallbackWith(JSON))
		assert.Error(t, err)
	})
}
