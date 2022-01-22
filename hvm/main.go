package main

import (
	"../hvmlib"
	"fmt"
	"os"
)

const (
	TITLE = "(\\_/)\n(^.^)\n(\")(\")\nHVM : THE HONNY VIRTUAL MACHINE\nVer. 1.0"
	USAGE = "USAGE : hvm.bin filename"
)

func main() {

	fmt.Println(TITLE)

	if len(os.Args) != 2 {
		fmt.Println(USAGE)
		return
	}

	h, err := hvmlib.Create(os.Args[1])

	if err != nil {
		fmt.Println(err)
		return
	}

	err = h.Execute()

	if err != nil {
		fmt.Println(err)
		return
	}

}
