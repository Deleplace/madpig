package madpig

import (
	"fmt"
	"testing"
)

func TestProcess(t *testing.T) {
	hits := process(t)
	for url, filehits := range hits {
		fmt.Println(url, "contains", filehits)
	}
}

func TestProcessFiles(t *testing.T) {
	hits := processFiles(t)
	for filepath, filehits := range hits {
		fmt.Println(filepath, "contains", filehits)
	}
}

func BenchmarkProcess(b *testing.B) {
	for i := 0; i < b.N; i++ {
		process(b)
	}
}

func BenchmarkProcessFiles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processFiles(b)
	}
}

func process(t testing.TB) (allhits map[string][]string) {
	allhits = make(map[string][]string, len(testfiles))
	for _, url := range pageURLs {
		hits, err := webpageFindWords(url, words)
		if err != nil {
			t.Error(err)
			continue
		}
		allhits[url] = hits
	}
	return allhits
}

func processFiles(t testing.TB) (allhits map[string][]string) {
	allhits = make(map[string][]string, len(testfiles))
	for _, filepath := range testfiles {
		hits, err := fileFindWords(filepath, words)
		if err != nil {
			t.Error(err)
			continue
		}
		allhits[filepath] = hits
	}
	return allhits
}

var pageURLs = []string{
	"https://en.wikipedia.org/wiki/Go_(programming_language)",
	"https://en.wikipedia.org/wiki/Benchmark",
	"https://en.wikipedia.org/wiki/Rickrolling",
}

var testfiles = []string{
	"testdata/Go_(programming_language).html",
	"testdata/Benchmark.html",
	"testdata/Rickrolling.html",
}

var words = []string{
	"measurements",
	"penguin",
	"gopher",
	"fox",
	"concurrency",
	"gromit",
	"deoxyribonucleic",
}
