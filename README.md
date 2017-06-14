# BookBrowser
A easy-to-use tool to generate a web-based epub book browser. All you need to do is download it into the folder with your ebooks, and start it.

## Features
- Search
- Responsive web interface
- Browse by:
    - Author
    - Series (from calibre metadata)
- Web based reader
    - Custom fonts, colors, sizing, spacing
    - Remembers your position
    - And more
- Search
- And more
- Easy-to-use
- Fast
- No extra dependencies

## Screenshots

| ![](screenshots/books-mobile.png) | ![](screenshots/authors-mobile.png) | ![](screenshots/search-mobile.png) | ![](screenshots/book-mobile.png) |
| --- | --- | --- | --- |
| ![](screenshots/books-desktop.png) | ![](screenshots/authors-desktop.png) | ![](screenshots/search-desktop.png) | ![](screenshots/book-desktop.png) |

## System Requirements
The server works on all platforms.

The web interface works on IE 9+, Edge, Firefox 3+, Chrome, Safari 5.1+, Opera 17+, and Android browser 4.4+.

The web-based reader works on IE 10+, Edge, Firefox 28+, Chrome 21+, Safari 9+, Opera 17+, and Android browser 4.4+.

## Usage
Run BookBrowser from the directory with the epub books.

You can also use the command line arguments below:

````
  -addr string
    	The address to bind to. (default ":8090")
  -bookdir string
    	The directory to get books from. This directory must exist. (default ".")
  -tempdir string
    	The directory to use for storing temporary files such as book cover thumbnails. This directory is create on start and deleted on exit. (default "./_temp")
````
