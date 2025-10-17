package stringUtil

import (
	"math/rand"
	"net"
	"net/url"
	"regexp"
	"strings"

	timeUtil "radius-server/src/utils/time"
)

const (
	FullAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	HalfAlphabet = "abcdefghijklmnopqrstuvwxyz0123456789"
	Alphabet     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Digits       = "0123456789"
)

func GenerateRandomString(symbols string, length int) string {
	src := rand.NewSource(timeUtil.NowUnixTime())
	r := rand.New(src)
	b := make([]byte, length)
	for i := range b {
		b[i] = symbols[r.Intn(len(symbols))]
	}
	return string(b)
}

var GetRandomCode = func(codeLength int) string {
	return GenerateRandomString(FullAlphabet, codeLength)
}

func ContainSubString(str string, substr string) bool {
	return strings.Contains(str, substr)
}

func ContainSubStringWithoutCase(str string, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

func IsValidIPAddress(ip string) bool {
	address := net.ParseIP(ip)
	return address != nil
}

func CapitalizeFirstChar(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

var youtubeVideoIdRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{11}$`)

func IsValidYoutubeUrl(urlStr string) bool {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Check if it's a YouTube URL
	host := parsedUrl.Host
	if host == "www.youtube.com" || host == "youtube.com" {
		// Format: https://www.youtube.com/watch?v=VIDEO_ID
		videoId := parsedUrl.Query().Get("v")
		return youtubeVideoIdRegex.MatchString(videoId)
	} else if host == "youtu.be" {
		// Format: https://youtu.be/VIDEO_ID
		path := strings.TrimPrefix(parsedUrl.Path, "/")
		return youtubeVideoIdRegex.MatchString(path)
	}

	return false
}

func ExtractYoutubeVideoId(urlStr string) (string, bool) {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		urlStr = "https://" + urlStr
	}
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return "", false
	}

	host := parsedUrl.Host
	if host == "www.youtube.com" || host == "youtube.com" {
		videoId := parsedUrl.Query().Get("v")
		if youtubeVideoIdRegex.MatchString(videoId) {
			return videoId, true
		}
	} else if host == "youtu.be" {
		path := strings.TrimPrefix(parsedUrl.Path, "/")
		if youtubeVideoIdRegex.MatchString(path) {
			return path, true
		}
	}

	return "", false
}
