package higurashi

import (
	"strconv"
	"strings"

	"github.com/VictoriqueMoe/umineko_script_parser/dialogue"
	"github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi/ast"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"
)

type (
	voiceSegment struct {
		path      string
		charID    string
		contentEN []dialogue.DialogueElement
		contentJP []dialogue.DialogueElement
	}

	rawQuote struct {
		characterName string
		character     character.Character
		segments      []voiceSegment
		unvoicedEN    []dialogue.DialogueElement
		unvoicedJP    []dialogue.DialogueElement
		soundEffects  []dto.SoundEffect
		arc           string
		episode       int
	}

	parseState struct {
		pendingCharName  string
		pendingChar      character.Character
		currentVoicePath string
		currentVoiceID   string
		segments         []voiceSegment
		unvoicedEN       []dialogue.DialogueElement
		unvoicedJP       []dialogue.DialogueElement
		soundEffects     []dto.SoundEffect
		currentArc       string
		currentEpisode   int
		pendingJPPart    string
		awaitingEN       bool
		hasVoice         bool
	}
)

var arcPrefixMap = map[string]struct {
	arc     string
	episode int
}{
	"onik": {"onikakushi", 1},
	"wata": {"watanagashi", 2},
	"tata": {"tatarigoroshi", 3},
	"hima": {"himatsubushi", 4},
	"meak": {"meakashi", 5},
	"tsum": {"tsumihoroboshi", 6},
	"mina": {"minagoroshi", 7},
	"mats": {"matsuribayashi", 8},
	"some": {"someutsushi", 9},
	"kage": {"kageboshi", 10},
	"tsuk": {"tsukiotoshi", 11},
	"tara": {"taraimawashi", 12},
	"yoig": {"yoigoshi", 13},
	"toki": {"tokihogushi", 14},
	"omot": {"miotsukushi_omote", 15},
	"kake": {"kakera", 16},
	"ura":  {"miotsukushi_ura", 17},
	"koto": {"kotohogushi", 18},
	"haji": {"hajisarashi", 19},
	"prol": {"prologue", 0},
}

func parse(lines []string) []rawQuote {
	var (
		state  parseState
		quotes []rawQuote
	)

	for i := 0; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])

		if len(line) == 0 {
			continue
		}

		if strings.HasPrefix(line, "//!file:") {
			state.currentArc, state.currentEpisode = arcFromFilename(line[8:])
			continue
		}

		if strings.HasPrefix(line, "if (GetGlobalFlag(GADVMode)) { OutputLine(\"<color=") {
			name := extractADVCharacterName(line)
			if name != "" {
				state.pendingCharName = name
				state.pendingChar = character.CharacterFromName(name)
			}
			continue
		}

		if strings.HasPrefix(line, "if (GetGlobalFlag(GADVMode)) { OutputLineAll(\"\", NULL") {
			state.pendingCharName = ""
			state.pendingChar = character.Narrator
			continue
		}

		if strings.HasPrefix(line, "ModPlayVoiceLS(") {
			voiceID, voicePath := parseVoiceLine(line)
			if voicePath != "" {
				state.currentVoicePath = voicePath
				state.currentVoiceID = strconv.Itoa(voiceID)
				state.hasVoice = true
			}
			if state.pendingChar == "" && voiceID >= 0 {
				state.pendingChar = character.CharacterFromID(strconv.Itoa(voiceID))
				if state.pendingCharName == "" {
					state.pendingCharName = character.CharacterNames.GetCharacterName(state.pendingChar)
				}
			}
			continue
		}

		if strings.HasPrefix(line, "PlaySE(") || strings.HasPrefix(line, "PlaySE( ") {
			se := parseSELine(line)
			if se.Filename != "" {
				state.soundEffects = append(state.soundEffects, se)
			}
			continue
		}

		if strings.HasPrefix(line, "OutputLine(NULL,") {
			state.pendingJPPart = extractStringArg(line, 16)
			state.awaitingEN = true
			continue
		}

		if state.awaitingEN && strings.Contains(line, "NULL,") {
			en := extractSecondStringArg(line)
			jp := state.pendingJPPart
			state.pendingJPPart = ""
			state.awaitingEN = false

			enElement := &ast.PlainText{Text: en}
			jpElement := &ast.PlainText{Text: jp}

			if state.hasVoice {
				state.segments = append(state.segments, voiceSegment{
					path:      state.currentVoicePath,
					charID:    state.currentVoiceID,
					contentEN: []dialogue.DialogueElement{enElement},
					contentJP: []dialogue.DialogueElement{jpElement},
				})
				state.currentVoicePath = ""
				state.currentVoiceID = ""
				state.hasVoice = false
			} else {
				state.unvoicedEN = append(state.unvoicedEN, enElement)
				state.unvoicedJP = append(state.unvoicedJP, jpElement)
			}
			continue
		}

		if line == "ClearMessage();" || strings.HasPrefix(line, "if (GetGlobalFlag(GADVMode)) { ClearMessage();") {
			if len(state.segments) > 0 || len(state.unvoicedEN) > 0 || len(state.unvoicedJP) > 0 {
				quotes = append(quotes, rawQuote{
					characterName: state.pendingCharName,
					character:     state.pendingChar,
					segments:      state.segments,
					unvoicedEN:    state.unvoicedEN,
					unvoicedJP:    state.unvoicedJP,
					soundEffects:  state.soundEffects,
					arc:           state.currentArc,
					episode:       state.currentEpisode,
				})
			}
			state.pendingCharName = ""
			state.pendingChar = ""
			state.segments = nil
			state.unvoicedEN = nil
			state.unvoicedJP = nil
			state.soundEffects = nil
			state.currentVoicePath = ""
			state.currentVoiceID = ""
			state.pendingJPPart = ""
			state.awaitingEN = false
			state.hasVoice = false
			continue
		}
	}

	return quotes
}

func arcFromFilename(filename string) (string, int) {
	name := strings.TrimSuffix(filename, ".txt")

	for len(name) > 0 && (name[0] == 'z' || name[0] == '_') {
		name = name[1:]
	}

	prefix := ""
	for _, ch := range name {
		if ch >= 'a' && ch <= 'z' {
			prefix += string(ch)
		} else {
			break
		}
	}

	if info, ok := arcPrefixMap[prefix]; ok {
		return info.arc, info.episode
	}

	if strings.HasPrefix(name, "kakera") {
		return "kakera", 16
	}
	if strings.HasPrefix(name, "omake") || strings.HasPrefix(name, "staffroom") {
		return "omake", 0
	}
	if strings.HasPrefix(name, "retrospective") {
		return "retrospective", 0
	}
	if strings.HasPrefix(name, "mio") {
		return "miotsukushi_omote", 15
	}
	if strings.HasPrefix(name, "dummy") {
		return "", 0
	}

	return "unknown", 0
}

func extractADVCharacterName(line string) string {
	marker := "NULL, \"<color=#"
	idx := strings.LastIndex(line, marker)
	if idx == -1 {
		return ""
	}

	rest := line[idx+len(marker):]
	gtIdx := strings.IndexByte(rest, '>')
	if gtIdx == -1 {
		return ""
	}
	rest = rest[gtIdx+1:]

	ltIdx := strings.Index(rest, "</color>")
	if ltIdx == -1 {
		return ""
	}

	return rest[:ltIdx]
}

func parseVoiceLine(line string) (int, string) {
	inner := line[len("ModPlayVoiceLS("):]
	closeIdx := strings.IndexByte(inner, ')')
	if closeIdx == -1 {
		return -1, ""
	}
	inner = inner[:closeIdx]

	parts := strings.Split(inner, ",")
	if len(parts) < 3 {
		return -1, ""
	}

	voiceID, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		return -1, ""
	}

	pathStr := strings.TrimSpace(parts[2])
	pathStr = strings.Trim(pathStr, "\"")

	return voiceID, pathStr
}

func parseSELine(line string) dto.SoundEffect {
	start := strings.IndexByte(line, '(')
	end := strings.IndexByte(line, ')')
	if start == -1 || end == -1 {
		return dto.SoundEffect{}
	}

	parts := strings.Split(line[start+1:end], ",")
	if len(parts) < 2 {
		return dto.SoundEffect{}
	}

	filename := strings.TrimSpace(parts[1])
	filename = strings.Trim(filename, "\"")

	return dto.SoundEffect{
		Filename:  filename,
		AfterClip: -1,
	}
}

func extractStringArg(line string, startIdx int) string {
	idx := strings.IndexByte(line[startIdx:], '"')
	if idx == -1 {
		return ""
	}
	start := startIdx + idx + 1

	var buf strings.Builder
	for i := start; i < len(line); i++ {
		if line[i] == '\\' && i+1 < len(line) {
			next := line[i+1]
			if next == '"' {
				buf.WriteByte('"')
				i++
				continue
			}
			if next == 'n' {
				buf.WriteByte('\n')
				i++
				continue
			}
			buf.WriteByte(line[i])
			continue
		}
		if line[i] == '"' {
			return buf.String()
		}
		buf.WriteByte(line[i])
	}

	return buf.String()
}

func extractSecondStringArg(line string) string {
	idx := strings.Index(line, "NULL,")
	if idx == -1 {
		return ""
	}
	return extractStringArg(line, idx+5)
}
