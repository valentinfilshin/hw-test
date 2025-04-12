package hw02unpackstring

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	counter := 0
	stringLen := utf8.RuneCountInString(str)
	result := make([]rune, 0, stringLen)
	lastRune := rune(0)

	for _, r := range str {
		if counter == 0 && unicode.IsDigit(r) {
			return "", ErrInvalidString
		}

		if counter > 0 {
			if unicode.IsDigit(r) && unicode.IsDigit(lastRune) {
				return "", ErrInvalidString
			}

			if unicode.IsDigit(r) {
				digit := int(r - '0')
				for i := 0; i < digit; i++ {
					result = append(result, lastRune)
				}
			}

			if !unicode.IsDigit(r) && !unicode.IsDigit(lastRune) {
				result = append(result, lastRune)
			}
		}

		counter++
		lastRune = r

		if counter == stringLen && !unicode.IsDigit(lastRune) {
			result = append(result, lastRune)
		}
	}

	return string(result), nil
}
