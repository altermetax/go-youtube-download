package goytdl

import (
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// ErrVideoInfoHTTPNotOK is returned when YouTube servers return a HTTP status code different than 200.
var ErrVideoInfoHTTPNotOK = errors.New("can't retrieve the video (HTTP error)")

// ErrVideoInfoFail is returned when YouTube servers return a 200 HTTP status code, but an error occurs while retrieving the video.
var ErrVideoInfoFail = errors.New("can't retrieve the video")

// GetVideoInfo returns the result of youtube.com/get_video_info, using the id parameter as the video_id parameter in the request.
// It may return ErrNotFound.
func GetVideoInfo(id string) (responseBody url.Values, err error) {
	resp, err := http.Get("http://youtube.com/get_video_info?video_id=" + id)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		err = ErrVideoInfoHTTPNotOK
		return
	}
	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	responseBody, err = url.ParseQuery(string(responseBytes))
	if err != nil {
		return
	}
	if responseBody.Get("status") == "fail" {
		err = ErrVideoInfoFail
		return
	}
	return
}

// Download is a ReadCloser for a HTTP download that also does a percentage calculation.
type Download struct {
	io.ReadCloser
	TotalDownloadedBytes int64
	TotalBytes           int64
}

// Read implements the Read method of the io.ReadCloser interface.
func (dl *Download) Read(p []byte) (n int, err error) {
	n, err = dl.ReadCloser.Read(p)
	if err != nil {
		return
	}
	dl.TotalDownloadedBytes += int64(n)
	return
}

// Close implements the Close method of the io.ReadCloser interface
func (dl *Download) Close() (err error) {
	err = dl.ReadCloser.Close()
	return
}

// GetPercentage retrieves the current percentage of a download.
func (dl *Download) GetPercentage() float64 {
	return (float64(dl.TotalDownloadedBytes) / float64(dl.TotalBytes)) * float64(100)
}

// GetDownloadFromRawURL sets up a download of the video from a raw URL, i. e. obtained by GetURLFromItag.
// You may want to use NewDownloadFromVideoID instead of this.
// Run Start() on the returned value to start the download.
func GetDownloadFromRawURL(url string) (dl *Download, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	dl = &Download{ReadCloser: resp.Body, TotalBytes: resp.ContentLength}
	return
}

// NewDownloadFromVideoID returns a new Download setup based on a video ID and a format itag value.
// All itag values can be found at https://en.wikipedia.org/wiki/YouTube#Quality_and_formats.
func NewDownloadFromVideoID(id string, itag string) (dl *Download, err error) {
	videoInfo, err := GetVideoInfo(id)
	if err != nil {
		return
	}
	fmtStreamMap, err := GetFmtStreamMap(videoInfo)
	if err != nil {
		return
	}
	url, err := GetURLFromFmtItag(fmtStreamMap, itag)
	if err != nil {
		return
	}
	dl, err = GetDownloadFromRawURL(url)
	return
}
