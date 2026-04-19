package transformer

import (
	"testing"

	"github.com/VictoriqueMoe/umineko_script_parser/ciconia/ast"
	"github.com/VictoriqueMoe/umineko_script_parser/dialogue"
)

func plainElements(texts ...string) []dialogue.DialogueElement {
	out := make([]dialogue.DialogueElement, len(texts))
	for i := 0; i < len(texts); i++ {
		out[i] = &ast.PlainText{Text: texts[i]}
	}
	return out
}

func TestPlainText_PreservesText(t *testing.T) {
	pt := NewPlainTextTransformer()
	got := pt.Transform(plainElements("Hello, world!"))
	if got != "Hello, world!" {
		t.Errorf("got %q", got)
	}
}

func TestPlainText_StripsColoredTag(t *testing.T) {
	pt := NewPlainTextTransformer()
	got := pt.Transform([]dialogue.DialogueElement{
		&ast.PlainText{Text: "Hello "},
		&ast.ColoredText{Text: "world", Hex: "#8df270"},
		&ast.PlainText{Text: "!"},
	})
	want := "Hello world!"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHtml_WrapsColoredSpan(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform([]dialogue.DialogueElement{
		&ast.PlainText{Text: "Hello "},
		&ast.ColoredText{Text: "world", Hex: "#8df270"},
		&ast.PlainText{Text: "!"},
	})
	want := `Hello <span style="color:#8df270">world</span>!`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHtml_EscapesAngleBracketsAroundColor(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform([]dialogue.DialogueElement{
		&ast.PlainText{Text: "<warning> "},
		&ast.ColoredText{Text: "danger", Hex: "#ff0000"},
	})
	want := `&lt;warning&gt; <span style="color:#ff0000">danger</span>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHtml_EscapesQuotesInColoredText(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform([]dialogue.DialogueElement{
		&ast.ColoredText{Text: `She said "hi"`, Hex: "#00ff00"},
	})
	want := `<span style="color:#00ff00">She said &#34;hi&#34;</span>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHtml_StripsTrimAroundColoredText(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform([]dialogue.DialogueElement{
		&ast.PlainText{Text: "  "},
		&ast.ColoredText{Text: "core", Hex: "#123456"},
		&ast.PlainText{Text: "  "},
	})
	want := `<span style="color:#123456">core</span>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHtml_PlainOnlyStillEscapes(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform(plainElements("Tom & Jerry"))
	want := "Tom &amp; Jerry"
	if got != want {
		t.Errorf("got %q", got)
	}
}
