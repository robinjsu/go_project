package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

type Pgraph struct {
	// single paragraph, split into slice of strings split at each space
	lines [][]string
	// number of bytes in the entire string
	num int
}
type Content struct {
	pgraph []Pgraph
	num    int
}

// func (c *Content) add(st []string) {
// 	new_sent :=
// 	c.pgraph = append(c.pgraph)
// }

func NewContent() *Content {
	c := Content{}
	return &c
}

func (c *Content) Store(pg []byte) error {
	newPg := Pgraph{}
	line := string(pg)
	// t := regexp.MustCompile(`[.|?|!]`)
	sentences := strings.SplitAfter(line, `[.|?|!]`)
	// splitAt := func(c rune) bool {
	// 	return c == '.' || c == '?' || c == '!'
	// }
	// sentences := strings.FieldsFunc(line, splitAt)
	for _, lin := range sentences {
		words := strings.Fields(lin)
		newPg.lines = append(newPg.lines, words)
		newPg.num = len(pg)
	}
	c.pgraph = append(c.pgraph, newPg)
	c.num++
	return nil
}

func main() {
	file := "./test.txt"
	c := NewContent()
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error not nil! %v\n", err)
	}

	buffer := bytes.NewBuffer(content)
	p, err := buffer.ReadBytes('\n')
	c.Store(p)
	for p != nil {
		p, err = buffer.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("Error! %v\n", err)
		}
		c.Store(p)
	}
	fmt.Print(c.pgraph[0])
}
