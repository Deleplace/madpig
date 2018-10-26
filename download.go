package madpig

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func download(url string) (filepath string, err error) {
	file, err := ioutil.TempFile("", "")
	if err != nil {
		return "", fmt.Errorf("could not create temp file")
	}
	defer file.Close()

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error fetching %q: %v", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("error fetching %q: %s", url, resp.Status)
	}
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error downloading contents of %q: %s", url, resp.Status)
	}
	return file.Name(), nil
}
