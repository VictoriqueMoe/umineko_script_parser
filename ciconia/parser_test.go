package ciconia

import (
	"strings"
	"testing"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
)

func TestParseScriptText_SimpleQuote(t *testing.T) {
	input := `*act1
langjp!s0【ナイマ】
langen!s0^Naima:^
langjp!s0「こんにちは」\
langen!s0^"Hello there."^\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.Text != `"Hello there."` {
		t.Errorf("Text = %q", q.Text)
	}
	if q.TextJP != "こんにちは" {
		t.Errorf("TextJP = %q", q.TextJP)
	}
	if q.CharacterID != "naima" {
		t.Errorf("CharacterID = %q", q.CharacterID)
	}
	if q.Character != "Naima" {
		t.Errorf("Character = %q", q.Character)
	}
	if q.Chapter != "01" {
		t.Errorf("Chapter = %q", q.Chapter)
	}
	if q.ContentType != "chapter" {
		t.Errorf("ContentType = %q", q.ContentType)
	}
	if q.Episode != 1 {
		t.Errorf("Episode = %d", q.Episode)
	}
	if !strings.HasPrefix(q.AudioID, "c01:") {
		t.Errorf("AudioID = %q, expected prefix c01:", q.AudioID)
	}
}

func TestParseScriptText_Narrator(t *testing.T) {
	input := `*prologue
langjp!s0　朝の空気は冷たい。\
langen!s0^The morning air is cold.^\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.CharacterID != "narrator" {
		t.Errorf("CharacterID = %q", q.CharacterID)
	}
	if q.Character != "Narrator" {
		t.Errorf("Character = %q", q.Character)
	}
	if q.Chapter != "00" {
		t.Errorf("Chapter = %q", q.Chapter)
	}
	if q.ContentType != "prologue" {
		t.Errorf("ContentType = %q", q.ContentType)
	}
	if q.Text != "The morning air is cold." {
		t.Errorf("Text = %q", q.Text)
	}
	if !strings.HasPrefix(q.AudioID, "pro:") {
		t.Errorf("AudioID = %q, expected prefix pro:", q.AudioID)
	}
}

func TestParseScriptText_KeropoyoDetection(t *testing.T) {
	input := `*act5
langen^#8df270A friend has arrived poyo♪^#FFFFFF\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.CharacterID != "keropoyo" {
		t.Errorf("CharacterID = %q, want 'keropoyo'", q.CharacterID)
	}
	if q.Character != "Keropoyo" {
		t.Errorf("Character = %q, want 'Keropoyo'", q.Character)
	}
	if q.Text != "A friend has arrived poyo♪" {
		t.Errorf("Text = %q", q.Text)
	}
}

func TestParseScriptText_KeropoyoDoesNotOverrideExplicitSpeaker(t *testing.T) {
	input := `*act1
langjp!s0【ナイマ】
langen!s0^#ffd6d6Naima:^
langjp!s0「ケロポヨが鳴いて聞こえない」#FFFFFF\
langen!s0^"My Keropoyo is croaking and I can't hear."^#FFFFFF\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.CharacterID != "naima" {
		t.Errorf("CharacterID = %q, want 'naima' (explicit speaker wins over color)", q.CharacterID)
	}
}

func TestParseScriptText_KeropoyoMultipleConsecutive(t *testing.T) {
	input := `*act5
langen^#8df270First poyo line.^#FFFFFF\
langen^#8df270Second poyo line.^#FFFFFF\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) < 1 {
		t.Fatalf("expected at least 1 quote, got %d", len(quotes))
	}

	for i := 0; i < len(quotes); i++ {
		if quotes[i].CharacterID != "keropoyo" {
			t.Errorf("quote %d CharacterID = %q, want 'keropoyo'", i, quotes[i].CharacterID)
		}
	}
}

func TestParseScriptText_ColorSpan(t *testing.T) {
	input := `*act5
langen^#8df270A friend has arrived poyo♪^#FFFFFF\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.Text != "A friend has arrived poyo♪" {
		t.Errorf("Text = %q", q.Text)
	}
	if !strings.Contains(q.TextHtml, `<span style="color:#8df270">A friend has arrived poyo♪</span>`) {
		t.Errorf("TextHtml missing color span: %q", q.TextHtml)
	}
	if q.Chapter != "05" {
		t.Errorf("Chapter = %q", q.Chapter)
	}
}

func TestParseScriptText_PauseDirective(t *testing.T) {
	input := `*act1
langen^"Wow?^@^ What?^@^ Me?"^\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	want := `"Wow? What? Me?"`
	if q.Text != want {
		t.Errorf("Text = %q, want %q", q.Text, want)
	}
}

func TestParseScriptText_WaitDirective(t *testing.T) {
	input := `*act1
langen^"Fast and ^!w600^slow speech"^\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	if !strings.Contains(quotes[0].Text, "Fast and") && !strings.Contains(quotes[0].Text, "slow speech") {
		t.Errorf("Text = %q", quotes[0].Text)
	}
	if strings.Contains(quotes[0].Text, "!w600") {
		t.Errorf("wait directive not stripped: %q", quotes[0].Text)
	}
}

func TestParseScriptText_MultiLineJP(t *testing.T) {
	input := `*act1
langjp!s0【ナレーター】
langen!s0^Narrator:^
langjp!s0「それが空挺機兵。
langen!s0^"That's it."^\
langjp　通称、ガントレット
langjp　ナイトだ！！」\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.CharacterID != "narrator" {
		t.Errorf("CharacterID = %q", q.CharacterID)
	}
	if !strings.Contains(q.TextJP, "それが空挺機兵") {
		t.Errorf("TextJP missing first part: %q", q.TextJP)
	}
	if !strings.Contains(q.TextJP, "ナイトだ") {
		t.Errorf("TextJP missing last part: %q", q.TextJP)
	}
}

func TestParseScriptText_DataFragment(t *testing.T) {
	input := `*tips003
langjp「データ断片」\
langen^"Data fragment content."^\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.Chapter != "df03" {
		t.Errorf("Chapter = %q", q.Chapter)
	}
	if q.ContentType != "data_fragment" {
		t.Errorf("ContentType = %q", q.ContentType)
	}
	if !strings.HasPrefix(q.AudioID, "df03:") {
		t.Errorf("AudioID = %q, expected prefix df03:", q.AudioID)
	}
}

func TestParseScriptText_BinaryRejection(t *testing.T) {
	_, _, err := scriptparser.ParseText("some \x00 binary", NewParser())
	if err == nil {
		t.Fatal("expected error for binary input")
	}
}

func TestParseScriptText_IDStable(t *testing.T) {
	input := `*act1
langjp!s0【みやお】
langen!s0^Miyao:^
langjp!s0「こんにちは」\
langen!s0^"Hi!"^\
langjp!s0「さようなら」\
langen!s0^"Bye!"^\
`

	a, _, _ := scriptparser.ParseText(input, NewParser())
	b, _, _ := scriptparser.ParseText(input, NewParser())

	if len(a) != len(b) {
		t.Fatalf("length mismatch: %d vs %d", len(a), len(b))
	}
	for i := 0; i < len(a); i++ {
		if a[i].AudioID != b[i].AudioID {
			t.Errorf("AudioID[%d] mismatch: %q vs %q", i, a[i].AudioID, b[i].AudioID)
		}
	}
}

func TestParseScriptText_AllBaseQuoteFieldsValid(t *testing.T) {
	input := `*act1
langjp!s0【みやお】
langen!s0^Miyao:^
langjp!s0「こんにちは」\
langen!s0^"Hi!"^\
`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.AudioID == "" {
		t.Error("AudioID empty — violates BaseQuote contract")
	}
	if q.Text == "" {
		t.Error("Text empty")
	}
	if q.TextHtml == "" {
		t.Error("TextHtml empty")
	}
	if q.CharacterID == "" {
		t.Error("CharacterID empty")
	}
	if q.Character == "" {
		t.Error("Character empty")
	}
	if q.Chapter == "" {
		t.Error("Chapter empty")
	}
	if q.ContentType == "" {
		t.Error("ContentType empty")
	}
	if q.Episode == 0 {
		t.Error("Episode = 0")
	}
}
