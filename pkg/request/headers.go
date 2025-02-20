package request

type RequestHeaders struct {
	headers map[string]string
}

func NewRequestHeaders() *RequestHeaders {
	return &RequestHeaders{headers: make(map[string]string)}
}

func (req *RequestHeaders) Set(name string, value string) {
	req.headers[name] = value
}

func (req *RequestHeaders) Get(name string) (string, bool) {
	val, exists := req.headers[name]
	return val, exists
}
