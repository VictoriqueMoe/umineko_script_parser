package transformer

import (
	"testing"

	"github.com/VictoriqueMoe/umineko_script_parser/dialogue"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi/ast"
)

func elements(texts ...string) []dialogue.DialogueElement {
	out := make([]dialogue.DialogueElement, len(texts))
	for i := 0; i < len(texts); i++ {
		out[i] = &ast.PlainText{Text: texts[i]}
	}
	return out
}

func TestPlainText_StripsItalics(t *testing.T) {
	pt := NewPlainTextTransformer()
	got := pt.Transform(elements(`<i>"...Never before has there been a corpse treated with such cruelty"</i>`))
	want := `"...Never before has there been a corpse treated with such cruelty"`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPlainText_StripsBold(t *testing.T) {
	pt := NewPlainTextTransformer()
	got := pt.Transform(elements(`This is <b>bold</b> text`))
	want := `This is bold text`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPlainText_StripsUnderline(t *testing.T) {
	pt := NewPlainTextTransformer()
	got := pt.Transform(elements(`Some <u>underlined</u> words`))
	want := `Some underlined words`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPlainText_PreservesAngleBrackets(t *testing.T) {
	pt := NewPlainTextTransformer()
	got := pt.Transform(elements(`<The Hidden Demon>`))
	want := `<The Hidden Demon>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestPlainText_MixedTagsAndBrackets(t *testing.T) {
	pt := NewPlainTextTransformer()
	got := pt.Transform(elements(`The <i>Kaitai Shinsho</i>, or the <i>New Text on Anatomy.</i>`))
	want := `The Kaitai Shinsho, or the New Text on Anatomy.`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTML_ConvertsItalics(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform(elements(`<i>He was killed because I told him.</i>`))
	want := `<em>He was killed because I told him.</em>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTML_ConvertsBold(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform(elements(`This is <b>important</b>`))
	want := `This is <strong>important</strong>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTML_EscapesAngleBrackets(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform(elements(`<The Hidden Demon>`))
	want := `&lt;The Hidden Demon&gt;`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTML_MixedTagsAndBrackets(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform(elements(`the <i>Kaitai Shinsho</i> and <The Hidden Demon>`))
	want := `the <em>Kaitai Shinsho</em> and &lt;The Hidden Demon&gt;`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTML_MultipleItalicSpans(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform(elements(`the <i>Kaitai Shinsho</i>, or the <i>New Text on Anatomy.</i>`))
	want := `the <em>Kaitai Shinsho</em>, or the <em>New Text on Anatomy.</em>`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTML_EscapesAmpersands(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform(elements(`Tom & Jerry`))
	want := `Tom &amp; Jerry`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestHTML_EscapesQuotes(t *testing.T) {
	ht := NewHtmlTransformer()
	got := ht.Transform(elements(`She said "hello"`))
	want := `She said &#34;hello&#34;`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
