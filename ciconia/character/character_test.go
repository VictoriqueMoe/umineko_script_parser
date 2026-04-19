package character

import "testing"

func TestMainCharacterNames_AllEnumsExistInFullMap(t *testing.T) {
	for ch := range MainCharacterNames {
		if _, ok := CharacterNames[ch]; !ok {
			t.Errorf("%q is in MainCharacterNames but missing from CharacterNames", ch)
		}
	}
}

func TestMainCharacterNames_SubsetOfCharacterNames(t *testing.T) {
	if len(MainCharacterNames) >= len(CharacterNames) {
		t.Errorf("MainCharacterNames should be a strict subset: main=%d, full=%d",
			len(MainCharacterNames), len(CharacterNames))
	}
}

func TestMainCharacterNames_DisplayNamesMatchFullMap(t *testing.T) {
	for ch, mainName := range MainCharacterNames {
		fullName, ok := CharacterNames[ch]
		if !ok {
			continue
		}
		if mainName != fullName {
			t.Errorf("%q display mismatch: main=%q full=%q", ch, mainName, fullName)
		}
	}
}

func TestMainCharacterNames_IncludesNarratorAndKeropoyo(t *testing.T) {
	if _, ok := MainCharacterNames[Narrator]; !ok {
		t.Error("Narrator missing from MainCharacterNames")
	}
	if _, ok := MainCharacterNames[Keropoyo]; !ok {
		t.Error("Keropoyo missing from MainCharacterNames")
	}
}

func TestMainCharacterNames_ExcludesRoleLabels(t *testing.T) {
	roleLabels := []Character{
		Officer, Officers, Reporter, Reporters, Newscaster, Doctor, Captain,
		AOUGeneral, AOUOfficer, ACROfficer, COUColonel, LATOOfficer, ABNOfficer,
		Mom, Dad, Man, Girl, Girls, Radio, News, Announcer,
	}
	for _, role := range roleLabels {
		if _, ok := MainCharacterNames[role]; ok {
			t.Errorf("%q is a role label but appears in MainCharacterNames", role)
		}
	}
}
