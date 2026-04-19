package ciconia

import (
	"os"
	"sort"
	"strings"
	"testing"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
)

func TestParseScriptText_RealData(t *testing.T) {
	data, err := os.ReadFile("../test/Ciconia/Ciconia Script Full.txt")
	if err != nil {
		t.Skip("test data not found, skipping")
	}

	quotes, _, err := scriptparser.ParseText(string(data), NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	t.Logf("Total quotes: %d", len(quotes))

	if len(quotes) < 5000 {
		t.Errorf("expected at least 5000 quotes, got %d", len(quotes))
	}

	chapterCounts := make(map[string]int)
	charCounts := make(map[string]int)
	contentTypeCounts := make(map[string]int)
	missingAudioID := 0
	missingText := 0
	missingChapter := 0
	missingContentType := 0
	episodeNotOne := 0
	coloredHtml := 0

	for i := 0; i < len(quotes); i++ {
		q := quotes[i]
		chapterCounts[q.Chapter]++
		charCounts[q.Character]++
		contentTypeCounts[q.ContentType]++

		if q.AudioID == "" {
			missingAudioID++
		}
		if q.Text == "" && q.TextJP == "" {
			missingText++
		}
		if q.Chapter == "" {
			missingChapter++
		}
		if q.ContentType == "" {
			missingContentType++
		}
		if q.Episode != 1 {
			episodeNotOne++
		}
		if strings.Contains(q.TextHtml, `<span style="color:#`) {
			coloredHtml++
		}
	}

	if missingAudioID > 0 {
		t.Errorf("%d quotes missing AudioID", missingAudioID)
	}
	if missingChapter > 0 {
		t.Errorf("%d quotes missing Chapter", missingChapter)
	}
	if missingContentType > 0 {
		t.Errorf("%d quotes missing ContentType", missingContentType)
	}
	if episodeNotOne > 0 {
		t.Errorf("%d quotes have Episode != 1", episodeNotOne)
	}

	t.Logf("Colored HTML quotes: %d", coloredHtml)

	t.Logf("Content types: %v", contentTypeCounts)

	t.Logf("Chapters: %d unique", len(chapterCounts))

	expectedChapters := []string{"00"}
	for i := 1; i <= 25; i++ {
		if i < 10 {
			expectedChapters = append(expectedChapters, "0"+itoa(i))
		} else {
			expectedChapters = append(expectedChapters, itoa(i))
		}
	}
	expectedChapters = append(expectedChapters, "25b")
	for i := 1; i <= 16; i++ {
		if i < 10 {
			expectedChapters = append(expectedChapters, "df0"+itoa(i))
		} else {
			expectedChapters = append(expectedChapters, "df"+itoa(i))
		}
	}

	for _, ch := range expectedChapters {
		if chapterCounts[ch] == 0 {
			t.Errorf("chapter %q has zero quotes", ch)
		} else {
			t.Logf("  %s: %d quotes", ch, chapterCounts[ch])
		}
	}

	topChars := topN(charCounts, 20)
	t.Logf("Top characters:")
	for _, p := range topChars {
		t.Logf("  %s: %d", p.name, p.count)
	}

	expectedMainCast := []string{"Miyao", "Jayden", "Narrator", "Koshka"}
	for _, name := range expectedMainCast {
		if charCounts[name] == 0 {
			t.Errorf("main cast character %q has zero quotes", name)
		}
	}

	idSet := make(map[string]bool, len(quotes))
	dupes := 0
	for i := 0; i < len(quotes); i++ {
		if idSet[quotes[i].AudioID] {
			dupes++
		}
		idSet[quotes[i].AudioID] = true
	}
	if dupes > 0 {
		t.Errorf("%d duplicate AudioIDs found", dupes)
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	return string(digits)
}

type charCountPair struct {
	name  string
	count int
}

func topN(m map[string]int, n int) []charCountPair {
	pairs := make([]charCountPair, 0, len(m))
	for k, v := range m {
		pairs = append(pairs, charCountPair{k, v})
	}
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].count > pairs[j].count
	})
	if len(pairs) < n {
		return pairs
	}
	return pairs[:n]
}
