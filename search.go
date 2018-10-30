package madpig

import (
	"fmt"
	"io"
	"os"
)

var count = make(map[string]int)

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
	count["fileContains"]++
	size, err := filesize(filepath)
	if err != nil {
		return false, err
	}

positions:
	for i := int64(0); i < size-int64(len(word)); i++ {
		count["os.Open"]++
		file, err := os.Open(filepath)
		if err != nil {
			return false, err
		}

		_, err = file.Seek(i, 0)
		if err != nil {
			file.Close()
			return false, err
		}

		for j := 0; j < len(word); j++ {
			c1, err := readByte(file)
			if err != nil {
				file.Close()
				return false, err
			}
			c2 := word[j]
			if c1 != c2 {
				// Word was not exactly found at position i
				file.Close()
				continue positions
			}
		}
		// All characters match!!
		file.Close()
		return true, nil
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

func readByte(r io.Reader) (byte, error) {
	buffer := make([]byte, 1)
	_, err := r.Read(buffer)
	return buffer[0], err
}
