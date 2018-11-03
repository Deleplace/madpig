package madpig

import (
	"bytes"
	"fmt"
	"os"
	"strings"
)

func webpageFindWords(url string, words []string) (hits []string, err error) {
	for _, word := range words {
		found, err := webpageContains(url, word)
		if err != nil {
			return hits, err
		}
		if found {
			hits = append(hits, fmt.Sprintf("Article %q contains %q :) \n", articleName(url), word))
		}
	}
	return hits, nil
}

func webpageContains(url string, word string) (bool, error) {
	tmpfile, err := download(url)
	if err != nil {
		return false, err
	}
	return fileContains(tmpfile, word)
}

func fileContains(filepath string, word string) (bool, error) {
	wordBytes := []byte(word)
	size, err := filesize(filepath)
	if err != nil {
		return false, err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return false, err
	}
	defer file.Close()
	for i := int64(0); i < size-int64(len(word)); i++ {
		_, err = file.Seek(i, 0)
		if err != nil {
			return false, err
		}

		buffer := make([]byte, len(word))
		_, err := file.Read(buffer)
		if err != nil {
			return false, err
		}
		if bytes.Equal(buffer, wordBytes) {
			return true, nil
		}
	}
	// No position i was a match
	return false, nil
}

func filesize(filepath string) (int64, error) {
	f, err := os.Open(filepath)
	defer f.Close()
	if err != nil {
		return 0, err
	}
	info, err := f.Stat()
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// e.g. "https://en.wikipedia.org/wiki/Go_(programming_language)" -> "Go (programming language)"
func articleName(wikipediaURL string) string {
	parts := strings.Split(wikipediaURL, "/")
	last := parts[len(parts)-1]
	name := strings.Replace(last, "_", " ", -1)
	return name
}
