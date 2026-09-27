// Package source validates the video source before handing it to yt-dlp.
package source

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

// ErrUnsupported reports a URL that does not point to a YouTube video.
var ErrUnsupported = errors.New("unsupported URL")

// yt-dlp supports hundreds of sites; restricting hosts keeps the tool within
// its stated scope and prevents arbitrary input from reaching the process.
var allowedHosts = map[string]bool{
	"youtube.com":       true,
	"www.youtube.com":   true,
	"m.youtube.com":     true,
	"music.youtube.com": true,
	"youtu.be":          true,
}

var idPathPrefixes = []string{"/shorts/", "/live/"}

// ValidateYouTube returns the normalized URL or ErrUnsupported.
func ValidateYouTube(raw string) (string, error) {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("%w: %q", ErrUnsupported, raw)
	}
	if u.Scheme != "https" && u.Scheme != "http" {
		return "", fmt.Errorf("%w: %q must use http or https", ErrUnsupported, raw)
	}
	host := strings.ToLower(u.Hostname())
	if !allowedHosts[host] {
		return "", fmt.Errorf("%w: %q is not a YouTube host", ErrUnsupported, raw)
	}
	if !hasVideoID(host, u) {
		return "", fmt.Errorf("%w: %q has no video identifier", ErrUnsupported, raw)
	}
	return u.String(), nil
}

func hasVideoID(host string, u *url.URL) bool {
	if host == "youtu.be" {
		return strings.Trim(u.Path, "/") != ""
	}
	if u.Path == "/watch" {
		return u.Query().Get("v") != ""
	}
	for _, prefix := range idPathPrefixes {
		if strings.HasPrefix(u.Path, prefix) && len(u.Path) > len(prefix) {
			return true
		}
	}
	return false
}
