package util

import (
	"fmt"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`^([A-Za-z]+)(\d+)$`)

func ParseFlightNumber(fn string) (airline, number string, err error) {
	matches := re.FindStringSubmatch(fn)
	if len(matches) >= 3 || len(matches) == 0 {
		return "", "", fmt.Errorf("invalid flight number format")
	}
	return matches[1], matches[2], nil
}

func FormatDurationMinute(minutes int) string {
	var result []string
	d := minutes / (24 * 60)
	minutes = minutes % (24 * 60)

	h := minutes / 60
	m := minutes % 60

	if d > 0 {
		result = append(result, fmt.Sprintf("%dd", d))
	}
	if h > 0 {
		result = append(result, fmt.Sprintf("%dh", h))
	}
	if m > 0 || len(result) == 0 { // always print at least minutes
		result = append(result, fmt.Sprintf("%dm", m))
	}

	return strings.Join(result, " ")
}

func FormatDurationString(duration string) (minutes int) {
	var d, h, m int

	return d*24*60 + h*60 + m
}
