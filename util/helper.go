package util

import "strings"

func getValue(str, findS string) (int, string) {
	sLeft := ""
	if str == "" {
		return 0, sLeft
	}

	nPos := 0
	if findS != "" {
		nPos = strings.Index(str, findS)
	}
	if nPos < 0 {
		return 0, sLeft
	}

	sLeft = str[nPos+len(findS):]
	sLeft = strings.TrimLeft(sLeft, " ")
	sLeft = strings.TrimLeft(sLeft, "\r\n")

	nPos = strings.Index(sLeft, "\r\n")
	if nPos >= 0 {
		sLeft = sLeft[:nPos]
	}

	sLeft = strings.TrimRight(sLeft, " ")
	sLeft = strings.TrimRight(sLeft, "\r\n")

	return len(sLeft), sLeft
}

func ContainsOne(str string, substrs ...string) bool {
	for _, substr := range substrs {
		if strings.Contains(str, substr) {
			return true
		}
	}
	return false
}
