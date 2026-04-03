package transformer

import "fmt"

type Format int

const (
	FormatPlainText Format = iota
	FormatHTML
)

type Factory struct {
	transformers map[Format]Transformer
}

func NewFactory() *Factory {
	return &Factory{
		transformers: make(map[Format]Transformer),
	}
}

func (f *Factory) Register(format Format, t Transformer) {
	f.transformers[format] = t
}

func (f *Factory) Get(format Format) (Transformer, error) {
	t, ok := f.transformers[format]
	if !ok {
		return nil, fmt.Errorf("no transformer registered for format %d", format)
	}
	return t, nil
}

func (f *Factory) MustGet(format Format) Transformer {
	t, err := f.Get(format)
	if err != nil {
		panic(err)
	}
	return t
}
