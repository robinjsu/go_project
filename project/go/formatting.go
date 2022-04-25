package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"

	tt "github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func isCommon(s string) bool {
	commonWords := [500]string{"the", "name", "of", "very", "to", "through", "and", "just", "a", "form", "in", "much", "is", "great", "it", "think", "you", "say", "that", "help", "he", "low", "was", "line", "for", "before", "on", "turn", "are", "cause", "with", "same", "as", "mean", "I", "differ", "his", "move", "they", "right", "be", "boy", "at", "old", "one", "too", "have", "does", "this", "tell", "from", "sentence", "or", "set", "had", "three", "by", "want", "hot", "air", "but", "well", "some", "also", "what", "play", "there", "small", "we", "end", "can", "put", "out", "home", "other", "read", "were", "hand", "all", "port", "your", "large", "when", "spell", "up", "add", "use", "even", "word", "land", "how", "here", "said", "must", "an", "big", "each", "high", "she", "such", "which", "follow", "do", "act", "their", "why", "time", "ask", "if", "men", "will", "change", "way", "went", "about", "light", "many", "kind", "then", "off", "them", "need", "would", "house", "write", "picture", "like", "try", "so", "us", "these", "again", "her", "animal", "long", "point", "make", "mother", "thing", "world", "see", "near", "him", "build", "two", "self", "has", "earth", "look", "father", "more", "head", "day", "stand", "could", "own", "go", "page", "come", "should", "did", "country", "my", "found", "sound", "answer", "no", "school", "most", "grow", "number", "study", "who", "still", "over", "learn", "know", "plant", "water", "cover", "than", "food", "call", "sun", "first", "four", "people", "thought", "may", "let", "down", "keep", "side", "eye", "been", "never", "now", "last", "find", "door", "any", "between", "new", "city", "work", "tree", "part", "cross", "take", "since", "get", "hard", "place", "start", "made", "might", "live", "story", "where", "saw", "after", "far", "back", "sea", "little", "draw", "only", "left", "round", "late", "man", "run", "year", "don't", "came", "while", "show", "press", "every", "close", "good", "night", "me", "real", "give", "life", "our", "few", "under", "stop", "open", "ten", "seem", "simple", "together", "several", "next", "vowel", "white", "toward", "children", "war", "begin", "lay", "got", "against", "walk", "pattern", "example", "slow", "ease", "center", "paper", "love", "often", "person", "always", "money", "music", "serve", "those", "appear", "both", "road", "mark", "map", "book", "science", "letter", "rule", "until", "govern", "mile", "pull", "river", "cold", "car", "notice", "feet", "voice", "care", "fall", "second", "power", "group", "town", "carry", "fine", "took", "certain", "rain", "fly", "eat", "unit", "room", "lead", "friend", "cry", "began", "dark", "idea", "machine", "fish", "note", "mountain", "wait", "north", "plan", "once", "figure", "base", "star", "hear", "box", "horse", "noun", "cut", "field", "sure", "rest", "watch", "correct", "color", "able", "face", "pound", "wood", "done", "main", "beauty", "enough", "drive", "plain", "stood", "girl", "contain", "usual", "front", "young", "teach", "ready", "week", "above", "final", "ever", "gave", "red", "green", "list", "oh", "though", "quick", "feel", "develop", "talk", "sleep", "bird", "warm", "soon", "free", "body", "minute", "dog", "strong", "family", "special", "direct", "mind", "pose", "behind", "leave", "clear", "song", "tail", "measure", "produce", "state", "fact", "product", "street", "black", "inch", "short", "lot", "numeral", "nothing", "class", "course", "wind", "stay", "question", "wheel", "happen", "full", "complete", "force", "ship", "blue", "area", "object", "half", "decide", "rock", "surface", "order", "deep", "fire", "moon", "south", "island", "problem", "foot", "piece", "yet", "told", "busy", "knew", "test", "pass", "record", "farm", "boat", "top", "common", "whole", "gold", "king", "possible", "size", "plane", "heard", "age", "best", "dry", "hour", "wonder", "better", "laugh", "true .", "thousand", "during", "ago", "hundred", "ran", "am", "check", "remember", "game", "step", "shape", "early", "yes", "hold", "hot", "west", "miss", "ground", "brought", "interest", "heat", "reach", "snow", "fast", "bed", "five", "bring", "sing", "sit", "listen", "perhaps", "six", "fill", "table", "east", "travel", "weight", "less", "language", "morning", "among"}
	for _, c := range commonWords {
		if s == c {
			return true
		}
	}
	return false
}

type Formatted struct {
	txt    string
	span   fixed.Int26_6
	bounds fixed.Rectangle26_6
}
type Content struct {
	fullText  []byte
	wrapped   []string
	formatted []Formatted
}

func NewContent() *Content {
	c := Content{}
	return &c
}

func (c *Content) parseText(filename string, face font.Face) (int, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("Error reading file! %v\n", err)
		return -1, err
	}
	c.fullText = content
	c.wrapped, c.formatted = formatLines(c.fullText, MAXLINE_TEXT)

	return 0, nil
}

func formatLines(fullText []byte, maxLineW int) ([]string, []Formatted) {
	var lines []string
	var fmtLines []Formatted
	var p []byte
	var lookAhead []byte
	var idx int
	var err error

	buffer := bufio.NewReaderSize(bytes.NewBuffer(fullText), len(fullText))
	length := buffer.Buffered()
	if length < maxLineW+1 {
		lookAhead, err = buffer.Peek(length)
	} else {
		lookAhead, err = buffer.Peek(maxLineW + 1)
	}
	if err != nil {
		panic(err)
	}
	idx = findWrapIdx(lookAhead, maxLineW)
	p = make([]byte, idx, idx)
	n, err := buffer.Read(p)
	if err != nil {
		panic(err)
	}
	ptrim := strings.TrimSuffix(string(p), "\n")
	fmtLines = append(fmtLines, Formatted{txt: ptrim})
	lines = append(lines, ptrim)

	for n != 0 {
		length = buffer.Buffered()
		if length < maxLineW {
			lookAhead, err = buffer.Peek(length)
		} else {
			lookAhead, err = buffer.Peek(maxLineW + 1)
		}
		if err != nil && err != io.EOF {
			panic(err)
		}
		idx = findWrapIdx(lookAhead, maxLineW)
		p = make([]byte, idx, idx)
		n, err = buffer.Read(p)
		if err != nil && err != io.EOF {
			panic(err)
		}
		ptrim := strings.TrimSuffix(string(p), "\n")
		fmtLines = append(fmtLines, Formatted{txt: ptrim})
		lines = append(lines, ptrim)
	}

	return lines, fmtLines
}

func findWrapIdx(b []byte, maxWidth int) int {
	if bytes.ContainsRune(b, rune('\n')) {
		return (bytes.IndexRune(b, rune('\n')) + 1)
	} else if !endsInSpace((b)) {
		return (bytes.LastIndexAny(b, " ") + 1)
	}
	return maxWidth
}

func endsInSpace(lookAhead []byte) bool {
	if len(lookAhead) > 1 {
		lastChar := lookAhead[len(lookAhead)-1]
		secondToLastChar := lookAhead[len(lookAhead)-2]
		return unicode.IsSpace(rune(lastChar)) && unicode.IsSpace(rune(secondToLastChar))
	}
	return true
}

func splitStr(lookup string) []string {
	var list []string
	splitWords := strings.Split(lookup, " ")
	for _, wd := range splitWords {
		word := strings.Trim(wd, " ,.!?';:“”’\"()")
		if !isCommon(word) {
			list = append(list, word)
		}
	}
	return list
}

func parseFont(file string) (*tt.Font, error) {
	ttfFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	ttf, err := tt.Parse(ttfFile)
	if err != nil {
		return nil, err
	}

	return ttf, nil
}

func loadFonts(fontSize float64, fonts ...string) map[string]font.Face {
	fontFaces := make(map[string]font.Face)
	for _, f := range fonts {
		// parse bytes and return a pointer to a Font type object
		fnt, err := parseFont(f)
		if err != nil {
			error.Error(err)
			panic("panic! TTF file not properly loaded")
		}
		// create face, which provides the `glyph mask images`
		face := tt.NewFace(fnt, &tt.Options{
			// options... here just font size (0 is 12-point default)
			Size: fontSize,
		})
		switch {
		case strings.Contains(f, "Regular"):
			fontFaces["regular"] = face
		case strings.Contains(f, "Bold"):
			fontFaces["bold"] = face
		}
	}
	return fontFaces
}

func (word *Word) formatDefs(maxLineW int) Word {
	for i, d := range word.Def {
		s := fmt.Sprintf(" - (%s) %s", d.PartOfSpeech, d.Definition)
		fmtDefs := wrapDef(s, maxLineW)
		word.Def[i].Wrapped = fmtDefs
	}
	return *word
}

func wrapDef(s string, wrapIdx int) []string {
	var lines []string
	if len(s) < wrapIdx {
		return append(lines, s)
	}
	for len(s) > 0 {
		if len(s) < 40 {
			lines = append(lines, s)
			break
		} else {
			lines = append(lines, s[:wrapIdx])
			s = s[wrapIdx:]
		}
	}
	return lines
}
