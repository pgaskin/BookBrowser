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
<meta name="apple-mobile-web-app-capable" content="yes">
<meta name="apple-mobile-web-app-status-bar-style" content="black">
</head>
<body>
<nav class="navbar navbar-default navbar-fixed-top">
	<div class="container">
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
<script>
    BookBrowserVersion = "` + curversion + `";

	var cachedXHR = function(url, cacheSeconds, callback, errorCallback) {
		var cacheKey = "xhrcache|" + url + "|";
		var cacheTimeKey = cacheKey + "time";
		var cacheValueKey = cacheKey + "value";

		var currentTime = Math.round(new Date().getTime() / 1000);
		var cacheValue = localStorage.getItem(cacheValueKey);
		var cacheTime = 0;
		
		try {
			cacheTime = parseInt(localStorage.getItem(cacheTimeKey));
		} catch (e) {
			localStorage.setItem(cacheTimeKey, 0);
			cacheTime = 0;
		}
		
		if (cacheValue === null || (currentTime - cacheTime) > cacheSeconds) {
			var xhttp = new XMLHttpRequest();
			xhttp.onreadystatechange = function() {
				if (this.readyState == 4 && this.status == 200) {
					var resp = this.responseText;
					localStorage.setItem(cacheTimeKey, currentTime);
					localStorage.setItem(cacheValueKey, resp);
					callback(resp);
				} else if (this.readyState == 4 && this.status > 399) {
					errorCallback("Error: HTTP status " + this.status.toString());
				}
			};
			xhttp.onerror = function() {
				errorCallback("Error: Network error");
			};
			xhttp.open("GET", url, true);
			xhttp.send();
		} else {
			callback(cacheValue);
		}
	};
	var HOUR = 3600;
	cachedXHR("https://api.github.com/repos/geek1011/BookBrowser/releases", HOUR / 2, function(respa) {
		cachedXHR("https://api.github.com/repos/geek1011/BookBrowser/releases/latest", HOUR / 2, function(respb) {
			try {
				var releases = JSON.parse(respa);
				var current = JSON.parse(respb);
				
				var currentVersion = BookBrowserVersion;

				var isDev = (currentVersion.indexOf("+")>-1);
				if (isDev) {
					console.warn("You are using a development version of BookBrowser");
					return;
				}
				
				var latestVersion = current["tag_name"];
				if (latestVersion == currentVersion) {
					console.info("You are using the latest version of BookBrowser: " + latestVersion);
					return;
				}

					console.info("You are not using the latest version of BookBrowser. Your current version is " + currentVersion + ", but the latest version is " + latestVersion);
					var releaseNotes = "";
					for (var i = 0; i < releases.length; i++) {
						var release = releases[i];
						if (release["tag_name"] == currentVersion) {
							break;
						}
						releaseNotes += "<b>" + release["tag_name"] + "</b><br><br>" + release.body.split("## Usage")[0].split("\n").filter(function(l) {
							return !(l.indexOf("Changes for")>-1) && (l !== "");
						}).map(function(l) {
							return l + "<br>";
						}).join("\n") + "<br><br><br>";
					}

					var message = "<b>You are not using the latest version of BookBrowser. Your current version is " + currentVersion + ", but the latest version is " + latestVersion + ".</b><br><br>You can download the latest version <a href=\"https://github.com/geek1011/BookBrowser/releases/latest\" target=\"_blank\">here</a>.<br><br>The release notes for the versions up to " + latestVersion + " are below.<br><br>";
					message += releaseNotes;
					
					console.log(message);

					var modalHTML = '<div class="modal fade" id="updateModal" tabindex="-1" role="dialog"> <div class="modal-dialog" role="document"> <div class="modal-content"> <div class="modal-header"> <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button> <h4 class="modal-title" id="updateModalLabel">BookBrowser Update Available</h4> </div> <div class="modal-body">' + message + '</div> <div class="modal-footer"> <button type="button" class="btn btn-default" data-dismiss="modal">Close</button> <a type="button" target="_blank" href="https://github.com/geek1011/BookBrowser/releases/latest" class="btn btn-primary">Update</a> </div> </div> </div> </div>';

					document.body.appendChild(document.createElement("div")).innerHTML = modalHTML;

					$('#updateModal').modal('show')
			} catch (err) {
				console.warn(err);
			}
		}, function(err){
			console.warn(err);
		});
	}, function(err){
		console.warn(err);
	});
</script>
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
	html.WriteString(`<a class="reader btn btn-default" href="/static/reader/#!/download/` + b.ID + `.epub">Read</a>`)
	html.WriteString(`</div>`)

	html.WriteString(`</div>`)
	return html.String()
}

func bookListPageHTML(books []Book, title string, notfoundtext string) (html string, notfound bool) {
	if len(books) == 0 {
		return pageHTML("Not Found", notfoundtext), true
	}

	var booksHTML bytes.Buffer
	booksHTML.WriteString(`<div class="books cards">`)
	for _, b := range books {
		booksHTML.WriteString(bookHTML(&b, true))
	}
	booksHTML.WriteString(`</div>`)

	return pageHTML(title, booksHTML.String()), false
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
