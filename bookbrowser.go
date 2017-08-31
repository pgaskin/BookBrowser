package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/urfave/cli"
)

var curversion = "dev"

// GetIP gets the preferred outbound ip of this machine
func GetIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func main() {
	workdir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Fatal error: %s\n", err)
	}

	tempdir, err := ioutil.TempDir("", "bookbrowser")
	if err != nil {
		tempdir = filepath.Join(workdir, "_temp")
	}

	app := cli.NewApp()
	app.Name = "BookBrowser"
	app.Usage = "Web-based eBook server supporting ePub and PDF."
	app.Version = curversion
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "bookdir, b",
			Value: workdir,
			Usage: "Load books from `DIR`. The directory must exist.",
		},
		cli.StringFlag{
			Name:  "tempdir, t",
			Value: tempdir,
			Usage: "Use `DIR` as the location for storing temporary files such as cover thumbnails. The directory is created on start and deleted on exit.",
		},
		cli.StringFlag{
			Name:  "addr, a",
			Value: ":8090",
			Usage: "`ADDR` is the address to bind the server to. It is in the format IP:PORT. The IP is optional.",
		},
	}
	app.HideHelp = true
	app.Action = func(c *cli.Context) {
		bookdir := c.String("bookdir")
		tempdir := c.String("tempdir")
		addr := c.String("addr")

		log.Printf("BookBrowser %s\n", curversion)

		if _, err := os.Stat(bookdir); err != nil {
			if os.IsNotExist(err) {
				log.Fatalf("Fatal error: book directory %s does not exist\n", bookdir)
			}
		}

		bookdir, err = filepath.Abs(bookdir)
		if err != nil {
			log.Fatalf("Fatal error: Could not resolve book directory %s: %s\n", bookdir, err)
		}

		if _, err := os.Stat(tempdir); os.IsNotExist(err) {
			os.Mkdir(tempdir, os.ModePerm)
		}

		tempdir, err = filepath.Abs(tempdir)
		if err != nil {
			log.Fatalf("Fatal error: Could not resolve temp directory %s: %s\n", tempdir, err)
		}

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigs
			log.Println("Cleaning up covers")
			os.RemoveAll(tempdir)
			os.Exit(0)
		}()

		if !strings.Contains(addr, ":") {
			log.Fatalln("Invalid listening address")
		}

		sp := strings.SplitN(addr, ":", 2)
		if sp[0] == "" {
			ip := GetIP()
			if ip != nil {
				log.Printf("This server can be accessed at http://%s:%s\n", ip.String(), sp[1])
			}
		}

		server := NewServer(addr, bookdir, tempdir, true)
		server.RefreshBookIndex()

		if len(*server.Books) == 0 {
			log.Fatalln("Fatal error: no books found")
		}

		sigsa := make(chan os.Signal, 1)
		signal.Notify(sigsa, syscall.SIGUSR1)
		go func() {
			for _ = range sigsa {
				go func() {
					log.Println("Booklist refresh triggered by SIGUSR1")
					server.RefreshBookIndex()
				}()
			}
		}()

		err = server.Serve()
		if err != nil {
			log.Fatalf("Error starting server: %s\n", err)
		}
	}
	app.Run(os.Args)
}
