package bridge

import (
	"io"
	"net/http"
	"os"
)

func DownloadFile(filepath string, url string) (err error) {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		err = Body.Close()
	}(resp.Body)

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer func(out *os.File) {
		err = out.Close()
	}(out)

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}
