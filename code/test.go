package main

import "char/markdown/code/proto"
import "fmt"

func main() {
	a := &proto.Test{
		Name: 1,
		Sex:  "ddd",
	}
	b, err := a.Marshal()
	fmt.Println(b, err)
}
