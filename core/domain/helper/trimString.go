package helper

import "unicode"

func TrimString(
	initialParam string,
) string {
	trimmedParam := make([]rune, 0, len(initialParam))

	for _, r := range initialParam {
		if !unicode.IsSpace(r) {
			trimmedParam = append(trimmedParam, r)
		}
	}

	return string(trimmedParam)
}
