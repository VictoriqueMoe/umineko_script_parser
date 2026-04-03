package dto

type HigurashiQuote struct {
	BaseQuote
	TextJP     string `json:"textJp"`
	TextJPHtml string `json:"textJpHtml"`
	Arc        string `json:"arc" example:"onikakushi"`
}
