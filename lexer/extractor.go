package lexer

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/lexer/ast"
	"github.com/VictoriqueMoe/umineko_script_parser/lexer/transformer"
)

type (
	QuoteExtractor struct {
		presets          *transformer.PresetContext
		strAliases       map[string]string
		subtitleRefs     []SubtitleRef
		validationErrors []ValidationError
		seFileMap        map[int]string
	}

	ExtractedQuote struct {
		Content      []ast.DialogueElement
		CharacterID  string
		AudioID      string
		AudioCharMap map[string]string
		AudioTextMap map[string][]ast.DialogueElement
		Episode      int
		ContentType  string
		Truth        TruthFlags
		SoundEffects []pendingSoundEffect
	}

	pendingSoundEffect struct {
		seNum     int
		afterClip int
	}

	SubtitleRef struct {
		SubPath     string
		AudioID     string
		CharacterID string
		Episode     int
		ContentType string
	}

	topLevelVoice struct {
		characterID string
		audioID     string
	}
)

func NewQuoteExtractor() *QuoteExtractor {
	return &QuoteExtractor{
		presets: transformer.NewPresetContext(),
	}
}

func (e *QuoteExtractor) SetSeFileMap(m map[int]string) {
	e.seFileMap = m
}

func (e *QuoteExtractor) ExtractQuotes(input string) []ExtractedQuote {
	script := Parse(input)
	e.validationErrors = Validate(script)
	return e.ExtractFromScript(script)
}

func (e *QuoteExtractor) ValidationErrors() []ValidationError {
	return e.validationErrors
}

func (e *QuoteExtractor) ExtractFromScript(script *ast.Script) []ExtractedQuote {
	e.presets.CollectFromScript(script)
	e.strAliases = make(map[string]string)
	e.subtitleRefs = nil

	for _, line := range script.Lines {
		if l, ok := line.(*ast.StraliasLine); ok {
			if l.Name != "" && l.Value != "" {
				e.strAliases[l.Name] = l.Value
			}
		}
	}

	var quotes []ExtractedQuote
	currentEpisode := 0
	currentContentType := ""
	var lastVoice *topLevelVoice

	var pendingSEForNext []pendingSoundEffect
	var postDialogueSEs []pendingSoundEffect
	inWaitOnDBlock := false
	currentWaitSegment := -1
	lastQuoteIdx := -1

	for _, line := range script.Lines {
		switch l := line.(type) {
		case *ast.EpisodeMarkerLine:
			currentEpisode = l.Episode
			if l.Type == "episode" {
				currentContentType = ""
			} else {
				currentContentType = l.Type
			}

		case *ast.LabelLine:
			if ep, ok := e.parseOmakeEpisode(l.Name); ok {
				currentEpisode = ep
				currentContentType = "omake"
			}

		case *ast.CommandLine:
			if l.Command == "lv" && len(l.Args) >= 3 {
				lastVoice = &topLevelVoice{
					characterID: strings.Trim(l.Args[1].Value, `"`),
					audioID:     strings.Trim(l.Args[2].Value, `"`),
				}
			}

		case *ast.SsaLoadLine:
			if l.SubAlias != "" && lastVoice != nil {
				resolved, ok := e.strAliases[l.SubAlias]
				if !ok {
					break
				}
				e.subtitleRefs = append(e.subtitleRefs, SubtitleRef{
					SubPath:     resolved,
					AudioID:     lastVoice.audioID,
					CharacterID: lastVoice.characterID,
					Episode:     currentEpisode,
					ContentType: currentContentType,
				})
			}

		case *ast.WaitOnDLine:
			inWaitOnDBlock = true
			currentWaitSegment = l.Segment

		case *ast.SeplayLine:
			if inWaitOnDBlock {
				postDialogueSEs = append(postDialogueSEs, pendingSoundEffect{
					seNum:     l.SeNum,
					afterClip: currentWaitSegment,
				})
			} else {
				pendingSEForNext = append(pendingSEForNext, pendingSoundEffect{
					seNum:     l.SeNum,
					afterClip: -1,
				})
			}

		case *ast.DialogueLine:
			if len(postDialogueSEs) > 0 && lastQuoteIdx >= 0 {
				quotes[lastQuoteIdx].SoundEffects = append(
					quotes[lastQuoteIdx].SoundEffects,
					postDialogueSEs...,
				)
				postDialogueSEs = nil
			}

			quote := e.extractFromDialogue(l)
			if quote != nil {
				if currentEpisode > 0 {
					quote.Episode = currentEpisode
				}
				quote.ContentType = currentContentType
				if len(pendingSEForNext) > 0 {
					quote.SoundEffects = append(quote.SoundEffects, pendingSEForNext...)
					pendingSEForNext = nil
				}
				quotes = append(quotes, *quote)
				lastQuoteIdx = len(quotes) - 1
			}

			inWaitOnDBlock = false
			currentWaitSegment = -1
		}
	}

	if len(postDialogueSEs) > 0 && lastQuoteIdx >= 0 {
		quotes[lastQuoteIdx].SoundEffects = append(
			quotes[lastQuoteIdx].SoundEffects,
			postDialogueSEs...,
		)
	}

	return quotes
}

func (e *QuoteExtractor) SubtitleRefs() []SubtitleRef {
	return e.subtitleRefs
}

func (e *QuoteExtractor) extractFromDialogue(d *ast.DialogueLine) *ExtractedQuote {
	voices := d.GetVoiceCommands()
	truth := DetectTruth(d.Content, e.presets)

	if len(voices) == 0 || hasWordsBeforeVoice(d.Content) {
		return &ExtractedQuote{
			Content:     d.Content,
			CharacterID: "narrator",
			Truth:       truth,
		}
	}

	characterID := voices[0].CharacterID

	seen := make(map[string]bool)
	var audioIDs []string
	multiChar := false
	for _, v := range voices {
		if !seen[v.AudioID] {
			seen[v.AudioID] = true
			audioIDs = append(audioIDs, v.AudioID)
			if v.CharacterID != characterID {
				multiChar = true
			}
		}
	}

	var audioCharMap map[string]string
	if multiChar {
		audioCharMap = make(map[string]string, len(audioIDs))
		for _, v := range voices {
			if audioCharMap[v.AudioID] == "" {
				audioCharMap[v.AudioID] = v.CharacterID
			}
		}
	}

	episode := 0
	if len(audioIDs) > 0 && len(audioIDs[0]) > 0 {
		ep := int(audioIDs[0][0] - '0')
		if ep >= 1 && ep <= 8 {
			episode = ep
		}
	}

	var audioTextMap map[string][]ast.DialogueElement
	if len(audioIDs) > 1 {
		audioTextMap = buildAudioTextMap(d.Content)
	}

	return &ExtractedQuote{
		Content:      d.Content,
		CharacterID:  characterID,
		AudioID:      strings.Join(audioIDs, ", "),
		AudioCharMap: audioCharMap,
		AudioTextMap: audioTextMap,
		Episode:      episode,
		Truth:        truth,
	}
}

func hasWordsBeforeVoice(elements []ast.DialogueElement) bool {
	for _, elem := range elements {
		switch el := elem.(type) {
		case *ast.VoiceCommand:
			return false
		case *ast.PlainText:
			if containsLetters(el.Text) {
				return true
			}
		case *ast.FormatTag:
			if result := hasWordsBeforeVoice(el.Content); result {
				return true
			}
		}
	}
	return false
}

func containsLetters(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) {
			return true
		}
	}
	return false
}

func buildAudioTextMap(elements []ast.DialogueElement) map[string][]ast.DialogueElement {
	result := make(map[string][]ast.DialogueElement)
	var currentAudioID string
	var currentFragment []ast.DialogueElement

	var walk func(elems []ast.DialogueElement)
	walk = func(elems []ast.DialogueElement) {
		for _, elem := range elems {
			switch el := elem.(type) {
			case *ast.VoiceCommand:
				if currentAudioID != "" && len(currentFragment) > 0 {
					result[currentAudioID] = currentFragment
				}
				currentAudioID = el.AudioID
				currentFragment = nil
			case *ast.FormatTag:
				if containsVoiceCommand(el.Content) {
					walk(el.Content)
				} else if currentAudioID != "" {
					currentFragment = append(currentFragment, el)
				}
			default:
				if currentAudioID != "" {
					currentFragment = append(currentFragment, elem)
				}
			}
		}
	}
	walk(elements)

	if currentAudioID != "" && len(currentFragment) > 0 {
		result[currentAudioID] = currentFragment
	}

	return result
}

func containsVoiceCommand(elements []ast.DialogueElement) bool {
	for _, elem := range elements {
		switch el := elem.(type) {
		case *ast.VoiceCommand:
			return true
		case *ast.FormatTag:
			if containsVoiceCommand(el.Content) {
				return true
			}
		}
	}
	return false
}

func (e *QuoteExtractor) Presets() *transformer.PresetContext {
	return e.presets
}

func (e *QuoteExtractor) ResolveSoundEffects(effects []pendingSoundEffect) []dto.SoundEffect {
	if len(effects) == 0 || len(e.seFileMap) == 0 {
		return nil
	}

	type seKey struct {
		filename  string
		afterClip int
	}
	seen := make(map[seKey]bool)
	var result []dto.SoundEffect
	for _, se := range effects {
		filename, ok := e.seFileMap[se.seNum]
		if !ok {
			continue
		}
		key := seKey{filename, se.afterClip}
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, dto.SoundEffect{
			Filename:  filename,
			AfterClip: se.afterClip,
		})
	}
	return result
}

func (*QuoteExtractor) parseOmakeEpisode(name string) (int, bool) {
	if len(name) < 3 || name[0] != 'o' {
		return 0, false
	}
	idx := strings.IndexByte(name[1:], '_')
	if idx <= 0 {
		return 0, false
	}
	ep, err := strconv.Atoi(name[1 : 1+idx])
	if err != nil {
		return 0, false
	}
	return ep, true
}
