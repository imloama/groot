package main

import "github.com/imloama/groot"

func main() {
	err := groot.NewDefaultServer().Run()
	if err != nil {
		panic(err)
	}
}
