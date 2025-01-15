package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var unpackedStr strings.Builder
	var lastRune rune
	var isEscapeSymbol bool

	for _, symbol := range str {
		switch {
		case isEscapeSymbol:
			{
				if !isDigit(symbol) && symbol != '\\' {
					return "", ErrInvalidString
				}

				isEscapeSymbol = false
				lastRune = symbol
			}

		case isDigit(symbol):
			{
				if lastRune == 0 {
					return "", ErrInvalidString
				}
				// because each symbol has a code in the ASCI/UTF-8 table, to convert from a numeric rune to an int
				// it is enough to subtract the code of the rune '0' from the code of this rune
				iterationCount := int(symbol - '0')
				for i := 0; i < iterationCount; i++ {
					unpackedStr.WriteRune(lastRune)
				}
				lastRune = 0
			}

		case symbol == '\\':
			isEscapeSymbol = true
			fallthrough

		default:
			if lastRune != 0 {
				unpackedStr.WriteRune(lastRune)
			}
			lastRune = symbol
		}
	}

	if lastRune != 0 {
		unpackedStr.WriteRune(lastRune)
	}
	return unpackedStr.String(), nil
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
