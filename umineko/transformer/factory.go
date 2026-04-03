package transformer

import sharedtransformer "github.com/VictoriqueMoe/umineko_script_parser/transformer"

func NewFactory(presets *PresetContext) *Factory {
	f := sharedtransformer.NewFactory()
	f.Register(FormatPlainText, NewPlainTextTransformer())
	f.Register(FormatHTML, NewHtmlTransformer(presets))
	return f
}
