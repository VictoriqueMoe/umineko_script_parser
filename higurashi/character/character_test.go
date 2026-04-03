package character

import "testing"

func TestCharacterFromID_MainCast(t *testing.T) {
	tests := []struct {
		id   string
		want Character
	}{
		{"1", Keiichi},
		{"2", Rena},
		{"3", Mion},
		{"4", Satoko},
		{"5", Rika},
		{"6", Shion},
		{"8", Tomitake},
		{"9", Takano},
		{"10", Irie},
		{"11", Ooishi},
		{"12", Hanyuu},
		{"13", Akasaka},
		{"0", MiscVoices},
		{"45", Hanyuu},
		{"43", Mion},
	}

	for _, tt := range tests {
		t.Run(tt.id, func(t *testing.T) {
			got := CharacterFromID(tt.id)
			if got != tt.want {
				t.Errorf("CharacterFromID(%q) = %q, want %q", tt.id, got, tt.want)
			}
		})
	}
}

func TestCharacterFromID_Unknown(t *testing.T) {
	got := CharacterFromID("999")
	if got != Character("999") {
		t.Errorf("expected Character(\"999\"), got %q", got)
	}
}

func TestCharacterFromName_MainCast(t *testing.T) {
	tests := []struct {
		name string
		want Character
	}{
		{"Keiichi", Keiichi},
		{"Rena", Rena},
		{"Mion", Mion},
		{"Ooishi", Ooishi},
		{"Adult Mion", Mion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CharacterFromName(tt.name)
			if got != tt.want {
				t.Errorf("CharacterFromName(%q) = %q, want %q", tt.name, got, tt.want)
			}
		})
	}
}

func TestCharacterFromName_MinorCharacter(t *testing.T) {
	got := CharacterFromName("Forensic Investigator")
	if got != Character("forensic investigator") {
		t.Errorf("expected lowercased fallback, got %q", got)
	}
}

func TestCharacterID(t *testing.T) {
	if Keiichi.ID() != "1" {
		t.Errorf("Keiichi.ID() = %q, want \"1\"", Keiichi.ID())
	}
	if Narrator.ID() != "narrator" {
		t.Errorf("Narrator.ID() = %q, want \"narrator\"", Narrator.ID())
	}
}

func TestCharacterNames_GetCharacterName(t *testing.T) {
	got := CharacterNames.GetCharacterName(Keiichi)
	if got != "Maebara Keiichi" {
		t.Errorf("got %q, want \"Maebara Keiichi\"", got)
	}
}

func TestCharacterNames_GetCharacterName_Unknown(t *testing.T) {
	got := CharacterNames.GetCharacterName(Character("forensic investigator"))
	if got != "forensic investigator" {
		t.Errorf("got %q, want \"forensic investigator\"", got)
	}
}

func TestGetAllCharacters_ReturnsCopy(t *testing.T) {
	all := CharacterNames.GetAllCharacters()
	all[Keiichi] = "modified"
	if CharacterNames[Keiichi] == "modified" {
		t.Error("GetAllCharacters should return a copy")
	}
}
