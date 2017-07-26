package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestEPUBMetadata(t *testing.T) {
	td, err := ioutil.TempDir("", "bookbrowser")
	if err != nil {
		t.Fatalf("Cannot create temp dir: %s", err.Error())
	}
	defer os.RemoveAll(td)

	book, err := NewBookFromFile("testdata/books/test1.epub", td)
	if err != nil {
		t.Fatalf("Could not load book: %s", err.Error())
	}

	coverfile := filepath.Join(td, book.ID+".jpg")
	if _, err := os.Stat(coverfile); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Cover file %s does not exist\n", coverfile)
		}
	}

	coverthumbfile := filepath.Join(td, book.ID+"_thumb.jpg")
	if _, err := os.Stat(coverthumbfile); err != nil {
		if os.IsNotExist(err) {
			t.Errorf("Cover thumbnail file %s does not exist\n", coverthumbfile)
		}
	}

	if book.Title != "BookBrowser Test Book 1" {
		t.Errorf("Incorrect title: %s", book.Title)
	}

	if book.Author != "Patrick G" {
		t.Errorf("Incorrect author: %s", book.Author)
	}

	if book.Publisher != "Patrick G" {
		t.Errorf("Incorrect publisher: %s", book.Publisher)
	}

	if book.Description != "<p>This is a test book for <i>BookBrowser</i>, a ebook content server.</p>" {
		t.Errorf("Incorrect description: %s", book.Description)
	}

	if book.FileType != "epub" {
		t.Errorf("Incorrect filetype: %s", book.FileType)
	}

	if book.HasCover != true {
		t.Errorf("Incorrect hascover value: %v", book.HasCover)
	}

	if book.Series.Name != "Test Series" {
		t.Errorf("Incorrect series name: %s", book.Series.Name)
	}

	if book.Series.Index != 1 {
		t.Errorf("Incorrect series index: %v", book.Series.Index)
	}

	if book.ID != "a611744562" {
		t.Errorf("Incorrect book id: %s", book.ID)
	}
}
