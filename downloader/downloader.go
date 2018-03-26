package downloader

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

// ProgressIndicator describes an individual progress indicator on the terminal
type ProgressIndicator struct {
	// position of the cursor to write indicator to
	// relative to top of all progress lines
	CursorPosition int64
	// total bytes we want to read
	TotalBytes int64
	// bytes we've read so far
	BytesRead int64
	// current transfer rate as at last update
	// bytes per second
	TransferRate float64
	// slice of all transfer rates received so far
	TransferRates []float64
	Done          bool
}

// ProgressIndicators describes a collection of ProgressIndicator structs
// , describing the overall set of indicators on the terminal
type ProgressIndicators struct {
	// indicator members
	Members []*ProgressIndicator
	// cursor position of the very top progress line
	TopCursorPosition int64
	Done              bool
}

// UpdateProgress updates a ProgressIndicator struct with the current progress
// information
func (pi *ProgressIndicator) UpdateProgress(read int, rate float64) {
	// update total bytes read to include current read
	pi.BytesRead += int64(read)
	// update transfer rates
	pi.TransferRate = rate
	pi.TransferRates = append(pi.TransferRates, rate)
}

// PrintProgress causes a ProgressIndicator to update the terminal with its
// current progress
func (pi *ProgressIndicator) PrintProgress() {

	if pi.TotalBytes >= pi.BytesRead {
		// fmt.Printf("\r%v KiB/s\n", KiB)
		b := New(int(pi.TotalBytes), pi.TransferRate)
		bars := NewBars(b)

		bars.Render(int(pi.BytesRead))
		if pi.TotalBytes == pi.BytesRead {
			pi.Done = true
		}
	}

}

// PrintProgress causes a ProgressIndicators struct to trigger the printing
// of current progress information for all its members
func (pis *ProgressIndicators) PrintProgress() {

	countDone := 0
	for _, member := range pis.Members {
		if member != nil {
			if member.Done {
				countDone = countDone + 1
			}
			if countDone == len(pis.Members) {
				pis.Done = true
				break
			} else {
				member.PrintProgress()
			}
		}

	}

	if !pis.Done {
		// fmt.Fprintf(os.Stdout, "%c[2K", 27)      // clear the line
		fmt.Printf("\x1b[%dA", len(pis.Members))
	} else {
		fmt.Println(" ")
		log.Println("finished downloading all chunk")
	}

}

var wg sync.WaitGroup

//ByteReader for custom io.Reader
type ByteReader struct {
	io.Reader
	total    int
	progress *ProgressIndicator
	t        time.Time
}

func (b ByteReader) Read(p []byte) (int, error) {
	chunkStartTime := time.Now()
	n, err := b.Reader.Read(p)
	chunkEndTime := time.Now()
	timeDifference := chunkEndTime.Sub(chunkStartTime).Seconds()

	transferSpeed := float64(n) / timeDifference

	b.progress.UpdateProgress(n, transferSpeed)

	return n, err
}

//Start starts the download process
func Start(url string, limit int) error {
	// done := make(chan bool)

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	h, err := http.Head(url)
	if err != nil {
		return err
	}

	contentDisposition := h.Header["Content-Disposition"][0]
	_, params, err := mime.ParseMediaType(contentDisposition)
	if err != nil {
		return err
	}

	filename := params["filename"]
	length, err := strconv.Atoi(h.Header["Content-Length"][0])

	if err != nil {
		log.Printf("%v", err)
		return err
	}

	chunk := length / limit
	diff := length % limit

	// setup progress indicators
	pis := new(ProgressIndicators)
	pis.Members = make([]*ProgressIndicator, limit)

	_, err = os.Create(filename)

	if err != nil {
		return err
	}

	wg.Add(limit)
	errch := make(chan error)

	for i := 0; i < limit; i++ {
		min := chunk * i
		max := chunk * (i + 1)

		if i == limit-1 {
			max += diff
		}

		go func(min int, max int, i int, url string) {

			client := &http.Client{}
			req, _ := http.NewRequest("GET", url, nil)

			rangeHeader := "bytes=" + strconv.Itoa(min) + "-" + strconv.Itoa(max-1)
			req.Header.Add("Range", rangeHeader)

			totalBytes := max - min

			resp, err := client.Do(req)

			if err != nil {
				// log.Println("%v", err)
				errch <- err
				return
			}

			defer resp.Body.Close()

			if err != nil {
				errch <- err
				return
			}

			f, err := os.OpenFile(filename, os.O_RDWR, 0644)

			if err != nil {
				errch <- err
				return
			}
			defer f.Close()

			defer f.Close()

			// setup progress bar
			pi := new(ProgressIndicator)
			pi.TotalBytes = int64(totalBytes)
			pis.Members[int64(i)] = pi

			br := &ByteReader{Reader: resp.Body, progress: pi}

			_, err = f.Seek(int64(min), 0)

			if err != nil {
				errch <- err
				return
			}
			_, err = io.Copy(f, br)
			if err != nil {
				errch <- err
				return
			}

			// wg.Done()
			select {
			case <-errch:
				return
			default:
				wg.Done()
			}

			// log.Printf("%v %v", i, " chunk download done")

		}(min, max, i, url)

	}

	stop := make(chan bool, 1)

	wg.Add(1)
	ticker := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				if pis.Done {
					stop <- true
					wg.Done()
				} else {
					if len(pis.Members) == limit {
						pis.PrintProgress()
					}
				}
			case <-c:
				fmt.Printf("\x1b[%dB", len(pis.Members))
				fmt.Printf("Download cancelled!\n")
				os.Exit(1)
			case <-stop:
				log.Println("debug stop")
				return
			}
		}
	}()

	//TODO handle error properly from error channel

	wg.Wait()

	log.Println("download done!")

	return nil

}
