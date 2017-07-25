package main

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/beevik/etree"
)

// PDFMeta is the metadata of a PDF file
type PDFMeta struct {
	Title  string
	Author string
}

// StringBetween gets the string in between two other strings, and returns an empty string if not found. It returns the first match.
func StringBetween(str, start, end string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	e := strings.Index(str, end)
	return str[s:e]
}

// StringAfter gets the string after another.
func StringAfter(str, start string) string {
	s := strings.Index(str, start)
	if s == -1 {
		return ""
	}
	s += len(start)
	return str[s:]
}

// GetPDFMeta gets the metadata of a PDF file.
func GetPDFMeta(fname string) (*PDFMeta, error) {
	c, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	s := string(c)
	c = []byte{}

	s = StringBetween(s, "<?xpacket begin", "</x:xmpmeta>")
	s = StringAfter(s, ">")

	xmp := etree.NewDocument()
	err = xmp.ReadFromString(s)
	if err != nil {
		return nil, err
	}

	v := &PDFMeta{}

	v.Title = filepath.Base(fname)

	for _, e := range xmp.FindElements("//format") {
		// Make sure it is a pdf, not another piece of embedded RDF metadata
		if e.Text() != "application/pdf" {
			s = ""
			debug.FreeOSMemory()

			return v, nil
		}
		break
	}

	for _, e := range xmp.FindElements("//title/Alt/li") {
		v.Title = e.Text()
		break
	}

	v.Author = ""
	for _, e := range xmp.FindElements("//creator/Seq/li") {
		v.Author = e.Text()
		break
	}

	s = ""
	debug.FreeOSMemory()

	return v, nil
}
