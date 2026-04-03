package transformer

import (
	"html"
	"strings"

	"github.com/VictoriqueMoe/umineko_script_parser/dialogue"
	sharedtransformer "github.com/VictoriqueMoe/umineko_script_parser/transformer"
)

type (
	Transformer = sharedtransformer.Transformer
	Format      = sharedtransformer.Format
	Factory     = sharedtransformer.Factory
)

var (
	FormatPlainText = sharedtransformer.FormatPlainText
	FormatHTML      = sharedtransformer.FormatHTML
)

type PlainTextTransformer struct{}

func NewPlainTextTransformer() *PlainTextTransformer {
	return &PlainTextTransformer{}
}

func (t *PlainTextTransformer) Transform(elements []dialogue.DialogueElement) string {
	var buf strings.Builder
	for i := 0; i < len(elements); i++ {
		if te, ok := elements[i].(dialogue.TextElement); ok {
			buf.WriteString(te.GetText())
		}
	}

	text := buf.String()
	text = strings.TrimSpace(text)
	text = strings.TrimLeft(text, "\u3000")
	text = strings.TrimSpace(text)
	return text
}

type HtmlTransformer struct{}

func NewHtmlTransformer() *HtmlTransformer {
	return &HtmlTransformer{}
}

func (t *HtmlTransformer) Transform(elements []dialogue.DialogueElement) string {
	var buf strings.Builder
	for i := 0; i < len(elements); i++ {
		if te, ok := elements[i].(dialogue.TextElement); ok {
			buf.WriteString(te.GetText())
		}
	}

	text := buf.String()
	text = strings.TrimSpace(text)
	text = strings.TrimLeft(text, "\u3000")
	text = strings.TrimSpace(text)
	return html.EscapeString(text)
}

func NewFactory() *Factory {
	f := sharedtransformer.NewFactory()
	f.Register(FormatPlainText, NewPlainTextTransformer())
	f.Register(FormatHTML, NewHtmlTransformer())
	return f
}
