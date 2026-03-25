package transformer

import "github.com/VictoriqueMoe/umineko_script_parser/lexer/ast"

type Transformer interface {
	Transform(elements []ast.DialogueElement) string
}
