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

type nameID struct {
	Name string
	ID   string
}

func sortedBookPropertyList(books []Book, getNameID func(Book) nameID, filterNameID func(nameID) bool, sortNameID func(nameID, nameID) bool) []nameID {
	doneItems := map[string]bool{}
	items := []nameID{}
	for _, b := range books {
		nid := getNameID(b)
		if doneItems[nid.ID] {
			continue
		}
		doneItems[nid.ID] = true
		items = append(items, nameID{
			Name: nid.Name,
			ID:   nid.ID,
		})
	}
	filteredItems := []nameID{}
	for _, ni := range items {
		if filterNameID(ni) {
			filteredItems = append(filteredItems, ni)
		}
	}
	sort.Slice(filteredItems, func(i, j int) bool {
		return sortNameID(filteredItems[i], filteredItems[j])
	})
	return filteredItems
}

func sortedBookList(books []Book, filterBook func(Book) bool, sortBook func(Book, Book) bool) []Book {
	filteredItems := []Book{}
	for _, book := range books {
		if filterBook(book) {
			filteredItems = append(filteredItems, book)
		}
	}
	sort.Slice(filteredItems, func(i, j int) bool {
		return sortBook(filteredItems[i], filteredItems[j])
	})
	return filteredItems
}

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

		authors := sortedBookPropertyList(books, func(b Book) nameID {
			return nameID{
				Name: b.Author,
				ID:   b.AuthorID,
			}
		}, func(ni nameID) bool {
			return ni.Name != ""
		}, func(a nameID, b nameID) bool {
			return a.Name < b.Name
		})
		for _, ni := range authors {
			listHTML.WriteString(itemCardHTML(ni.Name, "", "/authors/"+ni.ID))
		}

		io.WriteString(w, pageHTML("Authors", listHTML.String()))
		return
	}

	w.Header().Set("Content-Type", "text/html")

	matched := sortedBookList(books, func(book Book) bool {
		return book.AuthorID == aid
	}, func(a Book, b Book) bool {
		return a.Title < b.Title
	})

	aname := ""
	if len(matched) != 0 {
		aname = matched[0].Author
	}

	html, notfound := bookListPageHTML(matched, aname, "Author not found")

	if notfound {
		w.WriteHeader(http.StatusNotFound)
	}

	io.WriteString(w, html)
}

// SeriesHandler handles the series page
func SeriesHandler(w http.ResponseWriter, r *http.Request) {
	sid := filepath.Base(r.URL.Path)

	if sid == "series" {
		w.Header().Set("Content-Type", "text/html")
		var listHTML bytes.Buffer

		series := sortedBookPropertyList(books, func(b Book) nameID {
			return nameID{
				Name: b.Series.Name,
				ID:   b.Series.ID,
			}
		}, func(ni nameID) bool {
			return ni.Name != ""
		}, func(a nameID, b nameID) bool {
			return a.Name < b.Name
		})
		for _, ni := range series {
			listHTML.WriteString(itemCardHTML(ni.Name, "", "/series/"+ni.ID))
		}
		if len(series) == 0 {
			io.WriteString(w, pageHTML("Series", "No series have been found."))
			return
		}

		io.WriteString(w, pageHTML("Series", listHTML.String()))
		return
	}

	w.Header().Set("Content-Type", "text/html")

	matched := sortedBookList(books, func(book Book) bool {
		return book.Series.ID == sid
	}, func(a Book, b Book) bool {
		return a.Series.Index < b.Series.Index
	})

	sname := ""
	if len(matched) != 0 {
		sname = matched[0].Series.Name
	}

	html, notfound := bookListPageHTML(matched, sname, "Series not found")

	if notfound {
		w.WriteHeader(http.StatusNotFound)
	}

	io.WriteString(w, html)
}

// BooksHandler handles the books page
func BooksHandler(w http.ResponseWriter, r *http.Request) {
	bid := filepath.Base(r.URL.Path)

	if bid == "books" {
		w.Header().Set("Content-Type", "text/html")

		matched := sortedBookList(books, func(book Book) bool {
			return true
		}, func(a Book, b Book) bool {
			return a.ModTime.Unix() > b.ModTime.Unix()
		})

		html, notfound := bookListPageHTML(matched, "Books", "There are no books in your library.")

		if notfound {
			w.WriteHeader(http.StatusNotFound)
		}

		io.WriteString(w, html)
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
