package umineko

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	"github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/character"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/decoder"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/lexer"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/mutation"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/transformer"
)

type (
	ParsedQuote = dto.UminekoQuote

	parser struct {
		extractor *lexer.QuoteExtractor
		factory   *transformer.Factory
		mutations mutation.Pipeline
	}
)

func ParseScriptText(script string) ([]ParsedQuote, []lexer.SubtitleRef, []scriptparser.ValidationError, error) {
	if strings.ContainsRune(script, 0) {
		return nil, nil, nil, scriptparser.ErrBinaryInput
	}
	p := newParser()
	quotes, refs, validationErrors := p.parse(strings.Split(script, "\n"))
	return quotes, refs, validationErrors, nil
}

func ParseFile(r io.Reader) ([]ParsedQuote, []lexer.SubtitleRef, []scriptparser.ValidationError, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read: %w", err)
	}
	decoded, err := decoder.Decode(raw)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("decode: %w", err)
	}
	quotes, refs, validationErrors, err := ParseScriptText(string(decoded))
	if err != nil {
		return nil, nil, nil, err
	}
	return quotes, refs, validationErrors, nil
}

func newParser() *parser {
	extractor := lexer.NewQuoteExtractor()

	return &parser{
		extractor: extractor,
		factory:   transformer.NewFactory(extractor.Presets()),
		mutations: *mutation.NewPipeline(),
	}
}

func (p *parser) parse(lines []string) ([]ParsedQuote, []lexer.SubtitleRef, []scriptparser.ValidationError) {
	filtered := make([]string, 0, len(lines)/8)
	seFileMap := make(map[int]string)

	for _, line := range lines {
		if len(line) < 2 {
			continue
		}

		if seNum, filename, ok := lexer.ParseSeDispatchLine(line); ok {
			seFileMap[seNum] = filename
			continue
		}

		switch line[0] {
		case 'd':
			if line[1] == ' ' || (line[1] == '2' && len(line) > 2 && line[2] == ' ') {
				filtered = append(filtered, line)
			}
		case 'p':
			if len(line) > 13 && line[:13] == "preset_define" {
				filtered = append(filtered, line)
			}
		case 'n':
			if len(line) > 4 && line[:4] == "new_" {
				filtered = append(filtered, line)
			}
		case '*':
			filtered = append(filtered, line)
		case 's':
			if len(line) > 8 && line[:8] == "stralias" {
				filtered = append(filtered, line)
			} else if len(line) > 8 && line[:8] == "ssa_load" {
				filtered = append(filtered, line)
			} else if len(line) > 6 && line[:6] == "seplay" {
				filtered = append(filtered, line)
			}
		case 'l':
			if len(line) > 2 && line[:2] == "lv" && line[2] == ' ' {
				filtered = append(filtered, line)
			}
		case 'm':
			if len(line) > 6 && line[:6] == "meplay" {
				filtered = append(filtered, line)
			}
		case 'w':
			if len(line) > 9 && line[:9] == "wait_on_d" {
				filtered = append(filtered, line)
			}
		}
	}

	p.extractor.SetSeFileMap(seFileMap)

	input := strings.Join(filtered, "\n")
	extracted := p.extractor.ExtractQuotes(input)
	quotes := make([]ParsedQuote, len(extracted))

	plainText := p.factory.MustGet(transformer.FormatPlainText)
	htmlText := p.factory.MustGet(transformer.FormatHTML)

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (len(extracted) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(extracted) {
			end = len(extracted)
		}
		if start >= end {
			break
		}
		wg.Go(func() {
			for i := start; i < end; i++ {
				eq := &extracted[i]

				var audioTextMap map[string]string
				if len(eq.AudioTextMap) > 0 {
					audioTextMap = make(map[string]string, len(eq.AudioTextMap))
					for audioID, fragment := range eq.AudioTextMap {
						audioTextMap[audioID] = plainText.Transform(fragment)
					}
				}

				quotes[i] = ParsedQuote{
					BaseQuote: dto.BaseQuote{
						Text:         plainText.Transform(eq.Content),
						TextHtml:     htmlText.Transform(eq.Content),
						CharacterID:  eq.CharacterID,
						Character:    character.CharacterNames.GetCharacterName(character.CharacterFromID(eq.CharacterID)),
						AudioID:      eq.AudioID,
						AudioCharMap: eq.AudioCharMap,
						AudioTextMap: audioTextMap,
						Episode:      eq.Episode,
						ContentType:  eq.ContentType,
						SoundEffects: p.extractor.ResolveSoundEffects(eq.SoundEffects),
					},
					HasRedTruth:    eq.Truth.HasRed,
					HasBlueTruth:   eq.Truth.HasBlue,
					HasGoldTruth:   eq.Truth.HasGold,
					HasPurpleTruth: eq.Truth.HasPurple,
				}
			}
		})
	}
	wg.Wait()

	quotes = p.mutations.Apply(quotes)

	return quotes, p.extractor.SubtitleRefs(), p.extractor.ValidationErrors()
}
