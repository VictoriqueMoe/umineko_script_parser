package ciconia

import (
	"os"
	"testing"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
)

func TestEpilogueSplit_RealData(t *testing.T) {
	data, err := os.ReadFile("../test/Ciconia/Ciconia Script Full.txt")
	if err != nil {
		t.Skip("test data not found, skipping")
	}

	quotes, _, err := scriptparser.ParseText(string(data), NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var ep25b, ep int
	var epQuote *ParsedQuote
	for i := 0; i < len(quotes); i++ {
		if quotes[i].Chapter == "25b" {
			ep25b++
		}
		if quotes[i].Chapter == "ep" {
			ep++
			if epQuote == nil {
				epQuote = &quotes[i]
			}
		}
	}

	t.Logf("25b quotes: %d | ep quotes: %d", ep25b, ep)

	if ep == 0 {
		t.Error("expected epilogue (ep) chapter to have quotes — p1last movie boundary not triggered")
	}

	if epQuote != nil && epQuote.ContentType != "epilogue" {
		t.Errorf("epilogue ContentType = %q, want 'epilogue'", epQuote.ContentType)
	}

	if epQuote != nil {
		t.Logf("First epilogue quote: [%s] %s", epQuote.Character, epQuote.Text)
	}
}
