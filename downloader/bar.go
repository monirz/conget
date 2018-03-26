package downloader

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type Bar struct {
	Val, Max     int
	Override     string
	Name         string
	NamePadding  int
	TransferRate float64
}

type tsize struct {
	rows uint16
	cols uint16
	pixW uint16
	pixH uint16
}

func getTermSize() tsize {
	size := tsize{}
	_, _, err := syscall.Syscall(syscall.SYS_IOCTL, os.Stdout.Fd(), syscall.TIOCGWINSZ, uintptr(unsafe.Pointer(&size)))
	if err != 0 {
		panic(err)
	}
	return size
}

func New(max int, transferRate float64) *Bar {
	return &Bar{
		Max:          max,
		TransferRate: transferRate,
	}
}

func (b *Bar) generate(width int) string {
	percentage := fmt.Sprintf("%3d%s", int(float64(b.Val)/float64(b.Max)*100.0), "%")

	leftPadding := 2

	if leftPadding < 1 {
		leftPadding = 1
	}
	if len(b.Override) > 0 {
		return percentage + strings.Repeat(" ", leftPadding+8) + b.Override
	}

	counts := fmt.Sprintf("%3d", b.Val)
	KiB := int64(b.TransferRate) / 1024
	// log.Println(b.TransferRate)
	tr := fmt.Sprintf("%10d%s", KiB, "KB/s")
	barwidth := width - 60

	if barwidth < 0 {
		barwidth = 0
	}

	part := int(float64(b.Val) / float64(b.Max) * float64(barwidth))

	if part < 0 {
		part = 0
	}
	if part > barwidth {
		part = barwidth
	}

	return percentage + strings.Repeat(" ", leftPadding) + " [" + strings.Repeat("=", part) + ">" + strings.Repeat(" ", barwidth-part) + "] " + counts + strings.Repeat(" ", 4) + tr
}

func (b *Bar) Render() {
	size := getTermSize()
	fmt.Printf("\r%s", b.generate(int(size.cols)))
}

type Bars struct {
	bars     []*Bar
	rendered bool
}

func NewBars(bars ...*Bar) *Bars {
	b := &Bars{bars: bars}
	return b
}

func (b *Bars) Render(val int) {
	// if b.rendered {
	// 	fmt.Printf("\x1b[%dA", len(b.bars))
	// }
	maxLen := 0
	for _, bar := range b.bars {
		if len(bar.Name) > maxLen {
			maxLen = len(bar.Name)
		}
	}
	size := getTermSize()
	for _, bar := range b.bars {
		bar.NamePadding = maxLen + 1
		bar.Val = val
		fmt.Printf("\r%s", bar.generate(int(size.cols)))
		fmt.Println()
	}
	b.rendered = true
}

func (b *Bars) Add(bar *Bar) {
	b.bars = append(b.bars, bar)
	if b.rendered {
		fmt.Println()
	}
}
