package timecode

import (
	"errors"
	"testing"
	"time"
)

func TestParseValid(t *testing.T) {
	cases := map[string]time.Duration{
		"45":          45 * time.Second,
		"0":           0,
		"90":          90 * time.Second,
		"12.5":        12500 * time.Millisecond,
		"1:05":        65 * time.Second,
		"90:00":       90 * time.Minute,
		"00:59.250":   59250 * time.Millisecond,
		"01:02:03":    time.Hour + 2*time.Minute + 3*time.Second,
		" 2:00:00 ":   2 * time.Hour,
		"100:00:00.5": 100*time.Hour + 500*time.Millisecond,
	}
	for input, want := range cases {
		got, err := Parse(input)
		if err != nil {
			t.Errorf("Parse(%q) returned error: %v", input, err)
			continue
		}
		if got != want {
			t.Errorf("Parse(%q) = %v, want %v", input, got, want)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	inputs := []string{
		"", "   ", "abc", "-5", "1e3", "NaN", "Inf", "1.2.3",
		"1:60", "1:60:00", "1:2:3:4", ":30", "1::2", "a:30", "1:-5", ".",
	}
	for _, input := range inputs {
		if _, err := Parse(input); !errors.Is(err, ErrInvalid) {
			t.Errorf("Parse(%q) error = %v, want ErrInvalid", input, err)
		}
	}
}

func TestSeconds(t *testing.T) {
	cases := map[time.Duration]string{
		0:                        "0",
		65 * time.Second:         "65",
		12500 * time.Millisecond: "12.5",
	}
	for input, want := range cases {
		if got := Seconds(input); got != want {
			t.Errorf("Seconds(%v) = %q, want %q", input, got, want)
		}
	}
}
