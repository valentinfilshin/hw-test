package hw02unpackstring

import (
	"errors"
	"unicode"
	"unicode/utf8"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	stringLen := utf8.RuneCountInString(str)

	if stringLen == 0 {
		return "", nil
	}

	firstRune, _ := utf8.DecodeRuneInString(str)
	if unicode.IsDigit(firstRune) {
		return "", ErrInvalidString
	}

	result := make([]rune, 0, stringLen)
	lastRune := rune(0)
	counter := 0

	for _, r := range str {
		if counter == 0 {
			counter++
			lastRune = r
			continue
		}

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

		counter++
		lastRune = r
	}

	if counter == stringLen && !unicode.IsDigit(lastRune) {
		result = append(result, lastRune)
	}

	return string(result), nil
}
