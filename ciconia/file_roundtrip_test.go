package ciconia

import (
	"os"
	"testing"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
)

func TestParseFile_RealData(t *testing.T) {
	f, err := os.Open("../test/Ciconia/Ciconia Script Full.file")
	if err != nil {
		t.Skip("encoded test data not found, skipping")
	}
	defer f.Close()

	fileQuotes, _, err := scriptparser.ParseReader(f, NewParser())
	if err != nil {
		t.Fatalf("ParseReader failed: %v", err)
	}

	textData, err := os.ReadFile("../test/Ciconia/Ciconia Script Full.txt")
	if err != nil {
		t.Skip("plain test data not found, skipping")
	}

	textQuotes, _, err := scriptparser.ParseText(string(textData), NewParser())
	if err != nil {
		t.Fatalf("ParseText failed: %v", err)
	}

	if len(fileQuotes) != len(textQuotes) {
		t.Fatalf("quote count mismatch: .file=%d, .txt=%d", len(fileQuotes), len(textQuotes))
	}

	for i := 0; i < len(fileQuotes); i++ {
		if fileQuotes[i].AudioID != textQuotes[i].AudioID {
			t.Errorf("quote %d AudioID mismatch: .file=%q, .txt=%q",
				i, fileQuotes[i].AudioID, textQuotes[i].AudioID)
			break
		}
		if fileQuotes[i].Text != textQuotes[i].Text {
			t.Errorf("quote %d Text mismatch", i)
			break
		}
	}

	t.Logf("round-trip verified: %d quotes parse identically from .txt and .file", len(fileQuotes))
}
