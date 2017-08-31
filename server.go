package main

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
	"runtime/debug"
	"sort"
	"strings"
	"time"

	"github.com/geek1011/kepubify/kepub"
	"github.com/julienschmidt/httprouter"
)

//go:generate go-bindata-assetfs static/...

// Server is a BookBrowser server.
type Server struct {
	Books    *BookList
	BookDir  string
	CoverDir string
	Addr     string
	Verbose  bool
	router   *httprouter.Router
}

// NewServer creates a new BookBrowser server. It will not index the books automatically.
func NewServer(addr, bookdir, coverdir string, verbose bool) *Server {
	s := &Server{
		Books:    &BookList{},
		BookDir:  bookdir,
		Addr:     addr,
		CoverDir: coverdir,
		Verbose:  verbose,
		router:   httprouter.New(),
	}

	s.initRouter()

	return s
}

// printLog runs log.Printf if verbose is true.
func (s *Server) printLog(format string, v ...interface{}) {
	if s.Verbose {
		log.Printf(format, v...)
	}
}

// RefreshBookIndex refreshes the book index
func (s *Server) RefreshBookIndex() error {
	books, err := NewBookListFromDir(s.BookDir, s.CoverDir, s.Verbose)
	if err != nil {
		debug.FreeOSMemory()
		return err
	}

	s.Books = books
	debug.FreeOSMemory()
	return nil
}

// Serve starts the BookBrowser server. It does not return unless there is an error.
func (s *Server) Serve() error {
	s.printLog("Serving on %s\n", s.Addr)
	err := http.ListenAndServe(s.Addr, s.router)
	if err != nil {
		return err
	}
	return nil
}

// initRouter initializes the router for the BookBrowser server.
func (s *Server) initRouter() {
	s.router = httprouter.New()

	s.router.GET("/", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(w, r, "/books/", http.StatusTemporaryRedirect)
	})

	s.router.GET("/random", s.handleRandom)

	s.router.GET("/search", s.handleSearch)

	s.router.GET("/books.json", s.handleBooksJSON)

	s.router.GET("/books", s.handleBookList)
	s.router.GET("/books/:id", s.handleBook)

	s.router.GET("/authors", s.handleAuthorList)
	s.router.GET("/authors/:id", s.handleAuthor)

	s.router.GET("/series", s.handleSeriesList)
	s.router.GET("/series/:id", s.handleSeries)

	s.router.GET("/download", s.handleDownloadList)
	s.router.GET("/download/:filename", s.handleDownload)

	s.router.ServeFiles("/static/*filepath", assetFS())
	s.router.ServeFiles("/covers/*filepath", http.Dir(s.CoverDir))
}

type nameID struct {
	Name string
	ID   string
}

func (s *Server) sortedBookPropertyList(getNameID func(Book) nameID, filterNameID func(nameID) bool, sortNameID func(nameID, nameID) bool) []nameID {
	doneItems := map[string]bool{}
	items := []nameID{}
	for _, b := range *s.Books {
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

func (s *Server) sortedBookList(filterBook func(Book) bool, sortBook func(Book, Book) bool) []Book {
	filteredItems := []Book{}
	for _, book := range *s.Books {
		if filterBook(book) {
			filteredItems = append(filteredItems, book)
		}
	}
	sort.Slice(filteredItems, func(i, j int) bool {
		return sortBook(filteredItems[i], filteredItems[j])
	})
	return filteredItems
}

func (s *Server) handleDownloadList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	sbl := s.sortedBookList(func(b Book) bool {
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
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	bid := p.ByName("filename")
	bid = strings.Replace(strings.Replace(bid, filepath.Ext(bid), "", 1), ".kepub", "", -1)
	iskepub := false
	if strings.HasSuffix(p.ByName("filename"), ".kepub.epub") {
		iskepub = true
	}

	for _, b := range *s.Books {
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

func (s *Server) handleAuthorList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	var listHTML bytes.Buffer

	authors := s.sortedBookPropertyList(func(b Book) nameID {
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
}

func (s *Server) handleAuthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	aid := p.ByName("id")

	w.Header().Set("Content-Type", "text/html")

	matched := s.sortedBookList(func(book Book) bool {
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

func (s *Server) handleSeriesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")
	var listHTML bytes.Buffer

	series := s.sortedBookPropertyList(func(b Book) nameID {
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
}

func (s *Server) handleSeries(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sid := p.ByName("id")

	w.Header().Set("Content-Type", "text/html")

	matched := s.sortedBookList(func(book Book) bool {
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

func (s *Server) handleBookList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "text/html")

	matched := s.sortedBookList(func(book Book) bool {
		return true
	}, func(a Book, b Book) bool {
		return a.ModTime.Unix() > b.ModTime.Unix()
	})

	html, notfound := bookListPageHTML(matched, "Books", "There are no books in your library.", true)

	if notfound {
		w.WriteHeader(http.StatusNotFound)
	}

	io.WriteString(w, html)
}

func (s *Server) handleBook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	bid := p.ByName("id")

	w.Header().Set("Content-Type", "text/html")
	for _, b := range *s.Books {
		if b.ID == bid {
			io.WriteString(w, pageHTML(b.Title, bookHTML(&b, true), false, false))
			return
		}
	}
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	q := r.URL.Query().Get("q")
	ql := strings.ToLower(q)

	if len(q) != 0 {
		w.Header().Set("Content-Type", "text/html")
		var booksHTML bytes.Buffer
		booksHTML.WriteString(`<script>document.querySelector(".q").value="` + strings.Replace(q, `"`, `\"`, -1) + `";</script>`)
		booksHTML.WriteString(`<div class="books view cards">`)
		matched := false
		for _, b := range *s.Books {
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

func (s *Server) handleBooksJSON(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	b, _ := json.Marshal(s.Books)
	w.Write(b)
}

func (s *Server) handleRandom(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rand.Seed(time.Now().Unix())
	n := rand.Int() % len(*s.Books)
	http.Redirect(w, r, "/books/"+(*s.Books)[n].ID, http.StatusTemporaryRedirect)
}
