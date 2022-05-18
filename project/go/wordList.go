package main

import (
	"fmt"
	"os"

	"github.com/faiface/gui"
)

func WordList(env gui.Env, save <-chan Word, title string) {
	var wordList []Word
	for {
		select {
		case wordObj := <-save:
			wordList = append(wordList, wordObj)

		case _, ok := <-env.Events():
			if !ok {
				if len(wordList) > 0 {
					writeFile := fmt.Sprintf("%v-WordList.txt", title)
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
			// switch e := e.(type) {
			// case win.MoDown:
			// 	print(e)
			// }
		}
	}
}
