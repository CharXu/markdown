package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	var url string
	url = "https://baidu.com"
	// fmt.Println("Please enter url: ")
	// fmt.Scan(url)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	proto := resp.Proto
	header := resp.Status
	defer resp.Body.Close()
	fmt.Printf("%s\n%\n%v", body, proto, header)
	file, err := os.OpenFile("my.log", os.O_CREATE|os.O_RDWR, 122)
	defer file.Close()
	mylog := log.New(file, "", log.Ldate|log.Ltime|log.Llongfile)
	mylog.Printf("%s", body)

}
