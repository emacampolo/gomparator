package io

import (
	"bufio"
	"io"
	"log"
	"os"
)

func Close(c io.Closer) {
	err := c.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func ReadFile(file *os.File) <-chan string {
	out := make(chan string)
	scanner := bufio.NewScanner(file)
	go func() {
		for scanner.Scan() {
			out <- scanner.Text()
		}
		close(out)
	}()

	return out
}
