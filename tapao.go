package tapao

import (
	"encoding/json"
	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/vmihailenco/msgpack"
)

// SerializerType :nodoc:
type SerializerType string

const (
	// JSON :nodoc:
	JSON = SerializerType(0)
	// MessagePack :nodoc:
	MessagePack = SerializerType(1)
	// Protobuf :nodoc:
	Protobuf = SerializerType(2)
)

type (
	// Options :nodoc:
	Options struct {
		serializer   SerializerType
		fallbackOpts fallbackOptions
	}
	fallbackOptions struct {
		isSet      bool
		serializer SerializerType
	}
)

var defaultOptions = Options{
	serializer:   MessagePack,
	fallbackOpts: fallbackOptions{isSet: false},
}

// With define type of serializer to use
func With(s SerializerType) func(*Options) {
	return func(o *Options) {
		o.serializer = s
	}
}

// FallbackWith define type of serializer to use when primary option failed
func FallbackWith(s SerializerType) func(*Options) {
	return func(o *Options) {
		o.fallbackOpts = fallbackOptions{
			isSet:      true,
			serializer: s,
		}
	}
}

// Marshal with options
func Marshal(in interface{}, opts ...func(*Options)) (out []byte, err error) {
	o := applyOptions(opts)
	out, err = marshal(in, o.serializer)
	if err != nil && o.fallbackOpts.isSet {
		return marshal(in, o.fallbackOpts.serializer)
	}
	return
}

// Unmarshal with options
func Unmarshal(in []byte, out interface{}, opts ...func(*Options)) (err error) {
	o := applyOptions(opts)
	err = unmarshal(in, out, o.serializer)
	if err != nil && o.fallbackOpts.isSet {
		return unmarshal(in, out, o.fallbackOpts.serializer)
	}
	return
}

func applyOptions(opts []func(*Options)) Options {
	to := defaultOptions
	for _, o := range opts {
		o(&to)
	}

	return to
}

func marshal(in interface{}, serializer SerializerType) (out []byte, err error) {
	switch serializer {
	case JSON:
		out, err = json.Marshal(in)
	case MessagePack:
		out, err = msgpack.Marshal(in)
	case Protobuf:
		protoIn, ok := in.(proto.Message)
		if !ok {
			return nil, errors.New("cannot cast input struct to protobuf")
		}
		out, err = proto.Marshal(protoIn)
	default:
		err = errors.New("serializer is not recognized")
	}
	return
}

func unmarshal(in []byte, out interface{}, serializer SerializerType) (err error) {
	switch serializer {
	case JSON:
		err = json.Unmarshal(in, out)
	case MessagePack:
		err = msgpack.Unmarshal(in, out)
	case Protobuf:
		_, ok := out.(proto.Message)
		if !ok {
			return errors.New("cannot cast output struct to protobuf")
		}
		err = proto.Unmarshal(in, out.(proto.Message))
	default:
		err = errors.New("serializer is not recognized")
	}
	return
}
