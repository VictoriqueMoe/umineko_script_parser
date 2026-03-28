package loader

import (
	"io/fs"
	"log"
	"strings"
	"time"

	"github.com/VictoriqueMoe/umineko_script_parser/decoder"
	"github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/lexer"
)

type (
	ParseFunc func(lines []string) ([]dto.ParsedQuote, []lexer.SubtitleRef, []lexer.ValidationError)

	Loader struct {
		fs    fs.ReadFileFS
		parse ParseFunc
	}
)

func New(efs fs.ReadFileFS, parse ParseFunc) *Loader {
	return &Loader{
		fs:    efs,
		parse: parse,
	}
}

func (l *Loader) Load(lang string, path string) ([]dto.ParsedQuote, []lexer.SubtitleRef, []lexer.ValidationError) {
	raw, err := l.fs.ReadFile(path)
	if err != nil {
		log.Printf("[%s] failed to read %s: %v", lang, path, err)
		return nil, nil, nil
	}

	decodeStart := time.Now()
	decoded, err := decoder.Decode(raw)
	if err != nil {
		log.Printf("[%s] failed to decode %s: %v", lang, path, err)
		return nil, nil, nil
	}
	log.Printf("[%s] decoded %s (%d -> %d bytes) in %v", lang, path, len(raw), len(decoded), time.Since(decodeStart).Round(time.Millisecond))

	lines := strings.Split(string(decoded), "\n")

	parseStart := time.Now()
	parsed, subtitleRefs, validationErrors := l.parse(lines)
	log.Printf("[%s] parsed %d lines -> %d quotes in %v", lang, len(lines), len(parsed), time.Since(parseStart).Round(time.Millisecond))

	return parsed, subtitleRefs, validationErrors
}
