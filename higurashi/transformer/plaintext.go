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

	richTextStrip = strings.NewReplacer(
		"<i>", "",
		"</i>", "",
		"<b>", "",
		"</b>", "",
		"<u>", "",
		"</u>", "",
	)

	richTextToHTML = strings.NewReplacer(
		"<i>", "\x00em_open\x00",
		"</i>", "\x00em_close\x00",
		"<b>", "\x00strong_open\x00",
		"</b>", "\x00strong_close\x00",
		"<u>", "\x00u_open\x00",
		"</u>", "\x00u_close\x00",
	)

	richTextPlaceholderRestore = strings.NewReplacer(
		"\x00em_open\x00", "<em>",
		"\x00em_close\x00", "</em>",
		"\x00strong_open\x00", "<strong>",
		"\x00strong_close\x00", "</strong>",
		"\x00u_open\x00", "<u>",
		"\x00u_close\x00", "</u>",
	)
)

func trimText(text string) string {
	text = strings.TrimSpace(text)
	text = strings.TrimLeft(text, "\u3000")
	text = strings.TrimSpace(text)
	return text
}

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
	return trimText(richTextStrip.Replace(buf.String()))
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
	text := richTextToHTML.Replace(buf.String())
	text = trimText(text)
	text = html.EscapeString(text)
	text = richTextPlaceholderRestore.Replace(text)
	return text
}

func NewFactory() *Factory {
	f := sharedtransformer.NewFactory()
	f.Register(FormatPlainText, NewPlainTextTransformer())
	f.Register(FormatHTML, NewHtmlTransformer())
	return f
}
