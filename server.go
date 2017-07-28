package main

//go:generate go-bindata-assetfs static/...

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/geek1011/kepubify/kepub"
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
	bid = strings.Replace(strings.Replace(bid, filepath.Ext(bid), "", 1), ".kepub", "", -1)
	iskepub := false
	if strings.HasSuffix(r.URL.Path, ".kepub.epub") {
		iskepub = true
	}

	if bid == "download" {
		w.Header().Set("Content-Type", "text/html")
		var buf bytes.Buffer
		buf.WriteString(`
<!DOCTYPE html>
<html>
<head>
<title>BookBrowser</title>
<style>
a,
a:link,
a:visited {
    display:  block;
    white-space: nowrap;
    text-overflow: ellipsis;
    color: inherit;
    text-decoration: none;
    font-family: sans-serif;
    padding: 5px 7px;
    background:  #FAFAFA;
    border-bottom: 1px solid #DDDDDD;
    cursor: pointer;
}

a:hover,
a:active {
    background: #EEEEEE;
}

html, body {
    background: #FAFAFA;
    margin: 0;
    padding: 0;
}
</style>
</head>
<body>
		`)
		sbl := sortedBookList(books, func(b Book) bool {
			return true
		}, func(a Book, b Book) bool {
			return a.Title < b.Title
		})
		for _, b := range sbl {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s - %s - %s (%v)</a>", b.ID, b.FileType, b.Title, b.Author, b.Series.Name, b.Series.Index))
		}
		buf.WriteString(`
</body>
</html>
		`)
		io.WriteString(w, buf.String())
		return
	}

	for _, b := range books {
		if b.ID == bid {
			if !iskepub {
				rd, err := os.Open(b.Filepath)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					io.WriteString(w, "Error handling request")
					log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
					return
				}

				w.Header().Set("Content-Disposition", "attachment; filename="+url.PathEscape(b.Title)+"."+b.FileType)
				switch b.FileType {
				case "epub":
					w.Header().Set("Content-Type", "application/epub+zip")
				case "pdf":
					w.Header().Set("Content-Type", "application/pdf")
				default:
					w.Header().Set("Content-Type", "application/octet-stream")
				}
				_, err = io.Copy(w, rd)
				rd.Close()
				if err != nil {
					log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
				}
			} else {
				if b.FileType != "epub" {
					w.WriteHeader(http.StatusNotFound)
					io.WriteString(w, "Not found")
					return
				}
				td, err := ioutil.TempDir("", "kepubify")
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
					io.WriteString(w, "Internal Server Error")
					return
				}
				defer os.RemoveAll(td)
				kepubf := filepath.Join(td, bid+".kepub.epub")
				err = kepub.Kepubify(b.Filepath, kepubf, false)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
					io.WriteString(w, "Internal Server Error - Error converting book")
					return
				}
				rd, err := os.Open(kepubf)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					io.WriteString(w, "Error handling request")
					log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
					return
				}
				w.Header().Set("Content-Disposition", "attachment; filename="+url.PathEscape(b.Title)+".kepub.epub")
				w.Header().Set("Content-Type", "application/epub+zip")
				_, err = io.Copy(w, rd)
				rd.Close()
				if err != nil {
					log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
				}
			}
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, pageHTML("Not Found", "Could not find book with id "+bid, false, false))
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
		listHTML.WriteString(`<div class="items view cards">`)
		for _, ni := range authors {
			listHTML.WriteString(itemCardHTML(ni.Name, "/authors/"+ni.ID))
		}
		listHTML.WriteString(`</div>`)

		io.WriteString(w, pageHTML("Authors", listHTML.String(), true, false))
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

	html, notfound := bookListPageHTML(matched, aname, "Author not found", false)

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
		listHTML.WriteString(`<div class="items view cards">`)
		for _, ni := range series {
			listHTML.WriteString(itemCardHTML(ni.Name, "/series/"+ni.ID))
		}
		listHTML.WriteString(`</div>`)
		if len(series) == 0 {
			io.WriteString(w, pageHTML("Series", "No series have been found.", false, false))
			return
		}

		io.WriteString(w, pageHTML("Series", listHTML.String(), true, false))
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

	html, notfound := bookListPageHTML(matched, sname, "Series not found", false)

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

		html, notfound := bookListPageHTML(matched, "Books", "There are no books in your library.", true)

		if notfound {
			w.WriteHeader(http.StatusNotFound)
		}

		io.WriteString(w, html)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	for _, b := range books {
		if b.ID == bid {
			io.WriteString(w, pageHTML(b.Title, bookHTML(&b, true), false, false))
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, pageHTML("Not Found", "Could not find book with id "+bid, false, false))
}

// SearchHandler handles the search page
func SearchHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	ql := strings.ToLower(q)

	if len(q) != 0 {
		w.Header().Set("Content-Type", "text/html")
		var booksHTML bytes.Buffer
		booksHTML.WriteString(`<script>document.querySelector(".q").value="` + strings.Replace(q, `"`, `\"`, -1) + `";</script>`)
		booksHTML.WriteString(`<div class="books view cards">`)
		matched := false
		for _, b := range books {
			matches := false
			matches = matches || strings.Contains(strings.ToLower(b.Author), ql)
			matches = matches || strings.Contains(strings.ToLower(b.Title), ql)
			matches = matches || strings.Contains(strings.ToLower(b.Series.Name), ql)

			if matches {
				booksHTML.WriteString(bookHTML(&b, false))
				matched = true
			}
		}
		booksHTML.WriteString(`</div>`)
		if !matched {
			booksHTML.WriteString("No books matching your query have been found.")
		}
		io.WriteString(w, pageHTML("Search Results", booksHTML.String(), true, true))
	} else {
		w.Header().Set("Content-Type", "text/html")
		io.WriteString(w, pageHTML("Search", `<center><a href="/static/list.html">Advanced Search</a></center>`, false, true))
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
