package dto

type CicroniaQuote struct {
	BaseQuote
	TextJP     string `json:"textJp"`
	TextJPHtml string `json:"textJpHtml"`
	Chapter    string `json:"chapter" example:"01"`
}
