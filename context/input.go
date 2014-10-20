package context

// Request represents the request made from a client
type Input struct {
	ContentType string
	Format      string // html, xml, json, plain, etc...
}

func NewInput() *Input {
	return &Input{}
}
