package source

import (
	"errors"
	"testing"
)

func TestValidateYouTubeAccepts(t *testing.T) {
	inputs := []string{
		"https://www.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://youtube.com/watch?v=dQw4w9WgXcQ&t=42",
		"http://m.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://music.youtube.com/watch?v=dQw4w9WgXcQ",
		"https://youtu.be/dQw4w9WgXcQ",
		"https://www.youtube.com/shorts/abc123",
		"https://www.youtube.com/live/abc123",
		"  https://WWW.YOUTUBE.COM/watch?v=dQw4w9WgXcQ  ",
	}
	for _, input := range inputs {
		if _, err := ValidateYouTube(input); err != nil {
			t.Errorf("ValidateYouTube(%q) returned error: %v", input, err)
		}
	}
}

func TestValidateYouTubeRejects(t *testing.T) {
	inputs := []string{
		"",
		"%zz",
		"-oProxyCommand=touch /tmp/pwned",
		"ftp://youtube.com/watch?v=abc",
		"https://vimeo.com/123456",
		"https://youtube.com.evil.example/watch?v=abc",
		"https://www.youtube.com/",
		"https://www.youtube.com/watch",
		"https://www.youtube.com/shorts/",
		"https://youtu.be/",
	}
	for _, input := range inputs {
		if _, err := ValidateYouTube(input); !errors.Is(err, ErrUnsupported) {
			t.Errorf("ValidateYouTube(%q) error = %v, want ErrUnsupported", input, err)
		}
	}
}
