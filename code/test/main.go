package main

import "crypto/rand"
import "fmt"

func main() {
	r := make([]byte, 10)
	_, err := rand.Read(r)
	code := string(r[:])
	if err != nil {
		fmt.Println("err ocured: ", err)
	}
	fmt.Println(code)
}
