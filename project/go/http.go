package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
)

type WordDef struct {
	Word string       `json:"word"`
	Def  []Definition `json:"definitions"`
}

type Definition struct {
	Definition   string `json:"definition"`
	PartOfSpeech string `json:"partOfSpeech"`
}

func getDef(lookup string) (WordDef, error) {
	DICT_KEY := os.Getenv("DICT_API_KEY")
	if DICT_KEY == "" {
		return WordDef{}, errors.New("no api key provided")
	}

	url := fmt.Sprintf("https://wordsapiv1.p.rapidapi.com/words/%s/definitions", lookup)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return WordDef{}, err
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Host", "wordsapiv1.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", DICT_KEY)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return WordDef{}, err
	} else if res.StatusCode != 200 {
		return WordDef{}, errors.New("received a non-200 status code")
	}
	defer res.Body.Close()

	var worddef WordDef
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&worddef)
	if err != nil {
		return WordDef{}, err
	}
	return worddef, nil
}
