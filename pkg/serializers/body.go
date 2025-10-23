package serializers

import (
	"bytes"
	"github.com/Salah2Eddin/go-http/pkg/response"
)

// TODO: add serializers for other content encoding types
var bodySerializers = map[string]ISerializer[[]byte]{
	"default": BodySerializer{},
}

func bodySerializerFactory(resp *response.Response) ISerializer[[]byte] {
	encoding, _ := resp.Headers.Get("transfer-encoding")
	if serializer, ok := bodySerializers[encoding.Name()]; ok {
		return serializer
	}
	return bodySerializers["default"]
}

type BodySerializer struct {
}

func (b BodySerializer) Serialize(body *[]byte, buf *bytes.Buffer) {
	buf.Write(*body)
}
