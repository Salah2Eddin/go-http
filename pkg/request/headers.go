package request

type Headers struct {
	headers map[string][]string
}

func NewRequestHeaders() Headers {
	return Headers{headers: make(map[string][]string)}
}

func (req *Headers) Set(name string, value []string) {
	req.headers[name] = value
}

func (req *Headers) Get(name string) ([]string, bool) {
	val, exists := req.headers[name]
	return val, exists
}
