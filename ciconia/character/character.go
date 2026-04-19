package character

import "strings"

type (
	Character        string
	characterMapping map[Character]string
)

const (
	Narrator                  Character = "narrator"
	Keropoyo                  Character = "keropoyo"
	Miyao                     Character = "miyao"
	Jayden                    Character = "jayden"
	Gunhild                   Character = "gunhild"
	Toujirou                  Character = "toujirou"
	Chloe                     Character = "chloe"
	Rethabile                 Character = "rethabile"
	Koshka                    Character = "koshka"
	Lilja                     Character = "lilja"
	Lingji                    Character = "lingji"
	Stanislaw                 Character = "stanislaw"
	Meow                      Character = "meow"
	Seshat                    Character = "seshat"
	Okonogi                   Character = "okonogi"
	Maricarmen                Character = "maricarmen"
	Naima                     Character = "naima"
	Jestress                  Character = "jestress"
	Valentina                 Character = "valentina"
	Aysha                     Character = "aysha"
	Vier                      Character = "vier"
	Momotake                  Character = "momotake"
	Naomi                     Character = "naomi"
	Ishak                     Character = "ishak"
	Andry                     Character = "andry"
	Abdou                     Character = "abdou"
	Sujatha                   Character = "sujatha"
	Rukhshana                 Character = "rukhshana"
	Owner                     Character = "owner"
	Leah                      Character = "leah"
	Stephania                 Character = "stephania"
	KingOfSorrow              Character = "king_of_sorrow"
	Fatma                     Character = "fatma"
	Commentator               Character = "commentator"
	News                      Character = "news"
	Noor                      Character = "noor"
	StaffMember               Character = "staff_member"
	Mariana                   Character = "mariana"
	KingOfRidicule            Character = "king_of_ridicule"
	Gannet                    Character = "gannet"
	Captain                   Character = "captain"
	SeniorMilitaryOfficial    Character = "senior_military_official"
	Officer                   Character = "officer"
	KingOfFury                Character = "king_of_fury"
	Doctor                    Character = "doctor"
	Mario                     Character = "mario"
	Man                       Character = "man"
	COUColonel                Character = "cou_colonel"
	CheeryMan                 Character = "cheery_man"
	Instructor                Character = "instructor"
	Attendant                 Character = "attendant"
	COUOfficer                Character = "cou_officer"
	Reporter                  Character = "reporter"
	Engineers                 Character = "engineers"
	Setsugekka                Character = "setsugekka"
	Gray                      Character = "gray"
	GovernmentOfficial        Character = "government_official"
	AOUHighRankingOfficer     Character = "aou_high_ranking_officer"
	Analyst                   Character = "analyst"
	SuitedMan                 Character = "suited_man"
	AOUGeneral                Character = "aou_general"
	LATOOfficer               Character = "lato_officer"
	Announcer                 Character = "announcer"
	ACROfficer                Character = "acr_officer"
	ViceAdmiral               Character = "vice_admiral"
	Operator                  Character = "operator"
	Friends                   Character = "friends"
	AOUOperator               Character = "aou_operator"
	LiuTiankai                Character = "liu_tiankai"
	HighMilitaryOfficials     Character = "high_military_officials"
	AOUOfficer                Character = "aou_officer"
	Trainees                  Character = "trainees"
	Radio                     Character = "radio"
	InternationalVolunteer    Character = "international_volunteer"
	Friend                    Character = "friend"
	Cyril                     Character = "cyril"
	CommandCenterRadio        Character = "command_center_radio"
	ABNOfficer                Character = "abn_officer"
	Wicksell                  Character = "wicksell"
	Students                  Character = "students"
	StrategicOfficer          Character = "strategic_officer"
	Spectator                 Character = "spectator"
	Officers                  Character = "officers"
	Newscaster                Character = "newscaster"
	Mom                       Character = "mom"
	COUInternetDenizen        Character = "cou_internet_denizen"
	AOUPolitician             Character = "aou_politician"
	ACRLowerRankingOfficer    Character = "acr_lower_ranking_officer"
	YoungPeople               Character = "young_people"
	Waiter                    Character = "waiter"
	Reporters                 Character = "reporters"
	Receptionist              Character = "receptionist"
	NewsBulletin              Character = "news_bulletin"
	MilitaryRadio             Character = "military_radio"
	Kings                     Character = "kings"
	Girls                     Character = "girls"
	Girl                      Character = "girl"
	FlightAttendant           Character = "flight_attendant"
	COUPressSecretary         Character = "cou_press_secretary"
	CommandCenter             Character = "command_center"
	Chamberlain               Character = "chamberlain"
	BlackSuitedMen            Character = "black_suited_men"
	BabaYaga                  Character = "baba_yaga"
	AOUPressSecretary         Character = "aou_press_secretary"
	AOUAnalyst                Character = "aou_analyst"
	ABNStrategicOfficer       Character = "abn_strategic_officer"
	Youth                     Character = "youth"
	VoiceOfTheCityFolk        Character = "voice_of_the_city_folk"
	VideoUploader             Character = "video_uploader"
	Tripleburger              Character = "tripleburger"
	Trainee                   Character = "trainee"
	Secretary                 Character = "secretary"
	Sakurafubuki              Character = "sakurafubuki"
	ProvisionalDirector       Character = "provisional_director"
	PressSecretary            Character = "press_secretary"
	PresidentSSpeech          Character = "president_s_speech"
	PoliceChief               Character = "police_chief"
	MysteriousVoice           Character = "mysterious_voice"
	MenAlongTheWalls          Character = "men_along_the_walls"
	LATOPressSecretary        Character = "lato_press_secretary"
	LATOPolice                Character = "lato_police"
	Kingfisher                Character = "kingfisher"
	InterruptCM               Character = "interrupt_cm"
	InternetForum             Character = "internet_forum"
	HighSchoolGirl            Character = "high_school_girl"
	Hakurou                   Character = "hakurou"
	GovernmentSpokesperson    Character = "government_spokesperson"
	GirlSVoice                Character = "girl_s_voice"
	Fudou                     Character = "fudou"
	FatmaStephania            Character = "fatma_stephania"
	ElderlyGentleman          Character = "elderly_gentleman"
	Dad                       Character = "dad"
	COUStrategicOfficer       Character = "cou_strategic_officer"
	COUOperator               Character = "cou_operator"
	COUAssemblyPressSecretary Character = "cou_assembly_press_secretary"
	COUAnalyst                Character = "cou_analyst"
	Civilian                  Character = "civilian"
	Citizens                  Character = "citizens"
	Chairperson               Character = "chairperson"
	Chairman                  Character = "chairman"
	Both                      Character = "both"
	Blogger                   Character = "blogger"
	ACRStrategicOfficer       Character = "acr_strategic_officer"
	ACRAnalyst                Character = "acr_analyst"
	ABNSoldier                Character = "abn_soldier"
	ABNPressSecretary         Character = "abn_press_secretary"
	ABNNews                   Character = "abn_news"
)

var (
	CharacterNames = characterMapping{
		Narrator:                  "Narrator",
		Keropoyo:                  "Keropoyo",
		Miyao:                     "Miyao",
		Jayden:                    "Jayden",
		Gunhild:                   "Gunhild",
		Toujirou:                  "Toujirou",
		Chloe:                     "Chloe",
		Rethabile:                 "Rethabile",
		Koshka:                    "Koshka",
		Lilja:                     "Lilja",
		Lingji:                    "Lingji",
		Stanislaw:                 "Stanisław",
		Meow:                      "Meow",
		Seshat:                    "Seshat",
		Okonogi:                   "Okonogi",
		Maricarmen:                "Maricarmen",
		Naima:                     "Naima",
		Jestress:                  "Jestress",
		Valentina:                 "Valentina",
		Aysha:                     "Aysha",
		Vier:                      "Vier",
		Momotake:                  "Momotake",
		Naomi:                     "Naomi",
		Ishak:                     "Ishak",
		Andry:                     "Andry",
		Abdou:                     "Abdou",
		Sujatha:                   "Sujatha",
		Rukhshana:                 "Rukhshana",
		Owner:                     "Owner",
		Leah:                      "Leah",
		Stephania:                 "Stephania",
		KingOfSorrow:              "King of Sorrow",
		Fatma:                     "Fatma",
		Commentator:               "Commentator",
		News:                      "News",
		Noor:                      "Noor",
		StaffMember:               "Staff Member",
		Mariana:                   "Mariana",
		KingOfRidicule:            "King of Ridicule",
		Gannet:                    "Gannet",
		Captain:                   "Captain",
		SeniorMilitaryOfficial:    "Senior Military Official",
		Officer:                   "Officer",
		KingOfFury:                "King of Fury",
		Doctor:                    "Doctor",
		Mario:                     "Mario",
		Man:                       "Man",
		COUColonel:                "COU Colonel",
		CheeryMan:                 "Cheery Man",
		Instructor:                "Instructor",
		Attendant:                 "Attendant",
		COUOfficer:                "COU Officer",
		Reporter:                  "Reporter",
		Engineers:                 "Engineers",
		Setsugekka:                "Setsugekka",
		Gray:                      "Gray",
		GovernmentOfficial:        "Government Official",
		AOUHighRankingOfficer:     "AOU High-Ranking Officer",
		Analyst:                   "Analyst",
		SuitedMan:                 "Suited Man",
		AOUGeneral:                "AOU General",
		LATOOfficer:               "LATO Officer",
		Announcer:                 "Announcer",
		ACROfficer:                "ACR Officer",
		ViceAdmiral:               "Vice-Admiral",
		Operator:                  "Operator",
		Friends:                   "Friends",
		AOUOperator:               "AOU Operator",
		LiuTiankai:                "Liu Tiankai",
		HighMilitaryOfficials:     "High Military Officials",
		AOUOfficer:                "AOU Officer",
		Trainees:                  "Trainees",
		Radio:                     "Radio",
		InternationalVolunteer:    "International Volunteer",
		Friend:                    "Friend",
		Cyril:                     "Cyril",
		CommandCenterRadio:        "Command Center Radio",
		ABNOfficer:                "ABN Officer",
		Wicksell:                  "Wicksell",
		Students:                  "Students",
		StrategicOfficer:          "Strategic Officer",
		Spectator:                 "Spectator",
		Officers:                  "Officers",
		Newscaster:                "Newscaster",
		Mom:                       "Mom",
		COUInternetDenizen:        "COU Internet Denizen",
		AOUPolitician:             "AOU Politician",
		ACRLowerRankingOfficer:    "ACR Lower-Ranking Officer",
		YoungPeople:               "Young People",
		Waiter:                    "Waiter",
		Reporters:                 "Reporters",
		Receptionist:              "Receptionist",
		NewsBulletin:              "News Bulletin",
		MilitaryRadio:             "Military Radio",
		Kings:                     "Kings",
		Girls:                     "Girls",
		Girl:                      "Girl",
		FlightAttendant:           "Flight Attendant",
		COUPressSecretary:         "COU Press Secretary",
		CommandCenter:             "Command Center",
		Chamberlain:               "Chamberlain",
		BlackSuitedMen:            "Black-suited Men",
		BabaYaga:                  "Baba Yaga",
		AOUPressSecretary:         "AOU Press Secretary",
		AOUAnalyst:                "AOU Analyst",
		ABNStrategicOfficer:       "ABN Strategic Officer",
		Youth:                     "Youth",
		VoiceOfTheCityFolk:        "Voice of the City Folk",
		VideoUploader:             "Video Uploader",
		Tripleburger:              "Tripleburger",
		Trainee:                   "Trainee",
		Secretary:                 "Secretary",
		Sakurafubuki:              "Sakurafubuki",
		ProvisionalDirector:       "Provisional Director",
		PressSecretary:            "Press Secretary",
		PresidentSSpeech:          "President's Speech",
		PoliceChief:               "Police Chief",
		MysteriousVoice:           "Mysterious Voice",
		MenAlongTheWalls:          "Men along the walls",
		LATOPressSecretary:        "LATO Press Secretary",
		LATOPolice:                "LATO Police",
		Kingfisher:                "Kingfisher",
		InterruptCM:               "Interrupt CM",
		InternetForum:             "Internet Forum",
		HighSchoolGirl:            "High School Girl",
		Hakurou:                   "Hakurou",
		GovernmentSpokesperson:    "Government Spokesperson",
		GirlSVoice:                "Girl's Voice",
		Fudou:                     "Fudou",
		FatmaStephania:            "Fatma/Stephania",
		ElderlyGentleman:          "Elderly Gentleman",
		Dad:                       "Dad",
		COUStrategicOfficer:       "COU Strategic Officer",
		COUOperator:               "COU Operator",
		COUAssemblyPressSecretary: "COU Assembly Press Secretary",
		COUAnalyst:                "COU Analyst",
		Civilian:                  "Civilian",
		Citizens:                  "Citizens",
		Chairperson:               "Chairperson",
		Chairman:                  "Chairman",
		Both:                      "Both",
		Blogger:                   "Blogger",
		ACRStrategicOfficer:       "ACR Strategic Officer",
		ACRAnalyst:                "ACR Analyst",
		ABNSoldier:                "ABN Soldier",
		ABNPressSecretary:         "ABN Press Secretary",
		ABNNews:                   "ABN News",
	}

	MainCharacterNames = characterMapping{
		Narrator:       "Narrator",
		Keropoyo:       "Keropoyo",
		Miyao:          "Miyao",
		Jayden:         "Jayden",
		Gunhild:        "Gunhild",
		Toujirou:       "Toujirou",
		Chloe:          "Chloe",
		Rethabile:      "Rethabile",
		Koshka:         "Koshka",
		Lilja:          "Lilja",
		Lingji:         "Lingji",
		Stanislaw:      "Stanisław",
		Meow:           "Meow",
		Seshat:         "Seshat",
		Okonogi:        "Okonogi",
		Maricarmen:     "Maricarmen",
		Naima:          "Naima",
		Jestress:       "Jestress",
		Valentina:      "Valentina",
		Aysha:          "Aysha",
		Vier:           "Vier",
		Momotake:       "Momotake",
		Naomi:          "Naomi",
		Ishak:          "Ishak",
		Andry:          "Andry",
		Abdou:          "Abdou",
		Sujatha:        "Sujatha",
		Rukhshana:      "Rukhshana",
		Leah:           "Leah",
		Stephania:      "Stephania",
		Fatma:          "Fatma",
		Noor:           "Noor",
		Mariana:        "Mariana",
		KingOfSorrow:   "King of Sorrow",
		KingOfRidicule: "King of Ridicule",
		KingOfFury:     "King of Fury",
		Mario:          "Mario",
		Gannet:         "Gannet",
		Setsugekka:     "Setsugekka",
		Gray:           "Gray",
		LiuTiankai:     "Liu Tiankai",
		Cyril:          "Cyril",
		Wicksell:       "Wicksell",
		BabaYaga:       "Baba Yaga",
		Tripleburger:   "Tripleburger",
		Sakurafubuki:   "Sakurafubuki",
		Kingfisher:     "Kingfisher",
		Fudou:          "Fudou",
		FatmaStephania: "Fatma/Stephania",
		Hakurou:        "Hakurou",
	}

	advNameMap = map[string]Character{
		"narrator":                     Narrator,
		"miyao":                        Miyao,
		"jayden":                       Jayden,
		"gunhild":                      Gunhild,
		"toujirou":                     Toujirou,
		"chloe":                        Chloe,
		"rethabile":                    Rethabile,
		"koshka":                       Koshka,
		"lilja":                        Lilja,
		"lingji":                       Lingji,
		"stanisław":                    Stanislaw,
		"meow":                         Meow,
		"seshat":                       Seshat,
		"okonogi":                      Okonogi,
		"maricarmen":                   Maricarmen,
		"naima":                        Naima,
		"jestress":                     Jestress,
		"valentina":                    Valentina,
		"aysha":                        Aysha,
		"vier":                         Vier,
		"momotake":                     Momotake,
		"naomi":                        Naomi,
		"ishak":                        Ishak,
		"andry":                        Andry,
		"abdou":                        Abdou,
		"sujatha":                      Sujatha,
		"rukhshana":                    Rukhshana,
		"owner":                        Owner,
		"leah":                         Leah,
		"stephania":                    Stephania,
		"king of sorrow":               KingOfSorrow,
		"fatma":                        Fatma,
		"commentator":                  Commentator,
		"news":                         News,
		"noor":                         Noor,
		"staff member":                 StaffMember,
		"mariana":                      Mariana,
		"king of ridicule":             KingOfRidicule,
		"gannet":                       Gannet,
		"captain":                      Captain,
		"senior military official":     SeniorMilitaryOfficial,
		"officer":                      Officer,
		"king of fury":                 KingOfFury,
		"doctor":                       Doctor,
		"mario":                        Mario,
		"man":                          Man,
		"cou colonel":                  COUColonel,
		"cheery man":                   CheeryMan,
		"instructor":                   Instructor,
		"attendant":                    Attendant,
		"cou officer":                  COUOfficer,
		"reporter":                     Reporter,
		"engineers":                    Engineers,
		"setsugekka":                   Setsugekka,
		"gray":                         Gray,
		"government official":          GovernmentOfficial,
		"aou high-ranking officer":     AOUHighRankingOfficer,
		"analyst":                      Analyst,
		"suited man":                   SuitedMan,
		"aou general":                  AOUGeneral,
		"lato officer":                 LATOOfficer,
		"announcer":                    Announcer,
		"acr officer":                  ACROfficer,
		"vice-admiral":                 ViceAdmiral,
		"operator":                     Operator,
		"friends":                      Friends,
		"aou operator":                 AOUOperator,
		"liu tiankai":                  LiuTiankai,
		"high military officials":      HighMilitaryOfficials,
		"aou officer":                  AOUOfficer,
		"trainees":                     Trainees,
		"radio":                        Radio,
		"international volunteer":      InternationalVolunteer,
		"friend":                       Friend,
		"cyril":                        Cyril,
		"command center radio":         CommandCenterRadio,
		"abn officer":                  ABNOfficer,
		"wicksell":                     Wicksell,
		"students":                     Students,
		"strategic officer":            StrategicOfficer,
		"spectator":                    Spectator,
		"officers":                     Officers,
		"newscaster":                   Newscaster,
		"mom":                          Mom,
		"cou internet denizen":         COUInternetDenizen,
		"aou politician":               AOUPolitician,
		"acr lower-ranking officer":    ACRLowerRankingOfficer,
		"young people":                 YoungPeople,
		"waiter":                       Waiter,
		"reporters":                    Reporters,
		"receptionist":                 Receptionist,
		"news bulletin":                NewsBulletin,
		"military radio":               MilitaryRadio,
		"kings":                        Kings,
		"girls":                        Girls,
		"girl":                         Girl,
		"flight attendant":             FlightAttendant,
		"cou press secretary":          COUPressSecretary,
		"command center":               CommandCenter,
		"chamberlain":                  Chamberlain,
		"black-suited men":             BlackSuitedMen,
		"baba yaga":                    BabaYaga,
		"aou press secretary":          AOUPressSecretary,
		"aou analyst":                  AOUAnalyst,
		"abn strategic officer":        ABNStrategicOfficer,
		"youth":                        Youth,
		"voice of the city folk":       VoiceOfTheCityFolk,
		"video uploader":               VideoUploader,
		"tripleburger":                 Tripleburger,
		"trainee":                      Trainee,
		"secretary":                    Secretary,
		"sakurafubuki":                 Sakurafubuki,
		"provisional director":         ProvisionalDirector,
		"press secretary":              PressSecretary,
		"president's speech":           PresidentSSpeech,
		"police chief":                 PoliceChief,
		"mysterious voice":             MysteriousVoice,
		"men along the walls":          MenAlongTheWalls,
		"lato press secretary":         LATOPressSecretary,
		"lato police":                  LATOPolice,
		"kingfisher":                   Kingfisher,
		"interrupt cm":                 InterruptCM,
		"internet forum":               InternetForum,
		"high school girl":             HighSchoolGirl,
		"hakurou":                      Hakurou,
		"government spokesperson":      GovernmentSpokesperson,
		"girl's voice":                 GirlSVoice,
		"fudou":                        Fudou,
		"fatma/stephania":              FatmaStephania,
		"elderly gentleman":            ElderlyGentleman,
		"dad":                          Dad,
		"cou strategic officer":        COUStrategicOfficer,
		"cou operator":                 COUOperator,
		"cou assembly press secretary": COUAssemblyPressSecretary,
		"cou analyst":                  COUAnalyst,
		"civilian":                     Civilian,
		"citizens":                     Citizens,
		"chairperson":                  Chairperson,
		"chairman":                     Chairman,
		"both":                         Both,
		"blogger":                      Blogger,
		"acr strategic officer":        ACRStrategicOfficer,
		"acr analyst":                  ACRAnalyst,
		"abn soldier":                  ABNSoldier,
		"abn press secretary":          ABNPressSecretary,
		"abn news":                     ABNNews,
	}
)

func (c Character) ID() string {
	return string(c)
}

func CharacterFromName(name string) Character {
	lower := strings.ToLower(strings.TrimSpace(name))
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
