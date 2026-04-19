package ciconia

import (
	"os"
	"strings"
	"testing"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
)

func TestSpotCheck_KeropoyoGreenLine(t *testing.T) {
	data, err := os.ReadFile("../test/Ciconia/Ciconia Script Full.txt")
	if err != nil {
		t.Skip("test data not found, skipping")
	}

	quotes, _, err := scriptparser.ParseText(string(data), NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var found *ParsedQuote
	for i := 0; i < len(quotes); i++ {
		if quotes[i].Text == "A friend has arrived poyo♪" {
			found = &quotes[i]
			break
		}
	}

	if found == nil {
		t.Fatal("expected Keropoyo green line not found")
	}

	if !strings.Contains(found.TextHtml, `<span style="color:#8df270">`) {
		t.Errorf("expected green span wrapping, got TextHtml=%q", found.TextHtml)
	}

	if found.Chapter == "" {
		t.Error("Chapter empty")
	}
	if found.AudioID == "" {
		t.Error("AudioID empty")
	}
}
