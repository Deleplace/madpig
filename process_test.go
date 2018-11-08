package madpig

import (
	"fmt"
	"io/ioutil"
	"sync"
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

func TestProcessConcurrent(t *testing.T) {
	hits := processConcurrent(t)
	for url, filehits := range hits {
		fmt.Println(url, "contains", filehits)
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
func BenchmarkProcessConcurrent(b *testing.B) {
	for i := 0; i < b.N; i++ {
		processConcurrent(b)
	}
}

func BenchmarkProcessDocuments(b *testing.B) {
	buffers := make([][]byte, len(testfiles))
	for i, filepath := range testfiles {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			b.Fatal(err)
		}
		buffers[i] = data
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, doc := range buffers {
			hits := documentFindWords(doc, words)
			_ = hits
		}
	}
}

func BenchmarkProcessDocumentsConcurrentWords(b *testing.B) {
	buffers := make([][]byte, len(testfiles))
	for i, filepath := range testfiles {
		data, err := ioutil.ReadFile(filepath)
		if err != nil {
			b.Fatal(err)
		}
		buffers[i] = data
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		wg.Add(len(buffers) * len(words))
		for j := range buffers {
			for k := range words {
				doc, word := buffers[j], words[k]
				go func() {
					found := documentContains(doc, word)
					_ = found
					wg.Done()
				}()
			}
		}
		wg.Wait()
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

func processConcurrent(t testing.TB) (allhits map[string][]string) {
	allhits = make(map[string][]string, len(testfiles))
	var m sync.Mutex
	var wg sync.WaitGroup
	for _, url := range pageURLs {
		url := url
		wg.Add(1)
		go func() {
			defer wg.Done()
			hits, err := webpageFindWords(url, words)
			m.Lock()
			defer m.Unlock()
			if err != nil {
				t.Error(err)
				return
			}
			allhits[url] = hits
		}()
	}
	wg.Wait()
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
