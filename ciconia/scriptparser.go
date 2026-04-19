package ciconia

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"strings"

	scriptparser "github.com/VictoriqueMoe/umineko_script_parser"
	"github.com/VictoriqueMoe/umineko_script_parser/ciconia/character"
	citransformer "github.com/VictoriqueMoe/umineko_script_parser/ciconia/transformer"
	"github.com/VictoriqueMoe/umineko_script_parser/dto"
)

type (
	ParsedQuote = dto.CicroniaQuote

	Parser struct{}
)

const phaseEpisode = 1

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) ParseLines(lines []string) ([]ParsedQuote, []scriptparser.ValidationError) {
	raw := parse(lines)
	factory := citransformer.NewFactory()
	return buildQuotes(raw, factory), nil
}

func buildQuotes(raw []rawQuote, factory *citransformer.Factory) []ParsedQuote {
	plainText := factory.MustGet(citransformer.FormatPlainText)
	htmlText := factory.MustGet(citransformer.FormatHTML)

	quotes := make([]ParsedQuote, 0, len(raw))
	idCounts := make(map[string]int)

	for i := 0; i < len(raw); i++ {
		rq := &raw[i]

		charID := rq.character.ID()
		charName := character.CharacterNames.GetCharacterName(rq.character)
		if charName == string(rq.character) && rq.characterName != "" {
			charName = rq.characterName
		}

		textEN := plainText.Transform(rq.bodyEN)
		textENHtml := htmlText.Transform(rq.bodyEN)
		textJP := plainText.Transform(rq.bodyJP)
		textJPHtml := htmlText.Transform(rq.bodyJP)

		if textEN == "" && textJP == "" {
			continue
		}

		audioID := synthID(rq.chapter, charID, textEN, textJP, idCounts)

		quotes = append(quotes, ParsedQuote{
			BaseQuote: dto.BaseQuote{
				Text:        textEN,
				TextHtml:    textENHtml,
				CharacterID: charID,
				Character:   charName,
				AudioID:     audioID,
				Episode:     phaseEpisode,
				ContentType: rq.contentType,
			},
			TextJP:     textJP,
			TextJPHtml: textJPHtml,
			Chapter:    rq.chapter,
		})
	}

	return quotes
}

func synthID(chapter, charID, textEN, textJP string, counts map[string]int) string {
	prefix := chapterPrefix(chapter)

	h := sha1.New()
	h.Write([]byte(charID))
	h.Write([]byte{'|'})
	h.Write([]byte(textEN))
	h.Write([]byte{'|'})
	h.Write([]byte(textJP))
	sum := h.Sum(nil)
	hash8 := hex.EncodeToString(sum[:4])

	base := prefix + ":" + hash8
	counts[base]++
	n := counts[base]
	if n == 1 {
		return base
	}
	return fmt.Sprintf("%s:%d", base, n)
}

func chapterPrefix(chapter string) string {
	if chapter == "00" {
		return "pro"
	}
	if chapter == "ep" {
		return "ep"
	}
	if strings.HasPrefix(chapter, "df") {
		return chapter
	}
	return "c" + chapter
}
