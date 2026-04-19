package scriptparser

import (
	"fmt"
	"io"
	"strings"

	"github.com/VictoriqueMoe/umineko_script_parser/decoder"
)

type Parser[Q any] interface {
	ParseLines([]string) ([]Q, []ValidationError)
}

func ParseText[Q any](script string, p Parser[Q]) ([]Q, []ValidationError, error) {
	if strings.ContainsRune(script, 0) {
		return nil, nil, ErrBinaryInput
	}
	lines := strings.Split(script, "\n")
	quotes, validationErrors := p.ParseLines(lines)
	return quotes, validationErrors, nil
}

func ParseReader[Q any](r io.Reader, p Parser[Q]) ([]Q, []ValidationError, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, nil, fmt.Errorf("read: %w", err)
	}

	decoded, err := decoder.Decode(data)
	if err != nil {
		return nil, nil, fmt.Errorf("decode: %w", err)
	}

	return ParseText(string(decoded), p)
}
