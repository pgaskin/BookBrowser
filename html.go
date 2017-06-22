package main

import (
	"bytes"
	"fmt"
)

func pageHTML(title string, content string) string {
	var html bytes.Buffer
	html.WriteString(`<!DOCTYPE html>
<html>
<head>
<title>`)
	html.WriteString(title)
	html.WriteString(`</title>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<meta http-equiv="X-UA-Compatible" content="ie=edge">
<link href="https://fonts.googleapis.com/css?family=Open+Sans:300,400,400i,600,700" rel="stylesheet">
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
<link rel="stylesheet" href="/static/style.css">
</head>
<body>
<nav class="navbar navbar-default navbar-fixed-top">
	<div class="navbar-header">
		<button type="button" class="navbar-toggle collapsed" data-toggle="collapse" data-target="#navbar" aria-expanded="false" aria-controls="navbar">
		<span class="sr-only">Toggle navigation</span>
		<span class="icon-bar"></span>
		<span class="icon-bar"></span>
		<span class="icon-bar"></span>
		</button>
		<a class="navbar-brand" href="/books/">BookBrowser</a>
	</div>
	<div id="navbar" class="navbar-collapse collapse">
		<ul class="nav navbar-nav">
		<li><a href="/books/">Books</a></li>
		<li><a href="/authors/">Authors</a></li>
		<li><a href="/series/">Series</a></li>
		<li><a href="/random/">Random Book</a></li>
		</ul>
		<ul class="nav navbar-nav navbar-right" style="padding-right:20px">
			<form class="navbar-form" role="search" method="GET" action="/search/">
			<div class="input-group">
				<input type="text" class="form-control" placeholder="Search" name="q" id="q">
				<div class="input-group-btn">
					<button class="btn btn-default" type="submit"><i class="glyphicon glyphicon-search"></i></button>
				</div>
			</div>
			</form>
		</ul>
	</div>
</nav>
<div class="container">
<div class="page-header">
<h1>`)
	html.WriteString(title)
	html.WriteString(`</h1>
</div>`)
	html.WriteString(content)
	html.WriteString(`</div>
<footer class="footer">
	<div class="container">
		<p class="text-muted"><a href="https://github.com/geek1011/BookBrowser">BookBrowser</a> ` + curversion + `</p> 
		<p class="text-muted">Copyright 2017 <a href="https://geek1011.github.io">Patrick G</a></p>
	</div>
</footer>
<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/3.2.1/jquery.min.js"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/js/bootstrap.min.js" integrity="sha384-Tc5IQib027qvyjSMfHjOMaLkfuWVxZxUPnCJA7l2mCWNIpG9mGCD8wGNIcPD7Txa" crossorigin="anonymous"></script>
</body>
</html>`)
	return html.String()
}

func bookHTML(b *Book, isCard bool) string {
	var html bytes.Buffer
	if isCard {
		html.WriteString(`<div class="book card">`)
	} else {
		html.WriteString(`<div class="book info">`)
	}

	html.WriteString(`<a class="cover" href="/books/` + b.ID + `">`)
	if b.HasCover {
		if isCard {
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
	if !isCard && b.Description != "" {
		html.WriteString(`<div class="description">`)
		html.WriteString(b.Description)
		html.WriteString(`</div>`)
	}
	html.WriteString(`<a class="download btn btn-default" href="/download/` + b.ID + `.epub">Download EPUB</a>`)
	html.WriteString(`<a class="reader btn btn-default" href="/static/reader/#/download/` + b.ID + `.epub">Read</a>`)
	html.WriteString(`</div>`)

	html.WriteString(`</div>`)
	return html.String()
}

func itemCardHTML(title string, description string, link string) string {
	var html bytes.Buffer
	html.WriteString(`<a class="item card" href="` + link + `">`)
	html.WriteString(`<div class="title">`)
	html.WriteString(title)
	html.WriteString(`</div>`)
	html.WriteString(`<div class="description">`)
	html.WriteString(description)
	html.WriteString(`</div>`)
	html.WriteString(`</a>`)
	return html.String()
}
