package serializers

import "bytes"

type BodySerializer struct {
}

func (b BodySerializer) Serialize(body *[]byte, buf *bytes.Buffer) {
	buf.Write(*body)
}
