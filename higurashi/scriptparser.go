package higurashi

import (
	"fmt"
	"io"
	"runtime"
	"strings"
	"sync"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	"github.com/VictoriqueMoe/umineko_script_parser/dialogue"
	"github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/higurashi/character"
	hitransformer "github.com/VictoriqueMoe/umineko_script_parser/higurashi/transformer"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/decoder"
)

type ParsedQuote = dto.HigurashiQuote

func ParseScriptText(script string) ([]ParsedQuote, []scriptparser.ValidationError, error) {
	if strings.ContainsRune(script, 0) {
		return nil, nil, scriptparser.ErrBinaryInput
	}

	lines := strings.Split(script, "\n")
	raw := parse(lines)

	factory := hitransformer.NewFactory()
	quotes := buildQuotes(raw, factory)

	return quotes, nil, nil
}

func ParseFile(r io.Reader) ([]ParsedQuote, []scriptparser.ValidationError, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("read: %w", err)
	}

	decoded, err := decoder.Decode(data)
	if err != nil {
		return nil, nil, fmt.Errorf("decode: %w", err)
	}

	return ParseScriptText(string(decoded))
}

func buildQuotes(raw []rawQuote, factory *hitransformer.Factory) []ParsedQuote {
	quotes := make([]ParsedQuote, len(raw))

	plainText := factory.MustGet(hitransformer.FormatPlainText)
	htmlText := factory.MustGet(hitransformer.FormatHTML)

	numWorkers := runtime.GOMAXPROCS(0)
	chunkSize := (len(raw) + numWorkers - 1) / numWorkers

	var wg sync.WaitGroup
	for w := 0; w < numWorkers; w++ {
		start := w * chunkSize
		end := start + chunkSize
		if end > len(raw) {
			end = len(raw)
		}
		if start >= end {
			break
		}
		wg.Go(func() {
			for i := start; i < end; i++ {
				quotes[i] = transformQuote(&raw[i], plainText, htmlText)
			}
		})
	}
	wg.Wait()

	return quotes
}

func transformQuote(rq *rawQuote, plainText hitransformer.Transformer, htmlText hitransformer.Transformer) ParsedQuote {
	charID := rq.character.ID()
	charName := rq.characterName
	if charName == "" {
		charName = character.CharacterNames.GetCharacterName(rq.character)
	}

	var (
		allEN        []dialogue.DialogueElement
		allJP        []dialogue.DialogueElement
		voicePaths   []string
		audioCharMap map[string]string
		audioTextMap map[string]string
	)

	if len(rq.segments) > 0 {
		audioCharMap = make(map[string]string, len(rq.segments))
		audioTextMap = make(map[string]string, len(rq.segments))

		for j := 0; j < len(rq.segments); j++ {
			seg := &rq.segments[j]
			voicePaths = append(voicePaths, seg.path)
			audioCharMap[seg.path] = seg.charID
			audioTextMap[seg.path] = plainText.Transform(seg.contentEN)
			allEN = append(allEN, seg.contentEN...)
			allJP = append(allJP, seg.contentJP...)
		}
	}

	allEN = append(allEN, rq.unvoicedEN...)
	allJP = append(allJP, rq.unvoicedJP...)

	audioID := strings.Join(voicePaths, ", ")

	return ParsedQuote{
		BaseQuote: dto.BaseQuote{
			Text:         plainText.Transform(allEN),
			TextHtml:     htmlText.Transform(allEN),
			CharacterID:  charID,
			Character:    charName,
			AudioID:      audioID,
			AudioCharMap: audioCharMap,
			AudioTextMap: audioTextMap,
			Episode:      rq.episode,
			ContentType:  rq.arc,
			SoundEffects: rq.soundEffects,
		},
		TextJP:     plainText.Transform(allJP),
		TextJPHtml: htmlText.Transform(allJP),
		Arc:        rq.arc,
	}
}
