package main

import (
	"fmt"
	"path/filepath"
)

func Convert(filename string) {

	filedir := filepath.Dir(filename)
	fmt.Println(filedir)

}
