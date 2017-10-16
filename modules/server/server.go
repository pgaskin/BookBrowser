package server

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/geek1011/BookBrowser/models"
	"github.com/geek1011/BookBrowser/modules/booklist"
	"github.com/geek1011/kepubify/kepub"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

//go:generate go-bindata-assetfs -pkg server -prefix ../../  ../../public/...

// Server is a BookBrowser server.
type Server struct {
	Books     *booklist.BookList
	booksLock *sync.RWMutex
	BookDir   string
	CoverDir  string
	Addr      string
	Verbose   bool
	router    *httprouter.Router
	render    *render.Render
	version   string
}

// NewServer creates a new BookBrowser server. It will not index the books automatically.
func NewServer(addr, bookdir, coverdir, version string, verbose bool) *Server {
	s := &Server{
		Books:     &booklist.BookList{},
		booksLock: &sync.RWMutex{},
		BookDir:   bookdir,
		Addr:      addr,
		CoverDir:  coverdir,
		Verbose:   verbose,
		router:    httprouter.New(),
		version:   version,
	}

	s.initRender()
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
	s.printLog("Locking book index\n")
	s.booksLock.Lock()
	defer s.printLog("Unlocking book index\n")
	defer s.booksLock.Unlock()

	books, errs := booklist.NewBookListFromDir(s.BookDir, s.CoverDir, s.Verbose)
	if len(errs) != 0 {
		if s.Verbose {
			log.Printf("Indexing finished with %v errors", len(errs))
		}
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

// initRender initializes the renderer for the BookBrowser server.
func (s *Server) initRender() {
	s.render = render.New(render.Options{
		Directory:  "public/templates",
		Asset:      Asset,
		AssetNames: AssetNames,
		Layout:     "base",
		Extensions: []string{".tmpl"},
		Funcs: []template.FuncMap{
			template.FuncMap{
				"ToUpper": strings.ToUpper,
				"raw": func(s string) template.HTML {
					return template.HTML(s)
				},
			},
		},
		IsDevelopment: false,
	})
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

	s.router.GET("/static/*filepath", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		req.URL.Path = "/static/" + ps.ByName("filepath")
		http.FileServer(assetFS()).ServeHTTP(w, req)
	})
	s.router.ServeFiles("/covers/*filepath", http.Dir(s.CoverDir))
}

func (s *Server) handleDownloadList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

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
	sbl := s.Books.Sorted(func(a, b *models.Book) bool {
		return a.Title < b.Title
	})
	for _, b := range sbl {
		if b.Author != nil && b.Series != nil {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s - %s - %s (%v)</a>", b.ID, b.FileType, b.Title, b.Author.Name, b.Series.Name, b.Series.Index))
		} else if b.Author != nil && b.Series == nil {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s - %s</a>", b.ID, b.FileType, b.Title, b.Author.Name))
		} else if b.Author == nil && b.Series != nil {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s - %s (%v)</a>", b.ID, b.FileType, b.Title, b.Series.Name, b.Series.Index))
		} else if b.Author == nil && b.Series == nil {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s</a>", b.ID, b.FileType, b.Title))
		}
	}
	buf.WriteString(`
</body>
</html>
	`)
	io.WriteString(w, buf.String())
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

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
	io.WriteString(w, "Could not find book with id "+bid)
}

func (s *Server) handleAuthorList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	s.render.HTML(w, http.StatusOK, "authors", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Authors",
		"ShowBar":          true,
		"ShowSearch":       false,
		"ShowViewSelector": true,
		"Title":            "Authors",
		"Authors": s.Books.GetAuthors().Sorted(func(a, b *models.Author) bool {
			return a.Name < b.Name
		}),
	})
}

func (s *Server) handleAuthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	aname := ""
	for _, author := range *s.Books.GetAuthors() {
		if author.ID == p.ByName("id") {
			aname = author.Name
		}
	}

	if aname != "" {
		s.render.HTML(w, http.StatusOK, "author", map[string]interface{}{
			"CurVersion":       s.version,
			"PageTitle":        aname,
			"ShowBar":          true,
			"ShowSearch":       false,
			"ShowViewSelector": true,
			"Title":            aname,
			"Books": s.Books.Filtered(func(book *models.Book) bool {
				return book.Author != nil && book.Author.ID == p.ByName("id")
			}).Sorted(func(a, b *models.Book) bool {
				return a.Title < b.Title
			}),
		})
		return
	}

	s.render.HTML(w, http.StatusNotFound, "notfound", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Not Found",
		"ShowBar":          false,
		"ShowSearch":       false,
		"ShowViewSelector": false,
		"Title":            "Not Found",
		"Message":          "Author not found.",
	})
}

func (s *Server) handleSeriesList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	s.render.HTML(w, http.StatusOK, "seriess", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Series",
		"ShowBar":          true,
		"ShowSearch":       false,
		"ShowViewSelector": true,
		"Title":            "Series",
		"Series": s.Books.GetSeries().Sorted(func(a, b *models.Series) bool {
			return a.Name < b.Name
		}),
	})
}

func (s *Server) handleSeries(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	sname := ""
	for _, series := range *s.Books.GetSeries() {
		if series.ID == p.ByName("id") {
			sname = series.Name
		}
	}

	if sname != "" {
		s.render.HTML(w, http.StatusOK, "series", map[string]interface{}{
			"CurVersion":       s.version,
			"PageTitle":        sname,
			"ShowBar":          true,
			"ShowSearch":       false,
			"ShowViewSelector": true,
			"Title":            sname,
			"Books": s.Books.Filtered(func(book *models.Book) bool {
				return book.Series != nil && book.Series.ID == p.ByName("id")
			}).Sorted(func(a, b *models.Book) bool {
				return a.Series.Index < b.Series.Index
			}),
		})
		return
	}

	s.render.HTML(w, http.StatusNotFound, "notfound", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Not Found",
		"ShowBar":          false,
		"ShowSearch":       false,
		"ShowViewSelector": false,
		"Title":            "Not Found",
		"Message":          "Series not found.",
	})
}

func (s *Server) handleBookList(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	s.render.HTML(w, http.StatusOK, "books", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Books",
		"ShowBar":          true,
		"ShowSearch":       true,
		"ShowViewSelector": true,
		"Title":            "",
		"Books":            s.Books,
	})
}

func (s *Server) handleBook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	for _, b := range *s.Books {
		if b.ID == p.ByName("id") {
			s.render.HTML(w, http.StatusOK, "book", map[string]interface{}{
				"CurVersion":       s.version,
				"PageTitle":        b.Title,
				"ShowBar":          false,
				"ShowSearch":       false,
				"ShowViewSelector": false,
				"Title":            "",
				"Book":             b,
			})
			return
		}
	}

	s.render.HTML(w, http.StatusNotFound, "notfound", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Not Found",
		"ShowBar":          false,
		"ShowSearch":       false,
		"ShowViewSelector": false,
		"Title":            "Not Found",
		"Message":          "Book not found.",
	})
}

func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	q := r.URL.Query().Get("q")
	ql := strings.ToLower(q)

	if len(q) != 0 {
		s.render.HTML(w, http.StatusOK, "search", map[string]interface{}{
			"CurVersion":       s.version,
			"PageTitle":        "Search Results",
			"ShowBar":          true,
			"ShowSearch":       true,
			"ShowViewSelector": true,
			"Title":            "Search Results",
			"Query":            q,
			"Books": s.Books.Filtered(func(a *models.Book) bool {
				matches := false
				matches = matches || a.Author != nil && strings.Contains(strings.ToLower(a.Author.Name), ql)
				matches = matches || strings.Contains(strings.ToLower(a.Title), ql)
				matches = matches || a.Series != nil && strings.Contains(strings.ToLower(a.Series.Name), ql)
				return matches
			}),
		})
		return
	}

	s.render.HTML(w, http.StatusOK, "search", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Search",
		"ShowBar":          true,
		"ShowSearch":       true,
		"ShowViewSelector": false,
		"Title":            "Search",
		"Query":            "",
	})
}

func (s *Server) handleBooksJSON(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	s.render.JSON(w, http.StatusOK, s.Books)
}

func (s *Server) handleRandom(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s.booksLock.RLock()
	defer s.booksLock.RUnlock()

	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % len(*s.Books)
	http.Redirect(w, r, "/books/"+(*s.Books)[n].ID, http.StatusTemporaryRedirect)
}
