package goytdl

import (
	"io"
	"io/ioutil"
	"net/http"
)

// GetVideoInfo returns the result of youtube.com/get_video_info, using the id parameter as the video_id parameter in the request.
func GetVideoInfo(id string) (responseBody string, err error) {
	resp, err := http.Get("http://youtube.com/get_video_info?video_id=" + id)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	responseBody = string(responseBytes)
	return
}

// Download is a Reader for a download that also does a percentage calculation.
type Download struct {
	io.ReadCloser
	TotalDownloadedBytes int64
	TotalBytes           int64
}

func (dl *Download) Read(p []byte) (n int, err error) {
	n, err = dl.ReadCloser.Read(p)
	if err != nil {
		return
	}
	dl.TotalDownloadedBytes += int64(n)
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
	url, err := GetURLFromItag(fmtStreamMap, itag)
	if err != nil {
		return
	}
	dl, err = GetDownloadFromRawURL(url)
	return
}
