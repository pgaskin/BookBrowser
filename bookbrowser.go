package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

// Series represents a book series
type Series struct {
	Name  string  `json:"name,omitempty"`
	ID    string  `json:"id,omitempty"`
	Index float64 `json:"index,omitempty"`
}

// Book represents a book
type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author,omitempty"`
	AuthorID    string `json:"authorid"`
	Publisher   string `json:"publisher,omitempty"`
	Description string `json:"description,omitempty"`
	Series      Series `json:"series,omitempty"`
	Filepath    string `json:"filepath"`
	HasCover    bool   `json:"hascover"`
}

var bookdir *string
var tempdir *string
var addr *string

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err)
	}

	bookdir = flag.String("bookdir", wd, "The directory to get books from. This directory must exist.")
	tempdir = flag.String("tempdir", filepath.Join(wd, "_temp"), "The directory to use for storing temporary files such as book cover thumbnails. This directory is create on start and deleted on exit.")
	addr = flag.String("addr", ":8090", "The address to bind to.")
	flag.Parse()

	if _, err := os.Stat(*bookdir); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Fatal error: book directory %s does not exist\n", *bookdir)
		}
	}

	*bookdir, err = filepath.Abs(*bookdir)
	if err != nil {
		log.Fatalf("Fatal error: Could not resolve book directory %s: %s\n", *bookdir, err)
	}

	if _, err := os.Stat(*tempdir); os.IsNotExist(err) {
		os.Mkdir(*tempdir, os.ModePerm)
	}

	*tempdir, err = filepath.Abs(*tempdir)
	if err != nil {
		log.Fatalf("Fatal error: Could not resolve temp directory %s: %s\n", *tempdir, err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		log.Println("Cleaning up covers")
		os.RemoveAll(*tempdir)
		os.Exit(0)
	}()

	books, err := indexBooks()
	if err != nil {
		log.Fatalf("Error indexing books: %s\n", err)
	}

	runServer(books, *addr)
}
