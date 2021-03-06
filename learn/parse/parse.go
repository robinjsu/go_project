package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"unicode"
)

const (
	maxLineW = 150
)

type Sentence struct {
	line      string
	numBytes  int
	numPixels int
}

type Pgraph struct {
	// single paragraph, split into slice of strings split at each space
	lines []Sentence
	// number of bytes in the entire string
	num int
}
type Content struct {
	fullText []byte
	pgraph   []Pgraph
	num      int
}

func NewContent() *Content {
	c := Content{}
	return &c
}

func (c *Content) Store(pg []byte) error {
	newPg := Pgraph{}
	line := string(pg)
	sentences := strings.SplitAfter(line, ". ")
	for _, lin := range sentences {
		sentence := Sentence{line: lin, numBytes: len(lin), numPixels: 0}
		newPg.lines = append(newPg.lines, sentence)
	}
	newPg.num = len(sentences)
	c.pgraph = append(c.pgraph, newPg)
	c.num++
	return nil
}

type Formatted struct {
	txt  string
	span int
}

func endsInSpace(lookAhead []byte) bool {
	lastChar := lookAhead[maxLineW]
	secondToLastChar := lookAhead[maxLineW-1]
	return unicode.IsSpace(rune(lastChar)) && unicode.IsSpace(rune(secondToLastChar))
}

func formatLines(c *Content) []Formatted {
	var fmtLines []Formatted
	var p []byte
	var idx int

	buffer := bufio.NewReader(bytes.NewBuffer(c.fullText))
	lookAhead, err := buffer.Peek(maxLineW + 1)
	if err != nil {
		panic(err)
	}
	if !endsInSpace(lookAhead) {
		idx = bytes.LastIndexAny(lookAhead, " ")
	} else {
		idx = maxLineW
	}
	p = make([]byte, idx, idx)
	_, err = buffer.Read(p)
	if err != nil {
		panic(err)
	}
	fmtLines = append(fmtLines, Formatted{txt: string(lookAhead[0:idx]), span: idx})

	for i := 0; i < 10; i++ {
		lookAhead, err := buffer.Peek(maxLineW + 1)
		if err != nil {
			panic(err)
		}
		if !endsInSpace(lookAhead) {
			idx = bytes.LastIndex(lookAhead, []byte(" "))
		} else {
			idx = maxLineW
		}
		p = make([]byte, idx, idx)
		_, err = buffer.Read(p)
		if err != nil {
			panic(err)
		}
		fmtLines = append(fmtLines, Formatted{txt: string(p), span: idx})
	}

	return fmtLines
}

func main() {
	file := "./alice.txt"
	c := NewContent()
	content, err := os.ReadFile(file)
	if err != nil {
		fmt.Printf("Error not nil! %v\n", err)
	}
	c.fullText = content
	lines := formatLines(c)
	for _, l := range lines {
		fmt.Printf("%v %d\n\n", l.txt, l.span)
	}

	// buffer := bytes.NewBuffer(content)
	// p, err := buffer.ReadBytes('\n')
	// c.Store(p)
	// for p != nil {
	// 	p, err = buffer.ReadBytes('\n')
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		fmt.Printf("Error! %v\n", err)
	// 	}
	// 	c.Store(p)
	// }
	// fmt.Println(c.pgraph[7])
}
