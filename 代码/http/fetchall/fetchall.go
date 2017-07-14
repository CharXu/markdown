package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

func main() {

	logfile, err := os.OpenFile("my.log", os.O_CREATE|os.O_RDWR, 111)
	defer logfile.Close()

	start := time.Now()
	ch := make(chan string)
	url := os.Args[1]
	times, err := strconv.Atoi(os.Args[2])
	fmt.Println(times)
	if err != nil {
		fmt.Println(err)
	} else {
		for i := 0; i < times; i++ {
			go fetch(url, ch)
			fmt.Println(<-ch)
		}
	}

	fmt.Printf("%2.fs elapsed\n", time.Since(start).Seconds())
}

func fetch(url string, ch chan string) {
	start := time.Now()
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
		return
	}
	nbytes, err := io.Copy(ioutil.Discard, resp.Body)
	resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("While reading %s: %v", url, err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs  %7d  %s", secs, nbytes, url)
}
