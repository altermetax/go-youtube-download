package goytdl

import (
	"errors"
	"net/url"
	"strings"
)

// GetFmtStreamMap retrieves the stream map from a get_video_info request.
func GetFmtStreamMap(videoInfo string) (fmtStreamMap []string, err error) {
	values, err := url.ParseQuery(videoInfo)
	if err != nil {
		return
	}
	fmtStreamCSV := values.Get("url_encoded_fmt_stream_map")
	if err != nil {
		return
	}
	fmtStreamMap = strings.Split(fmtStreamCSV, ",")
	return
}

// GetURLFromItag retrieves the raw video download URL from a format itag and a format stream map.
// For available itags check https://en.wikipedia.org/wiki/YouTube#Quality_and_formats.
func GetURLFromItag(fmtStreamMap []string, itag string) (itagURL string, err error) {
	for _, format := range fmtStreamMap {
		var values url.Values
		values, err = url.ParseQuery(format)
		if err != nil {
			return
		}
		if values.Get("itag") == itag {
			itagURL = values.Get("url")
			return
		}
	}
	err = errors.New("itag not found in the video streams map")
	return
}
