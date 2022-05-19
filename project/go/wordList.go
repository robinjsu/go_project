package main

import (
	"fmt"
	"os"
	"strings"

	gui "github.com/faiface/gui"
	win "github.com/faiface/gui/win"
)

func WordList(env gui.Env, save <-chan Word) {
	var title string
	var wordList []Word
	for {
		select {
		case wordObj := <-save:
			wordList = append(wordList, wordObj)

		case e, ok := <-env.Events():
			if !ok {
				if len(wordList) > 0 {
					writeFile := fmt.Sprintf("%s-wordlist.txt", title)
					filePtr, err := os.Create(writeFile)
					defer filePtr.Close()
					if err != nil {
						panic(FileError{
							filename: writeFile,
							err:      err,
						})
					}
					filePtr.Chmod(os.ModeAppend)
					for _, w := range wordList {
						word := fmt.Sprintf("%v\n", w.String())
						_, err := filePtr.WriteString(word)
						if err != nil {
							panic(FileError{
								filename: writeFile,
								err:      err,
							})
						}
					}
				}
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.PathDrop:
				dirs := strings.Split(e.FilePath, "/")
				title = strings.TrimSuffix(dirs[len(dirs)-1], ".txt")
			}
		}
	}
}
