package transformer

import "fmt"

type Format int

const (
	FormatPlainText Format = iota
	FormatHTML
)

type Factory struct {
	presets      *PresetContext
	transformers map[Format]Transformer
}

func NewFactory(presets *PresetContext) *Factory {
	f := &Factory{
		presets:      presets,
		transformers: make(map[Format]Transformer),
	}

	f.transformers[FormatPlainText] = NewPlainTextTransformer()
	f.transformers[FormatHTML] = NewHtmlTransformer(presets)

	return f
}

func (f *Factory) Get(format Format) (Transformer, error) {
	if t, ok := f.transformers[format]; ok {
		return t, nil
	}
	return nil, fmt.Errorf("unknown transformer format: %d", format)
}

func (f *Factory) MustGet(format Format) Transformer {
	t, err := f.Get(format)
	if err != nil {
		panic(err)
	}
	return t
}

func (f *Factory) Register(format Format, t Transformer) {
	f.transformers[format] = t
}

func (f *Factory) Presets() *PresetContext {
	return f.presets
}
