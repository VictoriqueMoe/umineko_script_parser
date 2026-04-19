package higurashi

import (
	"os"
	"strings"
	"testing"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"
)

func TestParseScriptText_BasicQuote(t *testing.T) {
	input := `//!file:onik_001.txt
	if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#956f6e>圭一</color>", NULL, "<color=#956f6e>Keiichi</color>", NULL, Line_ContinueAfterTyping); }
	ModPlayVoiceLS(4, 1, "ps3/s20/01/440100001", 256, TRUE);
	OutputLine(NULL, "「こんにちは」",
			NULL, "\"Hello\"", Line_Normal);
	ClearMessage();`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.Text != "\"Hello\"" {
		t.Errorf("Text = %q", q.Text)
	}
	if q.TextJP != "「こんにちは」" {
		t.Errorf("TextJP = %q", q.TextJP)
	}
	if q.CharacterID != "1" {
		t.Errorf("CharacterID = %q, want \"1\"", q.CharacterID)
	}
	if q.Character != "Keiichi" {
		t.Errorf("Character = %q, want \"Keiichi\"", q.Character)
	}
	if q.Arc != "onikakushi" {
		t.Errorf("Arc = %q, want \"onikakushi\"", q.Arc)
	}
	if q.Episode != 1 {
		t.Errorf("Episode = %d, want 1", q.Episode)
	}
	if q.AudioID != "ps3/s20/01/440100001" {
		t.Errorf("AudioID = %q", q.AudioID)
	}
}

func TestParseScriptText_BinaryRejection(t *testing.T) {
	_, _, err := scriptparser.ParseText("some \x00 binary", NewParser())
	if err == nil {
		t.Fatal("expected error for binary input")
	}
}

func TestParseScriptText_NarratorQuote(t *testing.T) {
	input := `//!file:onik_001.txt
	if (GetGlobalFlag(GADVMode)) { OutputLineAll("", NULL, Line_ContinueAfterTyping); }
	OutputLine(NULL, "　朝の空気は冷たい。",
		   NULL, "The morning air was cold.", Line_Normal);
	ClearMessage();`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	if quotes[0].CharacterID != "narrator" {
		t.Errorf("CharacterID = %q, want narrator", quotes[0].CharacterID)
	}
	if quotes[0].Text != "The morning air was cold." {
		t.Errorf("Text = %q", quotes[0].Text)
	}
}

func TestParseScriptText_RealData(t *testing.T) {
	data, err := os.ReadFile("../test/higurashi/en.txt")
	if err != nil {
		t.Skip("test data not found, skipping")
	}

	quotes, _, err := scriptparser.ParseText(string(data), NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("Total quotes: %d", len(quotes))

	if len(quotes) < 50000 {
		t.Errorf("expected at least 50000 quotes, got %d", len(quotes))
	}

	arcCounts := make(map[string]int)
	charCounts := make(map[string]int)
	for i := 0; i < len(quotes); i++ {
		arcCounts[quotes[i].Arc]++
		charCounts[quotes[i].Character]++
	}

	t.Logf("Arcs: %d unique", len(arcCounts))
	for arc, count := range arcCounts {
		t.Logf("  %s: %d quotes", arc, count)
	}

	t.Logf("Top characters:")
	for _, name := range []string{"Keiichi", "Maebara Keiichi"} {
		if count, ok := charCounts[name]; ok {
			t.Logf("  %s: %d", name, count)
		}
	}

	found := false
	for i := 0; i < len(quotes); i++ {
		if quotes[i].Arc == "onikakushi" && quotes[i].Character != "" {
			found = true
			t.Logf("First Onikakushi quote: [%s] %s", quotes[i].Character, truncate(quotes[i].Text, 80))
			break
		}
	}
	if !found {
		t.Error("no onikakushi quotes found")
	}

	keiichiVoiced := 0
	for i := 0; i < len(quotes); i++ {
		if quotes[i].CharacterID == character.Keiichi.ID() && quotes[i].AudioID != "" {
			keiichiVoiced++
		}
	}
	t.Logf("Keiichi voiced quotes: %d", keiichiVoiced)
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}

func TestParseFile_RealData(t *testing.T) {
	f, err := os.Open("../test/higurashi/en.file")
	if err != nil {
		t.Skip("test data not found, skipping")
	}
	defer f.Close()

	quotes, _, err := scriptparser.ParseReader(f, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("Total quotes from .file: %d", len(quotes))

	if len(quotes) < 50000 {
		t.Errorf("expected at least 50000 quotes, got %d", len(quotes))
	}

	hasEN := 0
	hasJP := 0
	for i := 0; i < len(quotes); i++ {
		if quotes[i].Text != "" {
			hasEN++
		}
		if quotes[i].TextJP != "" {
			hasJP++
		}
	}
	t.Logf("Quotes with EN text: %d", hasEN)
	t.Logf("Quotes with JP text: %d", hasJP)

	arcSeen := make(map[string]bool)
	for i := 0; i < len(quotes); i++ {
		arcSeen[quotes[i].Arc] = true
	}
	expectedArcs := []string{"onikakushi", "watanagashi", "tatarigoroshi", "himatsubushi", "meakashi", "tsumihoroboshi", "minagoroshi", "matsuribayashi"}
	for _, arc := range expectedArcs {
		if !arcSeen[arc] {
			t.Errorf("missing arc: %s", arc)
		}
	}

	firstArcs := make([]string, 0)
	seen := make(map[string]bool)
	for i := 0; i < len(quotes); i++ {
		if !seen[quotes[i].Arc] && quotes[i].Arc != "" {
			seen[quotes[i].Arc] = true
			firstArcs = append(firstArcs, quotes[i].Arc)
		}
	}
	t.Logf("Arc order: %s", strings.Join(firstArcs, " → "))
}
