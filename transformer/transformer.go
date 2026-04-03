package transformer

import "github.com/VictoriqueMoe/umineko_script_parser/dialogue"

type Transformer interface {
	Transform(elements []dialogue.DialogueElement) string
}
