package main

import (
	"bytes"
	"fmt"
	"strings"
)

func pageHTML(title string, content string, containsview bool, showsearch bool) string {
	var html bytes.Buffer
	html.WriteString(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta http-equiv="X-UA-Compatible" content="ie=edge">
<title>`)
	html.WriteString(title)
	html.WriteString(`</title>
<link rel="stylesheet" href="/static/style.css">
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black">
</head>
<body class="light `)
	if containsview {
		html.WriteString("contains-view")
	} else {
		html.WriteString("no-contains-view")
	}
	html.WriteString(`">
<div class="container">
<div class="nav">
<div class="left">
<div class="title">
<a href="/books/" class="item">BookBrowser</a>
</div>
</div>
<div class="right">
<a href="/books/" class="item">Books</a>
<a href="/authors/" class="item">Authors</a>
<a href="/series/" class="item">Series</a>
<a href="/static/list.html" class="item">List</a>
<a href="/random/" class="item">Random</a>
<a href="/search/" class="item">Search</a>
</div>
</div>
<div class="section">
<div class="title">`)
	html.WriteString(title)
	html.WriteString(`</div>
<div class="body">
<div class="bar">`)
	if showsearch {
		html.WriteString(`<div class="search">
<form action="/search/" method="GET">
<input class="q" name="q" type="search" placeholder="Search books..." />
<button class="s" type="submit">
<i class="fa fa-search"></i>
</button>
</form>
</div>`)
	} else {
		html.WriteString(`<div style="flex:1"></div>`)
	}
	html.WriteString(`<div class="view-buttons">
<a href="javascript:void(0);" class="cards view-button" title="Cards">
<i class="fa fa-th"></i>
</a>
<a href="javascript:void(0);" class="list view-button" title="List">
<i class="fa fa-bars"></i>
</a>
</div>
</div>`)
	html.WriteString(content)
	html.WriteString(`</div>
</div>
<div class="footer section">
<div>BookBrowser ` + curversion + `</div>
<div>Copyright 2017 <a href="https://geek1011.github.io">Patrick G</a></div>
</div>
<script src="/static/view.js"></script>
<script>
    BookBrowserVersion = "` + curversion + `";
</script>
<script src="/static/picomodal.js"></script>
<script src="/static/updater.js"></script>
</body>
</html>`)
	return html.String()
}

func bookHTML(b *Book, isInfo bool) string {
	var html bytes.Buffer
	if isInfo {
		html.WriteString(`<div class="book info">`)
	} else {
		html.WriteString(`<div class="book">`)
	}

	html.WriteString(`<a class="cover" href="/books/` + b.ID + `">`)
	if b.HasCover {
		if !isInfo {
			html.WriteString(`<img alt="cover" src="/covers/` + b.ID + `_thumb.jpg" />`)
		} else {
			html.WriteString(`<img alt="cover" src="/covers/` + b.ID + `.jpg" />`)
		}
	} else {
		html.WriteString(`<img alt="cover" src="/static/nocover.jpg" />`)
	}
	html.WriteString(`</a>`)

	html.WriteString(`<div class="meta">`)
	html.WriteString(`<a class="title" href="/books/` + b.ID + `">`)
	html.WriteString(b.Title)
	html.WriteString(`</a>`)
	if b.Author != "" {
		html.WriteString(`<a class="author" href="/authors/` + b.AuthorID + `">`)
		html.WriteString(b.Author)
		html.WriteString(`</a>`)
	}
	if b.Series.Name != "" {
		html.WriteString(`<div class="series">`)
		html.WriteString(`<a class="name" href="/series/` + b.Series.ID + `">`)
		html.WriteString(b.Series.Name)
		html.WriteString(`</a> - <span class="index">`)
		html.WriteString(fmt.Sprintf("%v", b.Series.Index))
		html.WriteString(`</span>`)
		html.WriteString(`</div>`)
	}
	if isInfo && b.Description != "" {
		html.WriteString(`<div class="description">`)
		html.WriteString(b.Description)
		html.WriteString(`</div>`)
	}

	html.WriteString(`<div class="btn-group">`)
	html.WriteString(`<a class="download btn btn-default" href="/download/` + b.ID + `.` + b.FileType + `">Download ` + strings.ToUpper(b.FileType) + `</a>`)
	if b.FileType == "epub" {
		html.WriteString(`<a class="reader btn btn-default" href="/static/reader/epub/#!/download/` + b.ID + `.` + b.FileType + `">Read</a>`)
	}
	if b.FileType == "pdf" {
		html.WriteString(`<a class="reader btn btn-default" href="/static/reader/pdf/web/viewer.html?file=/download/` + b.ID + `.` + b.FileType + `">Read</a>`)
	}
	html.WriteString(`</div>`)

	html.WriteString(`</div>`)

	html.WriteString(`</div>`)
	return html.String()
}

func bookListPageHTML(books []Book, title string, notfoundtext string, showsearch bool) (html string, notfound bool) {
	if len(books) == 0 {
		return pageHTML("Not Found", notfoundtext, false, false), true
	}

	var booksHTML bytes.Buffer
	booksHTML.WriteString(`<div class="books view cards">`)
	for _, b := range books {
		booksHTML.WriteString(bookHTML(&b, false))
	}
	booksHTML.WriteString(`</div>`)

	return pageHTML(title, booksHTML.String(), true, showsearch), false
}

func itemCardHTML(title string, link string) string {
	var html bytes.Buffer
	html.WriteString(`<a class="item" href="` + link + `">`)
	html.WriteString(title)
	html.WriteString(`</a>`)
	return html.String()
}
