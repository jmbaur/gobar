package module

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/jmbaur/gobar/color"
	"github.com/jmbaur/gobar/i3"
	"golang.org/x/sys/unix"
)

type Battery struct {
	Name string
}

func getFileContents(f *os.File) (string, error) {
	defer f.Seek(0, io.SeekStart)
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bytes.Trim(data, "\n")), nil
}

func (b Battery) Run(c chan Update, position int) {
	capacityFile, err := os.Open(fmt.Sprintf("/sys/class/power_supply/%s/capacity", b.Name))
	if err != nil {
		log.Println(err)
		return
	}

	epfd, err := unix.EpollCreate(5)
	if err != nil {
		log.Println(err)
		return
	}

	events := make([]unix.EpollEvent, 1, 5)
	events[0] = unix.EpollEvent{
		Events: unix.EPOLLPRI, /* | unix.EPOLLHUP | unix.EPOLLERR */
		Fd:     int32(capacityFile.Fd()),
	}

	err = unix.EpollCtl(epfd, unix.EPOLL_CTL_ADD, int(capacityFile.Fd()), &events[0])
	if err != nil {
		log.Println(err)
		return
	}

	bs := make([]byte, 100, 100)

	for {
		col := color.Normal
		log.Println("waiting")
		n, err := unix.EpollWait(epfd, events, -1)
		if err != nil {
			log.Println(err)
			continue
		}

		for _, e := range events[:n] {
			log.Printf("%+v\n", e)
		}

		r, err := capacityFile.Read(bs)
		if err != nil {
			log.Println(err)
			continue
		}
		// seek to beginning to relatch event
		capacityFile.Seek(0, io.SeekStart)
		capacity, err := strconv.Atoi(string(bytes.TrimSpace(bs[0:r])))
		if err != nil {
			log.Println(err)
			continue
		}

		switch true {
		case capacity > 80:
			col = color.Green
		case capacity < 20:
			col = color.Red
		}

		log.Println("updating battery")
		c <- Update{
			Block: i3.Block{
				FullText: fmt.Sprintf("%s: %d%%", b.Name, capacity),
				Color:    col,
			},
			Position: position,
		}
	}
}
