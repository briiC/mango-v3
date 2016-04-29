package main

import (
	"fmt"

	"bitbucket.org/briiC/mango-v3"
)

func main() {
	ma := mango.NewServer()
	fmt.Println("Start listening on :" + ma.Port)
	panic(ma.Start())
}
