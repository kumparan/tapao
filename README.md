# tapao

Tapao is a simple wrapper for many serializers. First version of tapao supports `MessagePack` and `JSON`.

## Usage

When using `Marshal` or `Unmarshal`, tapao set the default serializer to `MessagePack`.
You can override it by providing `With` option.
```go
in := "Hello World"

// default
out, err := tapao.Marshal(in)

// override with JSON serializer
out, err := tapao.Marshal(in, tapao.With(tapao.JSON))
```

Sometimes, you need to be flexible when unmarshalling. With tapao, you could set the fallback when primary serializer is failed.

```go
var out string

// default to messagepack
err := tapao.Marshal(in, &out)

// try using different serializer when failed
out, err := tapao.Marshal(in, tapao.With(tapao.JSON), tapao.FallbackWith(tapao.MessagePack))
```