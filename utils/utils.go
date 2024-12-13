package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
)

type patternByte struct {
	isWildcard bool
	value      byte
}

func checkMatch(data []byte, pattern []patternByte) bool {
	for i, pb := range pattern {
		if !pb.isWildcard && data[i] != pb.value {
			return false
		}
	}
	return true
}

func parsePattern(pattern string) ([]patternByte, error) {
	parts := strings.Split(pattern, " ")
	var parsed []patternByte
	for _, part := range parts {
		if part == "?" {
			parsed = append(parsed, patternByte{isWildcard: true})
			continue
		}
		// Assume hex byte
		b, err := hex.DecodeString(part)
		if err != nil || len(b) != 1 {
			return nil, fmt.Errorf("invalid byte in pattern: %s", part)
		}
		parsed = append(parsed, patternByte{isWildcard: false, value: b[0]})
	}
	return parsed, nil
}

func MatchSignature(data []byte, pattern string) (bool, error) {
	patternBytes, err := parsePattern(pattern)
	if err != nil {
		return false, err
	}

	pLen := len(patternBytes)
	dLen := len(data)

	for start := 0; start <= dLen-pLen; start++ {
		if checkMatch(data[start:start+pLen], patternBytes) {
			return true, nil
		}
	}
	return false, nil
}
