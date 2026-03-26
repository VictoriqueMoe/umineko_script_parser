package transformer

import (
	"strconv"
	"strings"

	"github.com/VictoriqueMoe/umineko_script_parser/lexer/ast"
)

type PresetContext struct {
	SemanticPresets map[string]string
	DynamicColours  map[string]string
}

func DefaultSemanticPresets() map[string]string {
	return map[string]string{
		"1":  "red-truth",
		"2":  "blue-truth",
		"41": "gold-truth",
		"42": "purple-truth",
	}
}

func DefaultDynamicColours() map[string]string {
	return map[string]string{}
}

func NewPresetContext() *PresetContext {
	return &PresetContext{
		SemanticPresets: DefaultSemanticPresets(),
		DynamicColours:  DefaultDynamicColours(),
	}
}

func (p *PresetContext) CollectFromScript(script *ast.Script) {
	p.DynamicColours = DefaultDynamicColours()

	for _, line := range script.Lines {
		if preset, ok := line.(*ast.PresetDefineLine); ok {
			presetID := strconv.Itoa(preset.ID)

			if _, isSemantic := p.SemanticPresets[presetID]; isSemantic {
				continue
			}

			colour := strings.ToUpper(preset.Colour)
			if colour == "#FFFFFF" || colour == "" {
				continue
			}

			p.DynamicColours[presetID] = colour
		}
	}
}

func (p *PresetContext) GetSemanticClass(presetID string) string {
	return p.SemanticPresets[presetID]
}

func (p *PresetContext) GetDynamicColour(presetID string) string {
	return p.DynamicColours[presetID]
}
