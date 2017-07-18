package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

var bookdir *string
var tempdir *string
var addr *string

var curversion = "undefined"

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err)
	}

	td, err := ioutil.TempDir("", "bookbrowser")
	if err != nil {
		td = filepath.Join(wd, "_temp")
	}

	bookdir = flag.String("bookdir", wd, "The directory to get books from. This directory must exist.")
	tempdir = flag.String("tempdir", td, "The directory to use for storing temporary files such as book cover thumbnails. This directory is create on start and deleted on exit.")
	addr = flag.String("addr", ":8090", "The address to bind to.")
	flag.Parse()

	log.Printf("BookBrowser %s\n", curversion)

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

	books, err := NewBookListFromDir(*bookdir, *tempdir, true)
	if err != nil {
		log.Fatalf("Fatal error indexing books: %s\n", err)
	}

	if len(*books) == 0 {
		log.Fatalln("Fatal error: no books found")
	}

	runServer(*books, *addr)
}
