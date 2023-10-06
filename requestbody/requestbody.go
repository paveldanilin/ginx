package requestbody

type RequestBody interface {
	RequestBodyFormat() string
}

// JSON

type JSONData map[string]any

func (j JSONData) RequestBodyFormat() string {
	return "json"
}

type JSON struct{}

func (j JSON) RequestBodyFormat() string {
	return "json"
}

// XML

type XMLData map[string]any

func (x XMLData) RequestBodyFormat() string {
	return "xml"
}

type XML struct{}

func (x XML) RequestBodyFormat() string {
	return "xml"
}
