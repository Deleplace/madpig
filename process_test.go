package madpig

import (
	"fmt"
	"testing"
)

func TestProcess(t *testing.T) {
	hits := process(t)
	fmt.Println(hits)
}

func TestProcessFiles(t *testing.T) {
	hits := processFiles(t)
	fmt.Println(hits)
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

func process(t testing.TB) (allhits []string) {
	for _, url := range pageURLs {
		hits, err := webpageFindWords(url, words)
		if err != nil {
			t.Error(err)
			continue
		}
		allhits = append(allhits, hits...)
	}
	return allhits
}

func processFiles(t testing.TB) (allhits []string) {
	for _, filepath := range testfiles {
		hits, err := fileFindWords(filepath, words)
		if err != nil {
			t.Error(err)
			continue
		}
		allhits = append(allhits, hits...)
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
