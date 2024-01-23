package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"index/suffixarray"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var pagination Pagination

func main() {
	searcher := Searcher{}
	err := searcher.Load("completeworks.txt")
	if err != nil {
		log.Fatal(err)
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        handleSearch(searcher)(w, r)
    })

	http.HandleFunc("/loadMore", func(w http.ResponseWriter, r *http.Request) {
        enableCors(&w)
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        loadMore()(w, r)
    })

	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}

	fmt.Printf("shakesearch available at http://localhost:%s...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatal(err)
	}
}

func enableCors(w *http.ResponseWriter) {
    (*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
}

type Searcher struct {
	CompleteWorks string
	SuffixArray   *suffixarray.Index
}

type Pagination struct {
	data   []string
}

func handleSearch(searcher Searcher) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		query, ok := r.URL.Query()["q"]
		if !ok || len(query[0]) < 1 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("missing search query in URL params"))
			return
		}
		results := searcher.Search(query[0])
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func loadMore() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		results := pagination.data

		if len(results) > 20 {
			pagination.data = results[21:]
			results = results[:20]
		} else {
			pagination.data = []string{}
		}

		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		err := enc.Encode(results)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("encoding failure"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(buf.Bytes())
	}
}

func (s *Searcher) Load(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Load: %w", err)
	}
	s.CompleteWorks = string(dat)
	s.SuffixArray = suffixarray.New(dat)
	return nil
}

func (s *Searcher) Search(query string) []string {
	re := regexp.MustCompile(`(?i)` + regexp.QuoteMeta(query))
    idxs := re.FindAllIndex([]byte(s.CompleteWorks), -1)

    results := []string{}
    var lastStart, lastEnd int = -1, -1

    for _, idx := range idxs {
		startIdx := max(0, idx[0]-250)
        endIdx := min(len(s.CompleteWorks), idx[0]+250)

        // Merge overlapping or adjacent segments
        if startIdx <= lastEnd {
            lastEnd = endIdx
        } else {
            if lastStart != -1 {
                results = append(results, s.CompleteWorks[lastStart:lastEnd])
            }
            lastStart, lastEnd = startIdx, endIdx
        }
    }
    
    // Add the last segment if it exists
    if lastStart != -1 {
        results = append(results, s.CompleteWorks[lastStart:lastEnd])
    }

	if len(results) > 20 {
		pagination.data = results[20:]
		return results[:20]
	}

    return results
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}