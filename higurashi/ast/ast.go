package ast

type PlainText struct {
	Text string
}

func (p *PlainText) ElementType() string { return "plain_text" }
func (p *PlainText) GetText() string     { return p.Text }
