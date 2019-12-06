package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/monirz/conget/downloader"
	"github.com/spf13/cobra"
)

var (
	url                     string
	concurrentProcessNumber int
	rootCmd                 = &cobra.Command{
		Use:   "conget",
		Short: "Conget is a concurrent file downloader",
		Run:   download,
	}
)

func init() {

	rootCmd.Flags().StringVarP(&url, "url", "u", "", "URL of the source file, ex: http://example.com/file.mp4")
	rootCmd.MarkFlagRequired("url")
	rootCmd.Flags().IntVarP(&concurrentProcessNumber, "concurrent", "c", 5, "Concurrent process number, default is 5")

}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func download(cmd *cobra.Command, args []string) {

	//---------------------------------- start downloading processes ------------------------------------------
	err := downloader.Start(url, concurrentProcessNumber)
	if err != nil {
		log.Fatal("error ", err)
	}
}
