package lexer

import "testing"

func TestParseSeDispatchLine(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantNum  int
		wantFile string
		wantOk   bool
	}{
		{
			name:     "standard SE mapping",
			input:    `if %Se_Number = 47 mov $Se_Play,"sound\se\umise_047.ogg" : return`,
			wantNum:  47,
			wantFile: "umise_047",
			wantOk:   true,
		},
		{
			name:     "SE with 4-digit number",
			input:    `if %Se_Number = 1060 mov $Se_Play,"sound\se\umise_1060.ogg" : return`,
			wantNum:  1060,
			wantFile: "umise_1060",
			wantOk:   true,
		},
		{
			name:     "ME mapping to umilse",
			input:    `if %Me_Number =  5   mov $Me_Play,"sound\se\umilse_005.ogg" : return`,
			wantNum:  5,
			wantFile: "umilse_005",
			wantOk:   true,
		},
		{
			name:     "ME mapping to maria_v",
			input:    `if %Me_Number = 1054 mov $Me_Play,"sound\se\maria_v.ogg" : return`,
			wantNum:  1054,
			wantFile: "maria_v",
			wantOk:   true,
		},
		{
			name:     "SE cross-type mapping",
			input:    `if %Se_Number = 29 mov $Se_Play,"sound\se\umilse_022.ogg" : return`,
			wantNum:  29,
			wantFile: "umilse_022",
			wantOk:   true,
		},
		{
			name:     "ME cross-type mapping",
			input:    `if %Me_Number = 1137 mov $Me_Play,"sound\se\umise_037.ogg" : return`,
			wantNum:  1137,
			wantFile: "umise_037",
			wantOk:   true,
		},
		{
			name:     "SE ex file",
			input:    `if %Se_Number = 1064 mov $Se_Play,"sound\se\umise_ex01.ogg" : return`,
			wantNum:  1064,
			wantFile: "umise_ex01",
			wantOk:   true,
		},
		{
			name:     "ME off value",
			input:    `if %Me_Number =  0   mov $Me_Play,"off" : return`,
			wantNum:  0,
			wantFile: "",
			wantOk:   false,
		},
		{
			name:   "unrelated line",
			input:  `d [lv 0*"10"*"10100001"]` + "`Hello`",
			wantOk: false,
		},
		{
			name:   "empty line",
			input:  "",
			wantOk: false,
		},
		{
			name:     "leading whitespace",
			input:    `	if %Se_Number = 1 mov $Se_Play,"sound\se\umise_001.ogg" : return`,
			wantNum:  1,
			wantFile: "umise_001",
			wantOk:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotNum, gotFile, gotOk := ParseSeDispatchLine(tt.input)
			if gotOk != tt.wantOk {
				t.Fatalf("ok = %v, want %v", gotOk, tt.wantOk)
			}
			if !tt.wantOk {
				return
			}
			if gotNum != tt.wantNum {
				t.Errorf("seNum = %d, want %d", gotNum, tt.wantNum)
			}
			if gotFile != tt.wantFile {
				t.Errorf("filename = %q, want %q", gotFile, tt.wantFile)
			}
		})
	}
}
