package hw02unpackstring

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/rivo/uniseg"
)

var (
	ErrInvalidString = errors.New("invalid string")
	ErrUnpackString  = errors.New("unpack string error")
)

func Unpack(str string) (string, error) {
	var unpackedStr strings.Builder
	var lastSymbol string
	var isEscapeSymbol bool

	gr := uniseg.NewGraphemes(str)
	for gr.Next() {
		symbol := gr.Str()

		switch {
		case isEscapeSymbol:
			isEscapeSymbol = false
			lastSymbol = symbol

		case isDigit(gr.Runes()):
			{
				if lastSymbol == "" {
					return "", ErrInvalidString
				}

				// because each symbol has a code in the ASCI/UTF-8 table, to convert from a numeric rune to an int
				// it is enough to subtract the code of the rune '0' from the code of this rune
				iterationCount := int(gr.Runes()[0] - '0')
				for i := 0; i < iterationCount; i++ {
					_, err := unpackedStr.WriteString(lastSymbol)
					if err != nil {
						return "", fmt.Errorf("%w: %w", ErrUnpackString, err)
					}
				}
				lastSymbol = ""
			}

		case symbol == "\\":
			isEscapeSymbol = true
			fallthrough

		default:
			_, err := unpackedStr.WriteString(lastSymbol)
			if err != nil {
				return "", fmt.Errorf("%w: %w", ErrUnpackString, err)
			}

			lastSymbol = symbol
		}
	}

	unpackedStr.WriteString(lastSymbol)
	return unpackedStr.String(), nil
}

func isDigit(in []rune) bool {
	// digits in Unicode table are in range 48-57; we need to check only first rune
	return unicode.IsDigit(in[0])
}
