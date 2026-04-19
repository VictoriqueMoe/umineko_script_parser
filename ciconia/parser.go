package ciconia

import (
	"regexp"
	"strings"

	"github.com/VictoriqueMoe/umineko_script_parser/ciconia/ast"
	"github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"
	"github.com/VictoriqueMoe/umineko_script_parser/dialogue"
)

type (
	rawQuote struct {
		characterName string
		character     character.Character
		bodyEN        []dialogue.DialogueElement
		bodyJP        []dialogue.DialogueElement
		chapter       string
		contentType   string
	}

	parseState struct {
		currentChapter     string
		currentContentType string
		pendingCharName    string
		pendingChar        character.Character
		bufferEN           []dialogue.DialogueElement
		bufferJP           []dialogue.DialogueElement
		hasContent         bool
		pendingCharSet     bool
	}
)

var (
	actLabelRe     = regexp.MustCompile(`^\*act(\d+)(b)?$`)
	tipsLabelRe    = regexp.MustCompile(`^\*tips(\d+)$`)
	enCharMarkerRe = regexp.MustCompile(`^lang(?:en)(?:!s\d+)?\^(?:#[0-9A-Fa-f]{6})?([^^]+):\^$`)
	keropoyoOpenRe = regexp.MustCompile(`^langen(?:!s\d+)?\^#8df270`)
)

const prologueLabel = "*prologue"

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

		if strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "*") {
			chapter, contentType, ok := matchChapterLabel(line)
			if ok {
				flush(&state, &quotes)
				state.currentChapter = chapter
				state.currentContentType = contentType
				state.pendingCharName = ""
				state.pendingChar = ""
				state.pendingCharSet = false
			}
			continue
		}

		if state.currentChapter == "25b" && strings.HasPrefix(line, "movie ") && strings.Contains(line, `\p1last\`) {
			flush(&state, &quotes)
			state.currentChapter = "ep"
			state.currentContentType = "epilogue"
			state.pendingCharName = ""
			state.pendingChar = ""
			state.pendingCharSet = false
			continue
		}

		if state.currentChapter == "" {
			continue
		}

		if match := enCharMarkerRe.FindStringSubmatch(line); match != nil {
			flush(&state, &quotes)
			name := strings.TrimSpace(match[1])
			state.pendingCharName = name
			state.pendingChar = character.CharacterFromName(name)
			state.pendingCharSet = true
			continue
		}

		if isJPCharMarker(line) {
			continue
		}

		if strings.HasPrefix(line, "langen") {
			if keropoyoOpenRe.MatchString(line) && !state.pendingCharSet {
				flush(&state, &quotes)
				state.pendingCharName = "Keropoyo"
				state.pendingChar = character.Keropoyo
				state.pendingCharSet = true
			}
			body := stripLangPrefix(line, "langen")
			elements := parseEnBody(body)
			state.bufferEN = append(state.bufferEN, elements...)
			if len(elements) > 0 {
				state.hasContent = true
			}
			continue
		}

		if strings.HasPrefix(line, "langjp") {
			body := stripLangPrefix(line, "langjp")
			elements := parseJpBody(body)
			state.bufferJP = append(state.bufferJP, elements...)
			if len(elements) > 0 {
				state.hasContent = true
			}
			continue
		}

		flush(&state, &quotes)
		state.pendingCharName = ""
		state.pendingChar = ""
		state.pendingCharSet = false
	}

	flush(&state, &quotes)
	return quotes
}

func flush(state *parseState, quotes *[]rawQuote) {
	if !state.hasContent {
		state.bufferEN = nil
		state.bufferJP = nil
		return
	}

	char := state.pendingChar
	charName := state.pendingCharName
	if !state.pendingCharSet {
		char = character.Narrator
		charName = ""
	}

	*quotes = append(*quotes, rawQuote{
		characterName: charName,
		character:     char,
		bodyEN:        state.bufferEN,
		bodyJP:        state.bufferJP,
		chapter:       state.currentChapter,
		contentType:   state.currentContentType,
	})

	state.bufferEN = nil
	state.bufferJP = nil
	state.hasContent = false
}

func matchChapterLabel(line string) (string, string, bool) {
	if line == prologueLabel {
		return "00", "prologue", true
	}

	if m := actLabelRe.FindStringSubmatch(line); m != nil {
		num := m[1]
		if len(num) < 2 {
			num = "0" + num
		}
		chapter := num
		if m[2] != "" {
			chapter = num + "b"
		}
		return chapter, "chapter", true
	}

	if m := tipsLabelRe.FindStringSubmatch(line); m != nil {
		num := m[1]
		if len(num) >= 3 {
			num = num[len(num)-2:]
		}
		if len(num) < 2 {
			num = "0" + num
		}
		return "df" + num, "data_fragment", true
	}

	return "", "", false
}

func isJPCharMarker(line string) bool {
	body := stripLangPrefix(line, "langjp")
	if body == "" {
		return false
	}
	body = strings.TrimPrefix(body, "#")
	if len(body) >= 6 && isHex6(body[:6]) {
		body = body[6:]
	}
	if !strings.HasPrefix(body, "【") {
		return false
	}
	if !strings.HasSuffix(body, "】") {
		return false
	}
	return true
}

func stripLangPrefix(line, lang string) string {
	if !strings.HasPrefix(line, lang) {
		return line
	}
	rest := line[len(lang):]
	if strings.HasPrefix(rest, "!s") {
		idx := 2
		for idx < len(rest) && rest[idx] >= '0' && rest[idx] <= '9' {
			idx++
		}
		rest = rest[idx:]
	}
	return rest
}

func parseEnBody(body string) []dialogue.DialogueElement {
	body = strings.TrimSuffix(body, "\\")

	var (
		elements     []dialogue.DialogueElement
		buf          strings.Builder
		currentColor string
	)

	flushBuf := func() {
		if buf.Len() == 0 {
			return
		}
		text := buf.String()
		buf.Reset()
		if currentColor == "" {
			elements = append(elements, &ast.PlainText{Text: text})
		} else {
			elements = append(elements, &ast.ColoredText{Text: text, Hex: currentColor})
		}
	}

	i := 0
	for i < len(body) {
		if body[i] == '^' {
			if i+7 < len(body) && body[i+1] == '#' && isHex6(body[i+2:i+8]) {
				flushBuf()
				hex := body[i+1 : i+8]
				if strings.EqualFold(hex, "#FFFFFF") {
					currentColor = ""
				} else {
					currentColor = hex
				}
				i += 8
				continue
			}
			if i+1 < len(body) && body[i+1] == '@' {
				flushBuf()
				if i+2 < len(body) && body[i+2] == '^' {
					i += 3
				} else {
					i += 2
				}
				continue
			}
			if i+2 < len(body) && body[i+1] == '!' && body[i+2] == 'w' {
				end := strings.IndexByte(body[i+1:], '^')
				if end == -1 {
					i = len(body)
				} else {
					i += 1 + end + 1
				}
				continue
			}
			i++
			continue
		}
		buf.WriteByte(body[i])
		i++
	}
	flushBuf()
	return elements
}

func parseJpBody(body string) []dialogue.DialogueElement {
	body = strings.TrimSuffix(body, "\\")
	body = strings.TrimSuffix(body, "#FFFFFF")
	body = strings.TrimPrefix(body, "#")
	if len(body) >= 6 && isHex6(body[:6]) {
		body = body[6:]
	}
	body = strings.TrimLeft(body, "\u3000")
	body = strings.TrimPrefix(body, "「")
	body = strings.TrimSuffix(body, "」")

	var buf strings.Builder
	runes := []rune(body)
	i := 0
	for i < len(runes) {
		if runes[i] == '@' {
			i++
			continue
		}
		if runes[i] == '!' && i+1 < len(runes) && runes[i+1] == 'w' {
			j := i + 2
			for j < len(runes) && runes[j] >= '0' && runes[j] <= '9' {
				j++
			}
			i = j
			continue
		}
		buf.WriteRune(runes[i])
		i++
	}

	text := strings.TrimSpace(buf.String())
	if text == "" {
		return nil
	}
	return []dialogue.DialogueElement{&ast.PlainText{Text: text}}
}

func isHex6(s string) bool {
	if len(s) != 6 {
		return false
	}
	for i := 0; i < 6; i++ {
		c := s[i]
		isHex := (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')
		if !isHex {
			return false
		}
	}
	return true
}
