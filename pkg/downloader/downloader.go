package downloader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const xlsxMimeType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet;charset=UTF-8"

type RkiImpfquotenDownloader struct {
	url string
}

func NewDownloader(url string) (*RkiImpfquotenDownloader, error) {
	return &RkiImpfquotenDownloader{url: url}, nil
}

func (r *RkiImpfquotenDownloader) Download() ([]byte, error) {
	shouldRetry := true
	retryDelay := 100 * time.Millisecond
	const retryDelayMax = 10 * time.Minute

	var response []byte
	var err error

	for shouldRetry {
		response, err = r.tryDownload()
		if err != nil {
			if retryDelay > retryDelayMax {
				return nil, err
			}

			switch e := err.(type) {
			case *DownloadErrorContentType:
				shouldRetry = true
			case *DownloadErrorStatusCode:
				if e.statusCode >= 400 && e.statusCode < 500 {
					shouldRetry = false
				}
				if e.statusCode >= 500 {
					shouldRetry = true
				}
			}
		} else {
			shouldRetry = false
		}

		if shouldRetry {
			time.Sleep(retryDelay)
			retryDelay *= 2
		}
	}

	return response, err
}

func (r *RkiImpfquotenDownloader) tryDownload() ([]byte, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", r.url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("user-agent", "VaccineMonitoring/1.0 (github.com/patrick246/impfquotenmonitoring-de)")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, &DownloadErrorStatusCode{resp.StatusCode}
	}

	if contentType := resp.Header.Get(http.CanonicalHeaderKey("content-type")); contentType != xlsxMimeType {
		return nil, &DownloadErrorContentType{contentType: contentType}
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	return response, nil
}
