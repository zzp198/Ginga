package util

import (
	"fmt"
	"io"
	"net/http"
)

func FormatByte(b uint64) string {
	const (
		B = 1 << (10 * iota)
		KB
		MB
		GB
	)

	var value float64
	var unit string

	switch {
	case b >= GB:
		value = float64(b) / GB
		unit = "GB"
	case b >= MB:
		value = float64(b) / MB
		unit = "MB"
	case b >= KB:
		value = float64(b) / KB
		unit = "KB"
	case b >= B:
		value = float64(b) / B
		unit = "B"
	default:
		value = float64(b)
		unit = "B"
	}

	return fmt.Sprintf("%.2f %s", value, unit)
}

func HttpGet(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	return string(body), nil
}
