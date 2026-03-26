package lexer

import (
	"github.com/VictoriqueMoe/umineko_script_parser/lexer/ast"
	"github.com/VictoriqueMoe/umineko_script_parser/lexer/transformer"
)

type TruthFlags struct {
	HasRed    bool
	HasBlue   bool
	HasGold   bool
	HasPurple bool
}

func DetectTruth(elements []ast.DialogueElement, presets *transformer.PresetContext) TruthFlags {
	var flags TruthFlags
	detectInElements(elements, presets, &flags)
	return flags
}

func detectInElements(elements []ast.DialogueElement, presets *transformer.PresetContext, flags *TruthFlags) {
	for _, elem := range elements {
		if tag, ok := elem.(*ast.FormatTag); ok {
			if tag.Name == "p" || tag.Name == "preset" {
				switch presets.GetSemanticClass(tag.Param) {
				case "red-truth":
					flags.HasRed = true
				case "blue-truth":
					flags.HasBlue = true
				case "gold-truth":
					flags.HasGold = true
				case "purple-truth":
					flags.HasPurple = true
				}
			}
			detectInElements(tag.Content, presets, flags)
		}
	}
}
