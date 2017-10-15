package all

import (
	// All the imported formats register themselves with the RegisterFormat func.
	_ "github.com/geek1011/BookBrowser/formats/epub"
	_ "github.com/geek1011/BookBrowser/formats/pdf"
)
