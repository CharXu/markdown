package main

import (
	"fmt"
	"log"

	"char/markdown/code/protobuf/example"

	"github.com/golang/protobuf/proto"
)

func main() {
	test := &example.Test{
		Label: "Hello",
		Type:  17,
		OptionalGroup: &example.Test_Group{
			RequiedField: "good bye",
		},
	}

	data, err := proto.Marshal(test)
	if err != nil {
		fmt.Println("Marshal err :", err)
	}

	newTest := &example.Test{}
	err = proto.Unmarshal(data, newTest)
	if err != nil {
		log.Fatal("unmarshaling error: ", err)
	}
	if test.GetLabel() != newTest.GetLabel() {
		log.Fatalf("data mistmatch %q != %q", test.GetLabel(), newTest.GetLabel())
		return
	}
	fmt.Print("Label is :", test.GetLabel())
}
