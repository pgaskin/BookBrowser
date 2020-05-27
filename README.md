# BookBrowser
[![Build Status](https://travis-ci.org/geek1011/BookBrowser.svg?branch=master)](https://travis-ci.org/geek1011/BookBrowser)

**Note:** This project is no longer maintained, as I haven't had the time or motivation to continue working on it. BookBrowser will still work as-is (I still use it myself occasionally), but is unlikely to receive any new features or bugfixes.

An easy-to-use tool to generate a web-based ePub and PDF ebook browser. All you need to do is [download it](https://github.com/geek1011/BookBrowser/releases/latest) into the folder with your ebooks, and run it. There is also a [demo](https://bookbrowser-demo.geek1011.net/books/).

## Features
- Multiple book formats
    - epub
    - pdf
    - mobi (basic support)
- Search
- Advanced Search
    - Search any combination of fields
    - View all information in the results
- List view
- Responsive web interface
- Update notifications
- Browse by:
    - Author
    - Series (from calibre metadata)
- Sorted by:
    - Last added
    - Alphabetically
    - Series
- Web based reader
    - Custom fonts, colors, sizing, spacing
    - Remembers your position
    - Book search
    - And more
- Search
- And more
- Easy-to-use
- Fast
- No extra dependencies

## Screenshots

| ![](docs/screenshots/books-mobile.png) | ![](docs/screenshots/books-list-mobile.png) | ![](docs/screenshots/authors-mobile.png) | ![](docs/screenshots/book-mobile.png) |
| --- | --- | --- | --- |
| ![](docs/screenshots/books-desktop.png) | ![](docs/screenshots/books-list-desktop.png) | ![](docs/screenshots/authors-desktop.png) | ![](docs/screenshots/book-desktop.png) |

## Reader Screenshots

| Desktop | Mobile |
| --- | --- |
| ![](docs/screenshots/reader-desktop.png) | ![](docs/screenshots/reader-mobile.png) |

## Advanced Search

| ![](docs/screenshots/list-desktop.png) |
| --- |
| |

## System Requirements
The server works on all platforms.

The web interface works on IE 9+, Edge, Firefox 3+, Chrome, Safari 5.1+, Opera 17+, and Android browser 4.4+.

The web-based reader works on IE 10+, Edge, Firefox 28+, Chrome 21+, Safari 9+, Opera 17+, and Android browser 4.4+.

## Usage

```
Usage: BookBrowser [OPTIONS]

Options:
  -a, --addr string      the address to bind the server to ([IP]:PORT) (default ":8090")
  -b, --bookdir string   the directory to load books from (must exist) (default "/home/patrick/src/BookBrowser")
  -h, --help             Show this help text
  -n, --nocovers         do not index covers
  -t, --tempdir string   the directory to store temp files such as cover thumbnails (created on start, deleted on exit unless already exists) (default "/tmp/bookbrowser946254949")
      --version          Show the version
```
