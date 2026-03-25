package dto

type ParsedQuote struct {
	Text         string            `json:"text"`
	TextHtml     string            `json:"textHtml"`
	CharacterID  string            `json:"characterId"`
	Character    string            `json:"character"`
	AudioID      string            `json:"audioId"`
	AudioCharMap map[string]string `json:"audioCharMap,omitempty"`
	AudioTextMap map[string]string `json:"audioTextMap,omitempty"`
	Episode      int               `json:"episode"`
	ContentType  string            `json:"contentType"`
	HasRedTruth  bool              `json:"hasRedTruth,omitempty"`
	HasBlueTruth bool              `json:"hasBlueTruth,omitempty"`
}
