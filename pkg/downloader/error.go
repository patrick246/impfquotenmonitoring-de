package downloader

import (
	"fmt"
	"strconv"
)

type DownloadErrorStatusCode struct {
	statusCode int
}

func (e *DownloadErrorStatusCode) Error() string {
	return "status code mismatch, 200 != " + strconv.Itoa(e.statusCode)
}

func (e *DownloadErrorStatusCode) Code() int {
	return e.statusCode
}

type DownloadErrorContentType struct {
	contentType string
}

func (e *DownloadErrorContentType) Error() string {
	return fmt.Sprintf("content type mismatch, expected '%s', got '%s'", xlsxMimeType, e.contentType)
}
