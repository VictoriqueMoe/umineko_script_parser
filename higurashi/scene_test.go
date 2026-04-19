package higurashi

import (
	"strings"
	"testing"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
)

func TestParseScriptText_BreakfastScene(t *testing.T) {
	input := `//!file:onik_001.txt
	if (GetGlobalFlag(GADVMode)) { OutputLineAll("", NULL, Line_ContinueAfterTyping); }
	OutputLine(NULL, "　鼻歌混じりに味噌汁のなべを持ってくるお袋は、今朝も上機嫌な様子だった。",
		   NULL, "She hummed a little tune as she brought over the miso soup. It seemed like she was in a good mood today.", Line_Normal);
	ClearMessage();

	if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#f5e6d3>圭一の母</color>", NULL, "<color=#f5e6d3>Keiichi's mom</color>", NULL, Line_ContinueAfterTyping); }
	ModPlayVoiceLS(3, 0, "ps3/s19/00/992700001", 256, TRUE);
	OutputLine(NULL, "「こっちに引っ越してきてから、圭一が早起きになって嬉しいわね。」",
		   NULL, "\"I'm so happy you've been waking up early since we moved here, Keiichi.\"", GetGlobalFlag(GLinemodeSp));
	if (GetGlobalFlag(GADVMode)) { ClearMessage(); } else { OutputLineAll(NULL, "\n", Line_ContinueAfterTyping); }

	if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#956f6e>圭一</color>", NULL, "<color=#956f6e>Keiichi</color>", NULL, Line_ContinueAfterTyping); }
	ModPlayVoiceLS(3, 1, "ps3/s19/01/hr_kei00000", 256, TRUE);
	OutputLine(NULL, "「早起きしないと朝飯を食いそびれるんだよ。」",
		   NULL, "\"If I don't get up early, I don't make it in time for your early breakfasts!\"", GetGlobalFlag(GLinemodeSp));
	if (GetGlobalFlag(GADVMode)) { ClearMessage(); } else { OutputLineAll(NULL, "\n", Line_ContinueAfterTyping); }

	if (GetGlobalFlag(GADVMode)) { OutputLineAll("", NULL, Line_ContinueAfterTyping); }
	OutputLine(NULL, "　よい子ぶりを褒められ、ちょっと悪い子ぶった言い方をする自分がかわいかった。",
		   NULL, "I thought I was being cute, responding with a wise-crack after being praised for being good.", Line_Normal);
	ClearMessage();

	if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#f5e6d3>圭一の母</color>", NULL, "<color=#f5e6d3>Keiichi's mom</color>", NULL, Line_ContinueAfterTyping); }
	ModPlayVoiceLS(3, 0, "ps3/s19/00/992700002", 256, TRUE);
	OutputLine(NULL, "「ご飯はいっぱい？　それとも半分くらいでいい？」",
		   NULL, "\"Full bowl of rice? Or will half be enough?\"", GetGlobalFlag(GLinemodeSp));
	if (GetGlobalFlag(GADVMode)) { ClearMessage(); } else { OutputLineAll(NULL, "\n", Line_ContinueAfterTyping); }

	if (GetGlobalFlag(GADVMode)) { OutputLine("<color=#956f6e>圭一</color>", NULL, "<color=#956f6e>Keiichi</color>", NULL, Line_ContinueAfterTyping); }
	ModPlayVoiceLS(3, 1, "ps3/s19/01/hr_kei00010", 256, TRUE);
	OutputLine(NULL, "「山盛り。」",
		   NULL, "\"Pile it on.\"", GetGlobalFlag(GLinemodeSp));
	if (GetGlobalFlag(GADVMode)) { ClearMessage(); } else { OutputLineAll(NULL, "\n", Line_ContinueAfterTyping); }`

	quotes, _, err := scriptparser.ParseText(input, NewParser())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(quotes) != 6 {
		t.Fatalf("expected 6 quotes, got %d", len(quotes))
	}

	expected := []struct {
		character string
		textStart string
		hasVoice  bool
		arc       string
	}{
		{"Narrator", "She hummed a little tune", false, "onikakushi"},
		{"Keiichi's mom", "\"I'm so happy you've been waking up early", true, "onikakushi"},
		{"Keiichi", "\"If I don't get up early", true, "onikakushi"},
		{"Narrator", "I thought I was being cute", false, "onikakushi"},
		{"Keiichi's mom", "\"Full bowl of rice?", true, "onikakushi"},
		{"Keiichi", "\"Pile it on.\"", true, "onikakushi"},
	}

	for i, exp := range expected {
		q := quotes[i]
		if q.Character != exp.character {
			t.Errorf("quote[%d] Character = %q, want %q", i, q.Character, exp.character)
		}
		if !strings.HasPrefix(q.Text, exp.textStart) {
			t.Errorf("quote[%d] Text starts with %q, want prefix %q", i, truncate(q.Text, 50), exp.textStart)
		}
		if exp.hasVoice && q.AudioID == "" {
			t.Errorf("quote[%d] expected voice, got empty AudioID", i)
		}
		if !exp.hasVoice && q.AudioID != "" {
			t.Errorf("quote[%d] expected no voice, got %q", i, q.AudioID)
		}
		if q.Arc != exp.arc {
			t.Errorf("quote[%d] Arc = %q, want %q", i, q.Arc, exp.arc)
		}
		if q.TextJP == "" {
			t.Errorf("quote[%d] TextJP is empty", i)
		}
	}

	if quotes[1].AudioID != "ps3/s19/00/992700001" {
		t.Errorf("mom's voice = %q, want ps3/s19/00/992700001", quotes[1].AudioID)
	}
	if quotes[2].AudioID != "ps3/s19/01/hr_kei00000" {
		t.Errorf("keiichi's voice = %q, want ps3/s19/01/hr_kei00000", quotes[2].AudioID)
	}
}
