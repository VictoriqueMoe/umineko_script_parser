package transformer

import sharedtransformer "github.com/VictoriqueMoe/umineko_script_parser/transformer"

type (
	Transformer = sharedtransformer.Transformer
	Format      = sharedtransformer.Format
	Factory     = sharedtransformer.Factory
)

var (
	FormatPlainText = sharedtransformer.FormatPlainText
	FormatHTML      = sharedtransformer.FormatHTML
)
