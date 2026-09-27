// Package timecode converts human-readable timestamps into durations.
package timecode

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// ErrInvalid reports a timestamp outside the accepted formats.
var ErrInvalid = errors.New("invalid timecode")

// Parse accepts SS, MM:SS or HH:MM:SS; the seconds field may be fractional.
// Validation is strict so that a typo fails here, with a clear message,
// instead of silently producing the wrong clip.
func Parse(s string) (time.Duration, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("%w: empty value", ErrInvalid)
	}
	parts := strings.Split(s, ":")
	if len(parts) > 3 {
		return 0, fmt.Errorf("%w: %q has too many fields", ErrInvalid, s)
	}

	last := parts[len(parts)-1]
	if !isDecimal(last) {
		return 0, fmt.Errorf("%w: %q", ErrInvalid, s)
	}
	seconds, err := strconv.ParseFloat(last, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %q", ErrInvalid, s)
	}
	// The 60 limit applies only when a higher-order field exists:
	// "90" means ninety seconds, but "1:90" is ambiguous and is rejected.
	if len(parts) > 1 && seconds >= 60 {
		return 0, fmt.Errorf("%w: %q has seconds out of range", ErrInvalid, s)
	}

	total := seconds
	multiplier := 60.0
	for i := len(parts) - 2; i >= 0; i-- {
		value, err := parseField(parts[i])
		if err != nil {
			return 0, fmt.Errorf("%w: %q", ErrInvalid, s)
		}
		// Minutes are capped only when hours are present; "90:00" is ninety minutes.
		if len(parts) == 3 && i == 1 && value >= 60 {
			return 0, fmt.Errorf("%w: %q has minutes out of range", ErrInvalid, s)
		}
		total += float64(value) * multiplier
		multiplier *= 60
	}

	millis := math.Round(total * 1000)
	return time.Duration(millis) * time.Millisecond, nil
}

// Seconds formats the duration as decimal seconds, a format yt-dlp accepts.
func Seconds(d time.Duration) string {
	return strconv.FormatFloat(d.Seconds(), 'f', -1, 64)
}

func parseField(s string) (int, error) {
	if s == "" || strings.Trim(s, "0123456789") != "" {
		return 0, ErrInvalid
	}
	return strconv.Atoi(s)
}

// isDecimal rejects forms ParseFloat would accept, such as "NaN", "Inf" and "1e3".
func isDecimal(s string) bool {
	digits, dots := 0, 0
	for _, r := range s {
		switch {
		case r >= '0' && r <= '9':
			digits++
		case r == '.':
			dots++
		default:
			return false
		}
	}
	return digits > 0 && dots <= 1
}
