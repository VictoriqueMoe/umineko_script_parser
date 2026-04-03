package character

import "strings"

type (
	Character        string
	characterMapping map[Character]string
)

const (
	Narrator   Character = "narrator"
	MiscVoices Character = "misc_voices"
	Keiichi    Character = "keiichi"
	Rena       Character = "rena"
	Mion       Character = "mion"
	Satoko     Character = "satoko"
	Rika       Character = "rika"
	Shion      Character = "shion"
	Satoshi    Character = "satoshi"
	Tomitake   Character = "tomitake"
	Takano     Character = "takano"
	Irie       Character = "irie"
	Ooishi     Character = "ooishi"
	Hanyuu     Character = "hanyuu"
	Akasaka    Character = "akasaka"
	Okonogi    Character = "okonogi"
	Kasai      Character = "kasai"
	Kimiyoshi  Character = "kimiyoshi"
	Oryou      Character = "oryou"
	Teppei     Character = "teppei"
	Rina       Character = "rina"
	Chie       Character = "chie"
	Tomoe      Character = "tomoe"
	Madoka     Character = "madoka"
	Yamaoki    Character = "yamaoki"
	Fujita     Character = "fujita"
	Natsumi    Character = "natsumi"
	Chisato    Character = "chisato"
	Tamako     Character = "tamako"
	Akira      Character = "akira"
	Miyuki     Character = "miyuki"
	Otobe      Character = "otobe"
	Towada     Character = "towada"
	Riku       Character = "riku"
	Ouka       Character = "ouka"
	Kumagai    Character = "kumagai"
	Nagisa     Character = "nagisa"
	Akane      Character = "akane"
	Arakawa    Character = "arakawa"
	Maeno      Character = "maeno"
)

var (
	characterIDs = map[Character]string{
		MiscVoices: "0",
		Keiichi:    "1",
		Rena:       "2",
		Mion:       "3",
		Satoko:     "4",
		Rika:       "5",
		Shion:      "6",
		Tomitake:   "8",
		Takano:     "9",
		Irie:       "10",
		Ooishi:     "11",
		Hanyuu:     "12",
		Akasaka:    "13",
		Okonogi:    "14",
		Tomoe:      "28",
		Madoka:     "29",
		Yamaoki:    "30",
		Fujita:     "31",
		Natsumi:    "36",
		Chisato:    "37",
		Tamako:     "38",
		Akira:      "39",
		Miyuki:     "40",
		Otobe:      "41",
		Towada:     "42",
		Riku:       "46",
		Ouka:       "47",
	}

	idToCharacter map[string]Character

	nameToCharacter map[string]Character

	CharacterNames = characterMapping{
		Narrator:   "Narrator",
		MiscVoices: "Misc Voices",
		Keiichi:    "Maebara Keiichi",
		Rena:       "Ryuuguu Rena",
		Mion:       "Sonozaki Mion",
		Satoko:     "Houjou Satoko",
		Rika:       "Furude Rika",
		Shion:      "Sonozaki Shion",
		Satoshi:    "Houjou Satoshi",
		Hanyuu:     "Hanyuu",
		Tomitake:   "Tomitake Jirou",
		Takano:     "Takano Miyo",
		Irie:       "Irie Kyousuke",
		Ooishi:     "Ooishi Kuraudo",
		Akasaka:    "Akasaka Mamoru",
		Okonogi:    "Okonogi Tetsurou",
		Kasai:      "Kasai Tatsuyoshi",
		Kimiyoshi:  "Kimiyoshi Kiichiro",
		Oryou:      "Sonozaki Oryou",
		Teppei:     "Houjou Teppei",
		Rina:       "Mamiya Rina",
		Chie:       "Rumiko Chie",
		Tomoe:      "Minai Tomoe",
		Madoka:     "Ootaka Madoka",
		Yamaoki:    "Yamaoki",
		Fujita:     "Fujita",
		Natsumi:    "Kimiyoshi Natsumi",
		Chisato:    "Saeki Chisato",
		Tamako:     "Makimura Tamako",
		Akira:      "Toudou Akira",
		Miyuki:     "Akasaka Miyuki",
		Otobe:      "Otobe Yusuke",
		Towada:     "Towada Takuma",
		Riku:       "Furude Riku",
		Ouka:       "Furude Ouka",
		Kumagai:    "Kumagai Katsuya",
		Nagisa:     "Ozaki Nagisa",
		Akane:      "Sonozaki Akane",
		Arakawa:    "Arakawa",
		Maeno:      "Maeno",
	}

	advNameMap = map[string]Character{
		"keiichi":    Keiichi,
		"rena":       Rena,
		"mion":       Mion,
		"satoko":     Satoko,
		"rika":       Rika,
		"shion":      Shion,
		"satoshi":    Satoshi,
		"hanyuu":     Hanyuu,
		"tomitake":   Tomitake,
		"takano":     Takano,
		"irie":       Irie,
		"ooishi":     Ooishi,
		"akasaka":    Akasaka,
		"okonogi":    Okonogi,
		"kasai":      Kasai,
		"kimiyoshi":  Kimiyoshi,
		"oryou":      Oryou,
		"teppei":     Teppei,
		"rina":       Rina,
		"chie":       Chie,
		"tomoe":      Tomoe,
		"madoka":     Madoka,
		"yamaoki":    Yamaoki,
		"fujita":     Fujita,
		"natsumi":    Natsumi,
		"chisato":    Chisato,
		"tamako":     Tamako,
		"akira":      Akira,
		"miyuki":     Miyuki,
		"otobe":      Otobe,
		"towada":     Towada,
		"riku":       Riku,
		"ouka":       Ouka,
		"kumagai":    Kumagai,
		"nagisa":     Nagisa,
		"akane":      Akane,
		"arakawa":    Arakawa,
		"maeno":      Maeno,
		"narrator":   Narrator,
		"adult mion": Mion,
	}
)

func init() {
	idToCharacter = make(map[string]Character, len(characterIDs))
	for char, id := range characterIDs {
		idToCharacter[id] = char
	}
	idToCharacter["45"] = Hanyuu
	idToCharacter["43"] = Mion
	idToCharacter["narrator"] = Narrator
}

func (c Character) ID() string {
	if id, ok := characterIDs[c]; ok {
		return id
	}
	return string(c)
}

func CharacterFromID(id string) Character {
	if c, ok := idToCharacter[id]; ok {
		return c
	}
	return Character(id)
}

func CharacterFromName(name string) Character {
	lower := strings.ToLower(name)
	if c, ok := advNameMap[lower]; ok {
		return c
	}
	return Character(lower)
}

func (c characterMapping) GetCharacterName(ch Character) string {
	if name, ok := c[ch]; ok {
		return name
	}
	return string(ch)
}

func (c characterMapping) GetAllCharacters() map[Character]string {
	out := make(map[Character]string, len(c))
	for k, v := range c {
		out[k] = v
	}
	return out
}
