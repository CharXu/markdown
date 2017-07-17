package main

import (
	"example"
	"github.com/golang/protobuf/proto"
	"fmt"
)

func main() {
	test := &example.Test{
		Label: proto.String("Hello"),
		Type: proto.Int32(17)
		OptionalGroup: &example.Test_Group {
			RequiredField: proto.String("good bye"),
		},
	}

	data, err := proto.Marshal(test)
	if err != nil {
		fmt.Println("Marshal err :", err)
	}

	newTest
}
