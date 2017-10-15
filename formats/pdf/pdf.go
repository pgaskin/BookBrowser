package pdf

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/debug"
	"time"

	"github.com/beevik/etree"
	"github.com/geek1011/BookBrowser/formats"
	"github.com/geek1011/BookBrowser/models"
	"github.com/geek1011/BookBrowser/modules/util"
)

func indexer(filename string) (book *models.Book, cover *models.Cover, err error) {
	defer func() {
		if r := recover(); r != nil {
			book = nil
			cover = nil
			err = fmt.Errorf("Unknown error parsing book. Skipping. Error: %s", r)
		}
	}()

	var title string
	var author string
	var modtime time.Time

	if file, err := os.Stat(filename); err == nil {
		modtime = file.ModTime()
	}

	c, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	s := string(c)
	c = []byte{}

	s = util.StringBetween(s, "<?xpacket begin", "</x:xmpmeta>")
	s = util.StringAfter(s, ">")

	xmp := etree.NewDocument()
	err = xmp.ReadFromString(s)
	if err != nil {
		return nil, nil, err
	}

	title = filepath.Base(filename)

	for _, e := range xmp.FindElements("//format") {
		// Make sure it is a pdf, not another piece of embedded RDF metadata
		if e.Text() != "application/pdf" {
			s = ""
			debug.FreeOSMemory()

			return models.NewBook(title, author, "", "", 0, "", filename, false, modtime, "pdf"), nil, nil
		}
		break
	}

	for _, e := range xmp.FindElements("//title/Alt/li") {
		title = e.Text()
		break
	}

	for _, e := range xmp.FindElements("//creator/Seq/li") {
		author = e.Text()
		break
	}

	s = ""
	debug.FreeOSMemory()

	return models.NewBook(title, author, "", "", 0, "", filename, false, modtime, "pdf"), nil, nil
}

func init() {
	formats.RegisterFormat(&formats.Format{
		Glob:      "**/*.pdf",
		Extension: ".pdf",
		Indexer:   indexer,
	})
}
