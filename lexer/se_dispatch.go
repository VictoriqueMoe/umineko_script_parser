package lexer

import (
	"path/filepath"
	"strconv"
	"strings"
)

func ParseSeDispatchLine(line string) (int, string, bool) {
	trimmed := strings.TrimSpace(line)

	var numPrefix string
	if strings.HasPrefix(trimmed, "if %Se_Number") {
		numPrefix = "if %Se_Number"
	} else if strings.HasPrefix(trimmed, "if %Me_Number") {
		numPrefix = "if %Me_Number"
	} else {
		return 0, "", false
	}

	rest := strings.TrimSpace(trimmed[len(numPrefix):])
	if len(rest) == 0 || rest[0] != '=' {
		return 0, "", false
	}
	rest = strings.TrimSpace(rest[1:])

	numEnd := 0
	for numEnd < len(rest) && rest[numEnd] >= '0' && rest[numEnd] <= '9' {
		numEnd++
	}
	if numEnd == 0 {
		return 0, "", false
	}
	seNum, err := strconv.Atoi(rest[:numEnd])
	if err != nil {
		return 0, "", false
	}

	pathStart := strings.Index(rest, `"sound\se\`)
	if pathStart == -1 {
		return 0, "", false
	}
	pathStart += len(`"sound\se\`)
	pathEnd := strings.Index(rest[pathStart:], `"`)
	if pathEnd == -1 {
		return 0, "", false
	}

	fullFilename := rest[pathStart : pathStart+pathEnd]
	baseName := strings.TrimSuffix(filepath.Base(fullFilename), filepath.Ext(fullFilename))

	return seNum, baseName, true
}
