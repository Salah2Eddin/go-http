package httpheader

type Header struct {
	name   string
	values *[]Value
}

func NewHeader(name string, values *[]Value) *Header {
	return &Header{
		name:   name,
		values: values,
	}
}

func NewHeaderFromString(name string, value string) (*Header, error) {
	values, err := processHeaderValues([]byte(value))
	if err != nil {
		return nil, err
	}
	return NewHeader(name, values), err
}

func (h *Header) Name() string {
	return h.name
}

func (h *Header) Values() *[]Value {
	return h.values
}

func (h *Header) AddValue(value Value) {
	*h.values = append(*h.values, value)
}

func (h *Header) AddValues(values *[]Value) {
	*h.values = append(*h.values, *values...)
}
