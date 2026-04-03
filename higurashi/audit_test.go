package higurashi

import (
	"os"
	"testing"
)

func TestAudit_FieldPopulation(t *testing.T) {
	data, err := os.ReadFile("../test/higurashi/en.txt")
	if err != nil {
		t.Skip("test data not found")
	}

	quotes, _, _ := ParseScriptText(string(data))

	total := len(quotes)
	hasText := 0
	hasTextHtml := 0
	hasTextJP := 0
	hasCharacterID := 0
	hasCharacter := 0
	hasAudioID := 0
	hasAudioCharMap := 0
	hasAudioTextMap := 0
	hasEpisode := 0
	hasContentType := 0
	hasSoundEffects := 0
	hasArc := 0

	for i := 0; i < len(quotes); i++ {
		q := &quotes[i]
		if q.Text != "" {
			hasText++
		}
		if q.TextHtml != "" {
			hasTextHtml++
		}
		if q.TextJP != "" {
			hasTextJP++
		}
		if q.CharacterID != "" {
			hasCharacterID++
		}
		if q.Character != "" {
			hasCharacter++
		}
		if q.AudioID != "" {
			hasAudioID++
		}
		if len(q.AudioCharMap) > 0 {
			hasAudioCharMap++
		}
		if len(q.AudioTextMap) > 0 {
			hasAudioTextMap++
		}
		if q.Episode > 0 {
			hasEpisode++
		}
		if q.ContentType != "" {
			hasContentType++
		}
		if len(q.SoundEffects) > 0 {
			hasSoundEffects++
		}
		if q.Arc != "" {
			hasArc++
		}
	}

	t.Logf("Total quotes:    %d", total)
	t.Logf("---")
	t.Logf("Text:            %d (%.0f%%)", hasText, pct(hasText, total))
	t.Logf("TextHtml:        %d (%.0f%%)", hasTextHtml, pct(hasTextHtml, total))
	t.Logf("TextJP:          %d (%.0f%%)", hasTextJP, pct(hasTextJP, total))
	t.Logf("CharacterID:     %d (%.0f%%)", hasCharacterID, pct(hasCharacterID, total))
	t.Logf("Character:       %d (%.0f%%)", hasCharacter, pct(hasCharacter, total))
	t.Logf("AudioID:         %d (%.0f%%)", hasAudioID, pct(hasAudioID, total))
	t.Logf("AudioCharMap:    %d (%.0f%%)", hasAudioCharMap, pct(hasAudioCharMap, total))
	t.Logf("AudioTextMap:    %d (%.0f%%)", hasAudioTextMap, pct(hasAudioTextMap, total))
	t.Logf("Episode:         %d (%.0f%%)", hasEpisode, pct(hasEpisode, total))
	t.Logf("ContentType:     %d (%.0f%%)", hasContentType, pct(hasContentType, total))
	t.Logf("SoundEffects:    %d (%.0f%%)", hasSoundEffects, pct(hasSoundEffects, total))
	t.Logf("Arc:             %d (%.0f%%)", hasArc, pct(hasArc, total))
}

func pct(n, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(n) / float64(total) * 100
}
