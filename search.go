package madpig

import (
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

func webpageFindWords(url string, words []string) (hits []string, err error) {
	tmpfile, err := download(url)
	if err != nil {
		return nil, err
	}
	return fileFindWords(tmpfile, words)
}

func fileFindWords(filepath string, words []string) (hits []string, err error) {
	for _, word := range words {
		found, err := fileContains(filepath, word)
		if err != nil {
			return hits, err
		}
		if found {
			hits = append(hits, word)
		}
	}
	return hits, nil
}

func documentFindWords(doc []byte, words []string) (hits []string) {
	for _, word := range words {
		found := documentContains(doc, word)
		if found {
			hits = append(hits, word)
		}
	}
	return hits
}

func fileContains(filepath string, word string) (bool, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return false, err
	}
	return documentContains(data, word), nil
}

func documentContains(doc []byte, word string) bool {
	wordBytes := []byte(word)
	return bytes.Contains(doc, wordBytes)
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
