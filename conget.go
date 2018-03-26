package main

import (
	"fmt"
	"log"
	"os"

	"github.com/monirz/conget/downloader"
	"github.com/urfave/cli"
)

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	app := cli.NewApp()
	app.Name = "conget"
	app.Usage = "A concurrent file downloader."
	app.Author = "monir.smith@gmail.com - github.com/monirz"
	app.Version = "1.0.0"

	//---------------------------------------

	// setup cli flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Usage: "URL of the file.",
		},
		cli.IntFlag{
			Name:  "cuncurrency, c",
			Value: 10,
			Usage: "Number of concurrent process to be runned.",
		},
	}

	// url := "http://trailers.divx.com/divx_prod/profiles/Fashion_DivX720p_ASP.divx"
	var url string
	var concurrentProcessNumber int

	app.Action = func(c *cli.Context) error {
		if c.String("url") == "" {
			// Check if there is an update
			fmt.Printf("No URL provided! Ex. conget -u http://example.com/example.mp4 \n")

			os.Exit(1)
		} else {
			url = c.String("url")
		}

		if c.Int("cuncurrency") > 0 {
			concurrentProcessNumber = c.Int("cuncurrency")
		} else {
			concurrentProcessNumber = 10
		}
		return nil
	}
	app.Run(os.Args)
	fmt.Println("")

	// log.Fatal(url)

	//---------------------------
	err := downloader.Start(url, concurrentProcessNumber)
	if err != nil {
		log.Fatal("error ", err)
	}
}
