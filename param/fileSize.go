package param

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"unicode"
)

var suffixToShift = map[rune]uint{
	'K': 10, // kibibyte
	'M': 20, // mebibyte
	'G': 30, // gibibyte
	'T': 40, // tebibyte
	'P': 50, // pebibyte
	'E': 60, // exbibyte
}

func parseHumanReadableSize(humanReadable string) (int64, error) {
	if len(humanReadable) == 0 {
		return 0, errors.New("can't parse human readable size from empty string")
	}
	var scaleString string
	runes := []rune(humanReadable)
	suffix := runes[len(runes)-1]
	shift, ok := suffixToShift[suffix]
	if ok {
		scaleString = string(runes[:len(runes)-1])
	} else if unicode.IsDigit(rune(suffix)) {
		scaleString = humanReadable
		shift = 0
	} else {
		return 0, fmt.Errorf("human readable size can be an integer or float with optional suffix e.g. 1024, 1.5K, 4.2T, but found: %v", humanReadable)
	}
	scaleFloat, err := strconv.ParseFloat(scaleString, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse float from human readable size %v: %w\n", humanReadable, err)
	}
	result := int64(math.Round(scaleFloat * float64(uint(1)<<shift)))

	if result == math.MinInt64 {
		return 0, fmt.Errorf("Couldn't convert due to overflow: %v", humanReadable)
	}

	return result, nil
}
