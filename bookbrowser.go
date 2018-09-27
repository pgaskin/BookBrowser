package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	_ "github.com/geek1011/BookBrowser/formats/epub"
	_ "github.com/geek1011/BookBrowser/formats/mobi"
	_ "github.com/geek1011/BookBrowser/formats/pdf"
	"github.com/geek1011/BookBrowser/server"
	"github.com/geek1011/BookBrowser/util"
	"github.com/geek1011/BookBrowser/util/sigusr"
	"github.com/spf13/pflag"
)

var curversion = "dev"

func main() {
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err)
	}

	deftempdir, err := ioutil.TempDir("", "bookbrowser")
	if err != nil {
		deftempdir = filepath.Join(workdir, "_temp")
	}

	bookdir := pflag.StringP("bookdir", "b", workdir, "the directory to load books from (must exist)")
	tempdir := pflag.StringP("tempdir", "t", deftempdir, "the directory to store temp files such as cover thumbnails (created on start, deleted on exit unless already exists)")
	addr := pflag.StringP("addr", "a", ":8090", "the address to bind the server to ([IP]:PORT)")
	nocovers := pflag.BoolP("nocovers", "n", false, "do not index covers")
	help := pflag.BoolP("help", "h", false, "Show this help text")
	sversion := pflag.Bool("version", false, "Show the version")
	pflag.Parse()

	if *sversion {
		fmt.Printf("BookBrowser %s\n", curversion)
		os.Exit(0)
	}

	if *help || pflag.NArg() != 0 {
		fmt.Fprintf(os.Stderr, "Usage: BookBrowser [OPTIONS]\n\nVersion:\n  BookBrowser %s\n\nOptions:\n", curversion)
		pflag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
		if runtime.GOOS == "windows" {
			time.Sleep(time.Second * 2)
		}
		os.Exit(1)
	}

	noRemoveTempDir := false

	log.Printf("BookBrowser %s\n", curversion)

	if _, err := os.Stat(*bookdir); err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Error: book directory %s does not exist\n", *bookdir)
		}
	}

	if fi, err := os.Stat(*tempdir); err == nil || (fi != nil && fi.IsDir()) {
		noRemoveTempDir = true
		if *tempdir == deftempdir {
			noRemoveTempDir = false
		}
	}

	*bookdir, err = filepath.Abs(*bookdir)
	if err != nil {
		log.Fatalf("Error: could not resolve book directory %s: %v\n", *bookdir, err)
	}

	if _, err := os.Stat(*tempdir); os.IsNotExist(err) {
		os.Mkdir(*tempdir, os.ModePerm)
	}

	*tempdir, err = filepath.Abs(*tempdir)
	if err != nil {
		log.Fatalf("Error: could not resolve temp directory %s: %v\n", *tempdir, err)
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		if noRemoveTempDir {
			log.Println("Not removing temp dir because dir already existed at start")
		} else {
			log.Println("Cleaning up temp dir")
			os.RemoveAll(*tempdir)
		}
		os.Exit(0)
	}()

	if !strings.Contains(*addr, ":") {
		log.Fatalln("Error: invalid listening address")
	}

	sp := strings.SplitN(*addr, ":", 2)
	if sp[0] == "" {
		ip := util.GetIP()
		if ip != nil {
			log.Printf("This server can be accessed at http://%s:%s\n", ip.String(), sp[1])
		}
	}

	s := server.NewServer(*addr, *bookdir, *tempdir, curversion, true, *nocovers)
	go func() {
		s.RefreshBookIndex()
		if len(s.Indexer.BookList()) == 0 {
			log.Fatalln("Fatal error: no books found")
		}
		checkUpdate()
	}()

	sigusr.Handle(func() {
		log.Println("Booklist refresh triggered by SIGUSR1")
		s.RefreshBookIndex()
	})

	err = s.Serve()
	if err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}

func checkUpdate() {
	resp, err := http.Get("https://api.github.com/repos/geek1011/BookBrowser/releases/latest")
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		return
	}

	var obj struct {
		URL string `json:"html_url"`
		Tag string `json:"tag_name"`
	}
	if json.Unmarshal(buf, &obj) != nil {
		return
	}

	if curversion != "dev" {
		if !strings.HasPrefix(curversion, obj.Tag) {
			log.Printf("Running version %s. Latest version is %s: %s\n", curversion, obj.Tag, obj.URL)
		}
	}
}
