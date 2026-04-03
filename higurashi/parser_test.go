package higurashi

import (
	"strings"
	"testing"

	"github.com/VictoriqueMoe/umineko_script_parser/dialogue"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"
)

func elementText(elements []dialogue.DialogueElement) string {
	var buf strings.Builder
	for i := 0; i < len(elements); i++ {
		if te, ok := elements[i].(dialogue.TextElement); ok {
			buf.WriteString(te.GetText())
		}
	}
	return buf.String()
}

func TestScan_BasicVoicedQuote(t *testing.T) {
	input := `//!file:onik_001.txt
	if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#956f6e>圭一</color>", NULL, "<color=#956f6e>Keiichi</color>", NULL, Line_ContinueAfterTyping); }
	ModPlayVoiceLS(4, 1, "ps3/s20/01/440100001", 256, TRUE);
	OutputLine(NULL, "「テスト台詞」",
			NULL, "\"Test dialogue\"", Line_Normal);
	ClearMessage();`

	quotes := parse(strings.Split(input, "\n"))

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.character != character.Keiichi {
		t.Errorf("character = %q, want keiichi", q.character)
	}
	if q.characterName != "Keiichi" {
		t.Errorf("characterName = %q, want Keiichi", q.characterName)
	}
	if len(q.segments) != 1 {
		t.Fatalf("expected 1 segment, got %d", len(q.segments))
	}
	if q.segments[0].path != "ps3/s20/01/440100001" {
		t.Errorf("voice path = %q", q.segments[0].path)
	}
	if q.segments[0].charID != "1" {
		t.Errorf("voice charID = %q, want 1", q.segments[0].charID)
	}
	if elementText(q.segments[0].contentEN) != "\"Test dialogue\"" {
		t.Errorf("textEN = %q", elementText(q.segments[0].contentEN))
	}
	if elementText(q.segments[0].contentJP) != "「テスト台詞」" {
		t.Errorf("textJP = %q", elementText(q.segments[0].contentJP))
	}
	if q.arc != "onikakushi" {
		t.Errorf("arc = %q, want onikakushi", q.arc)
	}
	if q.episode != 1 {
		t.Errorf("episode = %d, want 1", q.episode)
	}
}

func TestScan_NarratorQuote(t *testing.T) {
	input := `//!file:onik_001.txt
	if (GetGlobalFlag(GADVMode)) { OutputLineAll("", NULL, Line_ContinueAfterTyping); }
	OutputLine(NULL, "　朝の空気は冷たい。",
		   NULL, "The morning air was cold.", Line_Normal);
	ClearMessage();`

	quotes := parse(strings.Split(input, "\n"))

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.character != character.Narrator {
		t.Errorf("character = %q, want narrator", q.character)
	}
	if len(q.segments) != 0 {
		t.Error("narrator should have no voice segments")
	}
	if len(q.unvoicedEN) != 1 {
		t.Errorf("expected 1 unvoiced EN element, got %d", len(q.unvoicedEN))
	}
}

func TestScan_MultiVoiceQuote(t *testing.T) {
	input := `//!file:onik_002.txt
	if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#5ec69a>魅音</color>", NULL, "<color=#5ec69a>Mion</color>", NULL, Line_ContinueAfterTyping); }
	ModPlayVoiceLS(4, 3, "ps3/s20/03/440300001", 256, TRUE);
	OutputLine(NULL, "「最初の台詞」",
			NULL, "\"First line\"", Line_WaitForInput);
	ModPlayVoiceLS(4, 3, "ps3/s20/03/440300002", 256, TRUE);
	OutputLine(NULL, "「二番目の台詞」",
			NULL, "\"Second line\"", Line_Normal);
	ClearMessage();`

	quotes := parse(strings.Split(input, "\n"))

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}

	q := quotes[0]
	if q.character != character.Mion {
		t.Errorf("character = %q, want mion", q.character)
	}
	if len(q.segments) != 2 {
		t.Fatalf("expected 2 segments, got %d", len(q.segments))
	}
	if q.segments[0].path != "ps3/s20/03/440300001" {
		t.Errorf("segment[0] path = %q", q.segments[0].path)
	}
	if q.segments[1].path != "ps3/s20/03/440300002" {
		t.Errorf("segment[1] path = %q", q.segments[1].path)
	}
	if elementText(q.segments[0].contentEN) != "\"First line\"" {
		t.Errorf("segment[0] textEN = %q", elementText(q.segments[0].contentEN))
	}
	if elementText(q.segments[1].contentEN) != "\"Second line\"" {
		t.Errorf("segment[1] textEN = %q", elementText(q.segments[1].contentEN))
	}
}

func TestScan_SoundEffect(t *testing.T) {
	input := `//!file:onik_001.txt
	PlaySE( 4, "wa_020", 56, 64 );
	if (GetGlobalFlag(GADVMode)) { OutputLineAll("", NULL, Line_ContinueAfterTyping); }
	OutputLine(NULL, "テスト",
		   NULL, "Test", Line_Normal);
	ClearMessage();`

	quotes := parse(strings.Split(input, "\n"))

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}
	if len(quotes[0].soundEffects) != 1 {
		t.Fatalf("expected 1 sound effect, got %d", len(quotes[0].soundEffects))
	}
	if quotes[0].soundEffects[0].Filename != "wa_020" {
		t.Errorf("SE filename = %q, want wa_020", quotes[0].soundEffects[0].Filename)
	}
}

func TestScan_ConditionalBoundary(t *testing.T) {
	input := `//!file:onik_001.txt
	if (GetGlobalFlag(GADVMode)) { OutputLineAll("", NULL, Line_ContinueAfterTyping); }
	OutputLine(NULL, "テスト",
		   NULL, "Test", GetGlobalFlag(GLinemodeSp));
	if (GetGlobalFlag(GADVMode)) { ClearMessage(); } else { OutputLineAll(NULL, "\n", Line_ContinueAfterTyping); }`

	quotes := parse(strings.Split(input, "\n"))

	if len(quotes) != 1 {
		t.Fatalf("expected 1 quote, got %d", len(quotes))
	}
}

func TestScan_ArcSwitching(t *testing.T) {
	input := `//!file:onik_001.txt
	if (GetGlobalFlag(GADVMode)) { OutputLineAll("", NULL, Line_ContinueAfterTyping); }
	OutputLine(NULL, "第一",
		   NULL, "First arc", Line_Normal);
	ClearMessage();
//!file:wata_001.txt
	if (GetGlobalFlag(GADVMode)) { OutputLineAll("", NULL, Line_ContinueAfterTyping); }
	OutputLine(NULL, "第二",
		   NULL, "Second arc", Line_Normal);
	ClearMessage();`

	quotes := parse(strings.Split(input, "\n"))

	if len(quotes) != 2 {
		t.Fatalf("expected 2 quotes, got %d", len(quotes))
	}
	if quotes[0].arc != "onikakushi" || quotes[0].episode != 1 {
		t.Errorf("first quote: arc=%q episode=%d", quotes[0].arc, quotes[0].episode)
	}
	if quotes[1].arc != "watanagashi" || quotes[1].episode != 2 {
		t.Errorf("second quote: arc=%q episode=%d", quotes[1].arc, quotes[1].episode)
	}
}

func TestScan_EmptyBoundarySkipped(t *testing.T) {
	input := `//!file:onik_001.txt
	ClearMessage();
	ClearMessage();`

	quotes := parse(strings.Split(input, "\n"))

	if len(quotes) != 0 {
		t.Errorf("expected 0 quotes, got %d", len(quotes))
	}
}

func TestArcFromFilename(t *testing.T) {
	tests := []struct {
		filename    string
		wantArc     string
		wantEpisode int
	}{
		{"onik_001.txt", "onikakushi", 1},
		{"wata_001.txt", "watanagashi", 2},
		{"tata_001.txt", "tatarigoroshi", 3},
		{"hima_001.txt", "himatsubushi", 4},
		{"meak_001.txt", "meakashi", 5},
		{"tsum_001.txt", "tsumihoroboshi", 6},
		{"mina_001.txt", "minagoroshi", 7},
		{"mats_001.txt", "matsuribayashi", 8},
		{"some_001.txt", "someutsushi", 9},
		{"kage_001.txt", "kageboshi", 10},
		{"zonik_003_vm00_n01.txt", "onikakushi", 1},
		{"_meak_001.txt", "meakashi", 5},
		{"omake_01.txt", "omake", 0},
		{"prol_001.txt", "prologue", 0},
		{"koto_001.txt", "kotohogushi", 18},
		{"haji_001.txt", "hajisarashi", 19},
		{"onik_tips_01.txt", "onikakushi", 1},
		{"_kakera01.txt", "kakera", 16},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			arc, ep := arcFromFilename(tt.filename)
			if arc != tt.wantArc {
				t.Errorf("arc = %q, want %q", arc, tt.wantArc)
			}
			if ep != tt.wantEpisode {
				t.Errorf("episode = %d, want %d", ep, tt.wantEpisode)
			}
		})
	}
}

func TestExtractADVCharacterName(t *testing.T) {
	tests := []struct {
		line string
		want string
	}{
		{
			`if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#956f6e>圭一</color>", NULL, "<color=#956f6e>Keiichi</color>", NULL, Line_ContinueAfterTyping); }`,
			"Keiichi",
		},
		{
			`if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#5ec69a>魅音</color>", NULL, "<color=#5ec69a>Mion</color>", NULL, Line_ContinueAfterTyping); }`,
			"Mion",
		},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := extractADVCharacterName(tt.line)
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
