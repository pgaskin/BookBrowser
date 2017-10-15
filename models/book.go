package models

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"time"

	"github.com/geek1011/BookBrowser/modules/util"
)

// NameID represents a name and an id
type NameID interface {
	GetName() string
	GetID() string
}

// Series represents a book series
type Series struct {
	Name  string  `json:"name,omitempty"`
	ID    string  `json:"id,omitempty"`
	Index float64 `json:"index,omitempty"`
}

// GetName gets the name of the Series
func (s *Series) GetName() string {
	return s.Name
}

// GetID gets the id of the Series
func (s *Series) GetID() string {
	return s.ID
}

// Author represents a book author
type Author struct {
	Name string `json:"name,omitempty"`
	ID   string `json:"id,omitempty"`
}

// GetName gets the name of the Author
func (a *Author) GetName() string {
	return a.Name
}

// GetID gets the id of the Author
func (a *Author) GetID() string {
	return a.ID
}

// Book is a book.
type Book struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      *Author   `json:"author,omitempty"`
	Publisher   string    `json:"publisher,omitempty"`
	Description string    `json:"description,omitempty"`
	Series      *Series   `json:"series,omitempty"`
	Filepath    string    `json:"filepath"`
	HasCover    bool      `json:"hascover"`
	ModTime     time.Time `json:"modtime,omitempty"`
	FileType    string    `json:"filetype,omitempty"` // Does not include leading period
}

// NewBook creates a new book
func NewBook(title, author, publisher, seriesName string, seriesIndex float64, description, filepath string, hascover bool, modtime time.Time, filetype string) *Book {
	book := &Book{
		Title: util.FixString(title),
		Author: &Author{
			Name: util.FixString(author),
		},
		Publisher:   util.FixString(publisher),
		Description: util.FixString(description),
		Series: &Series{
			Name:  util.FixString(seriesName),
			Index: seriesIndex,
		},
		Filepath: filepath,
		HasCover: hascover,
		ModTime:  modtime,
		FileType: filetype,
	}

	id := sha1.New()
	io.WriteString(id, book.Author.Name) // If empty, then it hashes an empty string to retain compatibility with old BookBrowser versions
	book.Author.ID = hex.EncodeToString(id.Sum(nil))[:10]
	io.WriteString(id, book.Series.Name) // If empty, then it hashes an empty string to retain compatibility with old BookBrowser versions
	io.WriteString(id, book.Title)
	book.ID = hex.EncodeToString(id.Sum(nil))[:10]

	id = sha1.New()
	io.WriteString(id, book.Series.Name)
	book.Series.ID = hex.EncodeToString(id.Sum(nil))[:10]

	if seriesName == "" {
		book.Series = nil
	}

	if author == "" {
		book.Author = nil
	}

	return book
}
