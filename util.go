package goytdl

import (
	"errors"
	"net/url"
	"regexp"
	"strings"
)

// ErrNotAYouTubeURL is returned when the input URL can't be recognized as a YouTube URL.
var ErrNotAYouTubeURL = errors.New("not recognized as a YouTube URL")

// ErrUnavailableFormat is returned when the input format itag is not found in the input format stream map.
var ErrUnavailableFormat = errors.New("unavailable video format")

var youtubeDotComURL = regexp.MustCompile("(https?://)?(www\\.)?(youtube\\.com(/watch\\?v=|/watch\\?vi=|/v/|/vi/)|youtu\\.be/)([a-zA-Z0-9_\\-]+)")

// GetVideoIDFromURL retrieves the video ID from a YouTube URL.
// May return ErrNotAYouTubeURL.
func GetVideoIDFromURL(youtubeURL string) (videoID string, err error) {
	if youtubeDotComURL.Match([]byte(youtubeURL)) {
		substrings := youtubeDotComURL.FindStringSubmatch(youtubeURL)
		videoID = substrings[5]
		return
	}
	err = ErrNotAYouTubeURL
	return
}

// GetFmtStreamMap retrieves the stream map from a get_video_info request.
func GetFmtStreamMap(videoInfo url.Values) (fmtStreamMap []string, err error) {
	fmtStreamCSV := videoInfo.Get("url_encoded_fmt_stream_map")
	if err != nil {
		return
	}
	fmtStreamMap = strings.Split(fmtStreamCSV, ",")
	return
}

// GetURLFromFmtItag retrieves the raw video download URL from a format itag and a format stream map.
// For available itags check https://en.wikipedia.org/wiki/YouTube#Quality_and_formats.
func GetURLFromFmtItag(fmtStreamMap []string, itag string) (streamURL string, err error) {
	for _, format := range fmtStreamMap {
		var values url.Values
		values, err = url.ParseQuery(format)
		if err != nil {
			return
		}
		if values.Get("itag") == itag {
			streamURL = values.Get("url")
			return
		}
	}
	err = ErrUnavailableFormat
	return
}
