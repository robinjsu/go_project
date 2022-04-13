package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type WordDef struct {
	Word string              `json:"word"`
	Def  []map[string]string `json:"definitions"`
}

type Definition struct {
	Definition   string
	PartOfSpeech string
}

func main() {
	DICT_KEY := os.Getenv("DICT_API_KEY")

	lookup := "blubber"
	url := fmt.Sprintf("https://wordsapiv1.p.rapidapi.com/words/%s/definitions", lookup)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("X-RapidAPI-Host", "wordsapiv1.p.rapidapi.com")
	req.Header.Add("X-RapidAPI-Key", DICT_KEY)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	var worddef WordDef
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&worddef)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i, def := range worddef.Def {
		fmt.Printf("%d - Def: %v, PoS: %v\n", i, def["definition"], def["partOfSpeech"])
	}
	// err = json.Unmarshal(jsonBody, &definitions)
	// if err != nil {
	// 	fmt.Print(err)
	// 	return
	// }
	// for def := range definitions["word"]
	// fmt.Println(res)
	// fmt.Println(string(body))
}
