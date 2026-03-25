package lexer

import (
	"github.com/VictoriqueMoe/umineko_script_parser/lexer/ast"
	"github.com/VictoriqueMoe/umineko_script_parser/lexer/transformer"
)

type TruthFlags struct {
	HasRed  bool
	HasBlue bool
}

func DetectTruth(elements []ast.DialogueElement, presets *transformer.PresetContext) TruthFlags {
	var flags TruthFlags
	detectInElements(elements, presets, &flags.HasRed, &flags.HasBlue)
	return flags
}

func detectInElements(elements []ast.DialogueElement, presets *transformer.PresetContext, hasRed, hasBlue *bool) {
	for _, elem := range elements {
		if tag, ok := elem.(*ast.FormatTag); ok {
			if tag.Name == "p" || tag.Name == "preset" {
				class := presets.GetSemanticClass(tag.Param)
				if class == "red-truth" {
					*hasRed = true
				} else if class == "blue-truth" {
					*hasBlue = true
				}
			}
			detectInElements(tag.Content, presets, hasRed, hasBlue)
		}
	}
}
