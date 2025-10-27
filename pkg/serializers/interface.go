package serializers

import "bytes"

type ISerializer[T any] interface {
	Serialize(obj T, buf *bytes.Buffer)
}
