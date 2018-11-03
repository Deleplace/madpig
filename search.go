package madpig

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func webpageFindWords(url string, words []string) (hits []string, err error) {
	tmpfile, err := download(url)
	if err != nil {
		return nil, err
	}
	for _, word := range words {
		found, err := fileContains(tmpfile, word)
		if err != nil {
			return hits, err
		}
		if found {
			hits = append(hits, fmt.Sprintf("Article %q contains %q :) \n", articleName(url), word))
		}
	}
	return hits, nil
}

func fileContains(filepath string, word string) (bool, error) {
	size, err := filesize(filepath)
	if err != nil {
		return false, err
	}

	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return false, err
	}
	m := len(word)
	for i := 0; i < int(size)-m; i++ {
		if string(data[i:i+m]) == word {
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
