package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"net/http"
	"os"
	"strings"

	"golang.org/x/image/font"
)

type Word struct {
	Word string       `json:"word"`
	Def  []Definition `json:"definitions"`
}

func (word *Word) formatDefs(face font.Face, areaR *image.Rectangle) Word {
	for i, d := range word.Def {
		s := fmt.Sprintf("[%s] %s", d.PartOfSpeech, d.Definition)
		fmtDefs := wrapDef(s, calculateLineWidth(face, areaR.Dx()))
		word.Def[i].Wrapped = fmtDefs
	}
	return *word
}

func (word *Word) String() string {
	var wl []string
	for _, d := range word.Def {
		f := fmt.Sprintf("(%s) %s", d.PartOfSpeech, d.Definition)
		wl = append(wl, f)
	}
	defs := strings.Join(wl, "; ")
	w := fmt.Sprintf("%s: %s", word.Word, defs)
	return w
}

type Definition struct {
	PartOfSpeech string `json:"partOfSpeech"`
	Definition   string `json:"definition"`
	Wrapped      []string
}

func getDef(lookup string) (Word, error) {
	DICT_KEY := os.Getenv("DICT_API_KEY")
	if DICT_KEY == "" {
		return Word{}, errors.New("API key missing")
	}

	url := fmt.Sprintf("https://wordsapiv1.p.rapidapi.com/words/%s/definitions", lookup)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Word{}, HttpError{
			method: "NewRequest",
			url:    url,
			err:    err,
		}
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Host", "wordsapiv1.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", DICT_KEY)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return Word{}, HttpError{
			method: "GET",
			url:    url,
			err:    err,
		}
	} else if res.StatusCode != 200 {
		return Word{}, HttpError{
			method: "Response",
			url:    url,
			err:    errors.New(fmt.Sprintf("received non-200 status code: %v", res.StatusCode)),
		}
	}
	defer res.Body.Close()

	var worddef Word
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&worddef)
	if err != nil {
		return Word{}, errors.New(fmt.Sprintf("error in decoding json: %s", err))
	}

	return worddef, nil
}
