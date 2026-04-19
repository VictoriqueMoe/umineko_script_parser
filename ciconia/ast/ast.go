package ast

type (
	PlainText struct {
		Text string
	}

	ColoredText struct {
		Text string
		Hex  string
	}
)

func (p *PlainText) ElementType() string {
	return "plain_text"
}

func (p *PlainText) GetText() string {
	return p.Text
}

func (c *ColoredText) ElementType() string {
	return "colored_text"
}

func (c *ColoredText) GetText() string {
	return c.Text
}
