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
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/geek1011/BookBrowser/booklist"
	"github.com/geek1011/BookBrowser/formats"
	"github.com/geek1011/BookBrowser/indexer"
	"github.com/geek1011/BookBrowser/public"
	"github.com/geek1011/kepubify/kepub"
	"github.com/julienschmidt/httprouter"
	"github.com/unrolled/render"
)

// Server is a BookBrowser server.
type Server struct {
	Indexer  *indexer.Indexer
	BookDir  string
	CoverDir string
	NoCovers bool
	Addr     string
	Verbose  bool
	router   *httprouter.Router
	render   *render.Render
	version  string
}

// NewServer creates a new BookBrowser server. It will not index the books automatically.
func NewServer(addr, bookdir, coverdir, version string, verbose, nocovers bool) *Server {
	i, err := indexer.New([]string{bookdir}, &coverdir, formats.GetExts())
	if err != nil {
		panic(err)
	}
	i.Verbose = verbose

	if verbose {
		log.Printf("Supported formats: %s", strings.Join(formats.GetExts(), ", "))
	}

	s := &Server{
		Indexer:  i,
		BookDir:  bookdir,
		Addr:     addr,
		CoverDir: coverdir,
		NoCovers: nocovers,
		Verbose:  verbose,
		router:   httprouter.New(),
		version:  version,
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
	errs, err := s.Indexer.Refresh()
	if err != nil {
		log.Printf("Error indexing: %s", err)
		return err
	}
	if len(errs) != 0 {
		if s.Verbose {
			log.Printf("Indexing finished with %v errors", len(errs))
		}
	} else {
		log.Printf("Indexing finished")
	}

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
		Directory:  "templates",
		Asset:      public.Box.MustBytes,
		AssetNames: public.Box.List,
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

	s.router.GET("/api/indexer", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"indexing": %t, "progress": %f}`, s.Indexer.Progress != 0, s.Indexer.Progress)
	})

	s.router.GET("/books", s.handleBooks)
	s.router.GET("/books/:id", s.handleBook)

	s.router.GET("/authors", s.handleAuthors)
	s.router.GET("/authors/:id", s.handleAuthor)

	s.router.GET("/series", s.handleSeriess)
	s.router.GET("/series/:id", s.handleSeries)

	s.router.GET("/download", s.handleDownloads)
	s.router.GET("/download/:filename", s.handleDownload)

	s.router.GET("/static/*filepath", func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		http.FileServer(public.Box).ServeHTTP(w, req)
	})
	s.router.ServeFiles("/covers/*filepath", http.Dir(s.CoverDir))
}

func (s *Server) handleDownloads(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
	sbl := s.Indexer.BookList().Sorted(func(a, b *booklist.Book) bool {
		return a.Title < b.Title
	})
	for _, b := range sbl {
		if b.Author != "" && b.Series != "" {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s - %s - %s (%v)</a>", b.ID(), b.FileType(), b.Title, b.Author, b.Series, b.SeriesIndex))
		} else if b.Author != "" && b.Series != "" {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s - %s</a>", b.ID(), b.FileType(), b.Title, b.Author))
		} else if b.Author == "" && b.Series != "" {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s - %s (%v)</a>", b.ID(), b.FileType(), b.Title, b.Series, b.SeriesIndex))
		} else if b.Author == "" && b.Series == "" {
			buf.WriteString(fmt.Sprintf("<a href=\"/download/%s.%s\">%s</a>", b.ID(), b.FileType(), b.Title))
		}
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

	for _, b := range s.Indexer.BookList() {
		if b.ID() == bid {
			if !iskepub {
				rd, err := os.Open(b.FilePath)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					io.WriteString(w, "Error handling request")
					log.Printf("Error handling request for %s: %s\n", r.URL.Path, err)
					return
				}

				w.Header().Set("Content-Disposition", `attachment; filename="`+regexp.MustCompile("[[:^ascii:]]").ReplaceAllString(b.Title, "_")+`.`+b.FileType()+`"`)
				switch b.FileType() {
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
				if b.FileType() != "epub" {
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
				err = (&kepub.Converter{}).Convert(b.FilePath, kepubf)
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

func (s *Server) handleAuthors(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	al := s.Indexer.BookList().Authors().Sorted(func(a, b struct{ Name, ID string }) bool {
		return a.Name < b.Name
	})

	pagination := NewPagination(r.URL.Query(),len(al))
	al = al.Skip(pagination.ItemOffset).Take(pagination.ItemLimit)

	s.render.HTML(w, http.StatusOK, "authors", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Authors",
		"ShowBar":          true,
		"ShowSearch":       false,
		"ShowViewSelector": true,
		"Title":            "Authors",
		"Authors":			al,
		"Pagination":		pagination,
	})
}

func (s *Server) handleAuthor(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	aname := ""
	for _, author := range *s.Indexer.BookList().Authors() {
		if author.ID == p.ByName("id") {
			aname = author.Name
		}
	}

	if aname != "" {
		bl := s.Indexer.BookList().Filtered(func(book *booklist.Book) bool {
			return book.Author != "" && book.AuthorID() == p.ByName("id")
		})
		bl, _ = bl.SortBy("title-asc")
		bl, _ = bl.SortBy(r.URL.Query().Get("sort"))

		pagination := NewPagination(r.URL.Query(),len(bl))
		bl = bl.Skip(pagination.ItemOffset).Take(pagination.ItemLimit)

		s.render.HTML(w, http.StatusOK, "author", map[string]interface{}{
			"CurVersion":       s.version,
			"PageTitle":        aname,
			"ShowBar":          true,
			"ShowSearch":       false,
			"ShowViewSelector": true,
			"Title":            aname,
			"Books":            bl,
			"Pagination":		pagination,
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

func (s *Server) handleSeriess(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	seriess := s.Indexer.BookList().Series().Sorted(func(a, b struct{ Name, ID string }) bool {
		return a.Name < b.Name
	})

	pagination := NewPagination(r.URL.Query(),len(seriess))
	seriess = seriess.Skip(pagination.ItemOffset).Take(pagination.ItemLimit)

	s.render.HTML(w, http.StatusOK, "seriess", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Series",
		"ShowBar":          true,
		"ShowSearch":       false,
		"ShowViewSelector": true,
		"Title":            "Series",
		"Series":			seriess,
		"Pagination":		pagination,
	})
}

func (s *Server) handleSeries(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sname := ""
	for _, series := range *s.Indexer.BookList().Series() {
		if series.ID == p.ByName("id") {
			sname = series.Name
		}
	}

	if sname != "" {
		/* the bl variable created here was unused by the original s.render.HTML() call below and seems to be
			dead code... @geek1011, safe to remove this?

		bl := s.Indexer.BookList().Filtered(func(book *booklist.Book) bool {
			return book.Series != "" && book.SeriesID() == p.ByName("id")
		})
		bl, _ = bl.SortBy("seriesindex-asc")
		bl, _ = bl.SortBy(r.URL.Query().Get("sort"))
		*/


		bl := s.Indexer.BookList().Filtered(func(book *booklist.Book) bool {
			return book.Series != "" && book.SeriesID() == p.ByName("id")
		}).Sorted(func(a, b *booklist.Book) bool {
			return a.SeriesIndex < b.SeriesIndex
		})

		pagination := NewPagination(r.URL.Query(),len(bl))
		bl = bl.Skip(pagination.ItemOffset).Take(pagination.ItemLimit)

		s.render.HTML(w, http.StatusOK, "series", map[string]interface{}{
			"CurVersion":       s.version,
			"PageTitle":        sname,
			"ShowBar":          true,
			"ShowSearch":       false,
			"ShowViewSelector": true,
			"Title":            sname,
			"Books":			bl,
			"Pagination":		pagination,
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

func (s *Server) handleBooks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	bl, _ := s.Indexer.BookList().SortBy("modified-desc")
	bl, _ = bl.SortBy(r.URL.Query().Get("sort"))

	pagination := NewPagination(r.URL.Query(),len(bl))
	bl = bl.Skip(pagination.ItemOffset).Take(pagination.ItemLimit)

	s.render.HTML(w, http.StatusOK, "books", map[string]interface{}{
		"CurVersion":       s.version,
		"PageTitle":        "Books",
		"ShowBar":          true,
		"ShowSearch":       true,
		"ShowViewSelector": true,
		"Title":            "",
		"Books":            bl,
		"Pagination":		pagination,
	})
}

func (s *Server) handleBook(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	for _, b := range s.Indexer.BookList() {
		if b.ID() == p.ByName("id") {
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
	q := r.URL.Query().Get("q")
	ql := strings.ToLower(q)

	if len(q) != 0 {
		bl := s.Indexer.BookList().Filtered(func(a *booklist.Book) bool {
			matches := false
			matches = matches || a.Author != "" && strings.Contains(strings.ToLower(a.Author), ql)
			matches = matches || strings.Contains(strings.ToLower(a.Title), ql)
			matches = matches || a.Series != "" && strings.Contains(strings.ToLower(a.Series), ql)
			return matches
		})
		bl, _ = bl.SortBy("title-asc")
		bl, _ = bl.SortBy(r.URL.Query().Get("sort"))

		pagination := NewPagination(r.URL.Query(),len(bl))
		bl = bl.Skip(pagination.ItemOffset).Take(pagination.ItemLimit)

		s.render.HTML(w, http.StatusOK, "search", map[string]interface{}{
			"CurVersion":       s.version,
			"PageTitle":        "Search Results",
			"ShowBar":          true,
			"ShowSearch":       true,
			"ShowViewSelector": true,
			"Title":            "Search Results",
			"Query":            q,
			"Books":            bl,
			"Pagination":		pagination,
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

func (s *Server) handleRandom(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int() % len(s.Indexer.BookList())
	http.Redirect(w, r, "/books/"+(s.Indexer.BookList())[n].ID(), http.StatusTemporaryRedirect)
}
