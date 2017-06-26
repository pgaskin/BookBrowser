package main

//go:generate go-bindata-assetfs static/...

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DownloadHandler handles file download
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	bid := filepath.Base(r.URL.Path)
	bid = strings.Replace(bid, filepath.Ext(bid), "", 1)

	if bid == "download" {
		w.Header().Set("Content-Type", "text/html")
		for _, b := range books {
			io.WriteString(w, fmt.Sprintf("<a href=\"/download/%s\">%s</a><br>", b.ID, b.Title))
		}
		return
	}

	for _, b := range books {
		if b.ID == bid {
			rd, err := os.Open(b.Filepath)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "Error handling request")
				log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
			}

			w.Header().Set("Content-Disposition", "attachment; filename="+url.PathEscape(b.Title)+".epub")
			w.Header().Set("Content-Type", "application/epub+zip")
			_, err = io.Copy(w, rd)
			rd.Close()
			if err != nil {
				log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, pageHTML("Not Found", "Could not find book with id "+bid))
}

// AuthorsHandler handles the authors page
func AuthorsHandler(w http.ResponseWriter, r *http.Request) {
	aid := filepath.Base(r.URL.Path)

	if aid == "authors" {
		w.Header().Set("Content-Type", "text/html")
		var listHTML bytes.Buffer
		doneAuthors := map[string]bool{}
		for _, b := range books {
			if doneAuthors[b.AuthorID] {
				continue
			}
			doneAuthors[b.AuthorID] = true
			listHTML.WriteString(itemCardHTML(b.Author, "", "/authors/"+b.AuthorID))
		}
		io.WriteString(w, pageHTML("Authors", listHTML.String()))
		return
	}

	found := false
	w.Header().Set("Content-Type", "text/html")
	var booksHTML bytes.Buffer
	booksHTML.WriteString(`<div class="books cards">`)
	aname := ""
	for _, b := range books {
		if b.AuthorID == aid {
			aname = b.Author
			booksHTML.WriteString(bookHTML(&b, true))
			found = true
		}
	}
	booksHTML.WriteString(`</div>`)
	if found != true {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, pageHTML("Not Found", "Could not find author with id "+aid))
	}
	io.WriteString(w, pageHTML(aname, booksHTML.String()))
}

// SeriesHandler handles the series page
func SeriesHandler(w http.ResponseWriter, r *http.Request) {
	sid := filepath.Base(r.URL.Path)

	if sid == "series" {
		w.Header().Set("Content-Type", "text/html")
		var listHTML bytes.Buffer
		doneSeries := map[string]bool{}
		for _, b := range books {
			if b.Series.Name == "" || doneSeries[b.Series.ID] {
				continue
			}
			doneSeries[b.Series.ID] = true
			listHTML.WriteString(itemCardHTML(b.Series.Name, "", "/series/"+b.Series.ID))
		}
		if len(doneSeries) == 0 {
			io.WriteString(w, pageHTML("Series", "No series have been found."))
			return
		}
		io.WriteString(w, pageHTML("Series", listHTML.String()))
		return
	}

	found := false
	w.Header().Set("Content-Type", "text/html")
	var booksHTML bytes.Buffer
	booksHTML.WriteString(`<div class="books cards">`)
	sname := ""
	matched := []Book{}
	for _, b := range books {
		if b.Series.ID == sid {
			sname = b.Series.Name
			matched = append(matched, b)
			found = true
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Series.Index < matched[j].Series.Index
	})
	for _, b := range matched {
		booksHTML.WriteString(bookHTML(&b, true))
	}
	booksHTML.WriteString(`</div>`)
	if found != true {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, pageHTML("Not Found", "Could not find series with id "+sid))
	}
	io.WriteString(w, pageHTML(sname, booksHTML.String()))
}

// BooksHandler handles the books page
func BooksHandler(w http.ResponseWriter, r *http.Request) {
	bid := filepath.Base(r.URL.Path)

	if bid == "books" {
		w.Header().Set("Content-Type", "text/html")
		var booksHTML bytes.Buffer
		booksHTML.WriteString(`<div class="books cards">`)
		for _, b := range books {
			booksHTML.WriteString(bookHTML(&b, true))
		}
		booksHTML.WriteString(`</div>`)
		io.WriteString(w, pageHTML("Books", booksHTML.String()))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	for _, b := range books {
		if b.ID == bid {
			io.WriteString(w, pageHTML(b.Title, bookHTML(&b, false)))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, pageHTML("Not Found", "Could not find book with id "+bid))
}

// SearchHandler handles the search page
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	ql := strings.ToLower(q)

	if len(q) != 0 {
		w.Header().Set("Content-Type", "text/html")
		var booksHTML bytes.Buffer
		booksHTML.WriteString(`<form role="search" method="GET" action="/search/">
<div class="input-group">
<input type="text" class="form-control" placeholder="Search" name="q" id="q" value="` + strings.Replace(q, `"`, "&quot;", -1) + `">
<div class="input-group-btn">
<button class="btn btn-default" type="submit"><i class="glyphicon glyphicon-search"></i></button>
</div>
</div>
</form><br>`)
		booksHTML.WriteString(`<div class="books cards">`)
		matched := false
		for _, b := range books {
			matches := false
			matches = matches || strings.Contains(strings.ToLower(b.Author), ql)
			matches = matches || strings.Contains(strings.ToLower(b.Title), ql)
			matches = matches || strings.Contains(strings.ToLower(b.Series.Name), ql)

			if matches {
				booksHTML.WriteString(bookHTML(&b, true))
				matched = true
			}
		}
		booksHTML.WriteString(`</div>`)
		if !matched {
			booksHTML.WriteString("No books matching your query have been found.")
			return
		}
		io.WriteString(w, pageHTML("Search Results: "+q, booksHTML.String()))
	} else {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, pageHTML("Search", `<form role="search" method="GET" action="/search/">
<div class="input-group">
<input type="text" class="form-control" placeholder="Search" name="q" id="q">
<div class="input-group-btn">
<button class="btn btn-default" type="submit"><i class="glyphicon glyphicon-search"></i></button>
</div>
</div>
</form>`))
	}
}

// JSONHandler handles the books.json file
func JSONHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(books)
	w.Write(b)
}

var books []Book

func runServer(bks []Book, addr string) {
	books = bks

	StaticHandler := http.StripPrefix("/static", http.FileServer(assetFS()))
	CoversHandler := http.StripPrefix("/covers", http.FileServer(http.Dir(*tempdir)))

	http.Handle("/static/", StaticHandler)
	http.Handle("/covers/", CoversHandler)
	http.HandleFunc("/download/", DownloadHandler)
	http.HandleFunc("/authors/", AuthorsHandler)
	http.HandleFunc("/series/", SeriesHandler)
	http.HandleFunc("/books/", BooksHandler)
	http.HandleFunc("/search/", SearchHandler)
	http.HandleFunc("/books.json", JSONHandler)
	http.HandleFunc("/random/", func(w http.ResponseWriter, r *http.Request) {
		rand.Seed(time.Now().Unix())
		n := rand.Int() % len(books)
		http.Redirect(w, r, "/books/"+books[n].ID, http.StatusTemporaryRedirect)
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/books/", http.StatusTemporaryRedirect)
	})

	log.Printf("Serving on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
