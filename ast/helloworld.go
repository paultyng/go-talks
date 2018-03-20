package main

import "os"

func main() {
	data := []byte("hello world")
	os.Stdout.Write(data)
}
