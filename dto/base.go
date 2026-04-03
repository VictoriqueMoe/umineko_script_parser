package dto

type (
	BaseQuote struct {
		Text         string            `json:"text" example:"You're insane!! Stop it already, damn iiiiiiiiit!!"`
		TextHtml     string            `json:"textHtml" example:"You&#39;re insane!! Stop it already, damn iiiiiiiiit!!"`
		CharacterID  string            `json:"characterId" example:"10"`
		Character    string            `json:"character" example:"Ushiromiya Battler"`
		AudioID      string            `json:"audioId" example:"30101088, 30101089"`
		AudioCharMap map[string]string `json:"audioCharMap,omitempty" example:"30101088:10,30101089:10"`
		AudioTextMap map[string]string `json:"audioTextMap,omitempty" example:"30101088:You're insane!!,30101089:Stop it already"`
		Episode      int               `json:"episode" example:"3"`
		ContentType  string            `json:"contentType" example:""`
		SoundEffects []SoundEffect     `json:"soundEffects,omitempty"`
	}

	SoundEffect struct {
		Filename  string `json:"filename" example:"umise_047"`
		AfterClip int    `json:"afterClip" example:"0"`
	}
)
