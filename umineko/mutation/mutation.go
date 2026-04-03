package mutation

import (
	"github.com/VictoriqueMoe/umineko_script_parser/dto"
	"github.com/VictoriqueMoe/umineko_script_parser/umineko/mutation/engine"
)

type Engine interface {
	Apply(quotes []dto.UminekoQuote) []dto.UminekoQuote
}

type Pipeline struct {
	engines []Engine
}

func NewPipeline() *Pipeline {
	p := &Pipeline{
		engines: []Engine{
			&engine.KanonAttributionEngine{},
		},
	}
	return p
}

func (p *Pipeline) Apply(quotes []dto.UminekoQuote) []dto.UminekoQuote {
	for _, e := range p.engines {
		quotes = e.Apply(quotes)
	}
	return quotes
}
