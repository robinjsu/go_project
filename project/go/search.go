package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"strings"

	"github.com/faiface/gui"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

const (
	MIN_X = 900
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

// TODO: a bit more string clean-up to do
// TODO: display words to side
// TODO pull definitions

// https://github.com/faiface/gui/blob/master/examples/imageviewer/util.go#L66
func drawText(s string, face font.Face) image.Image {
	text := &font.Drawer{
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(0, face.Metrics().Height.Ceil()),
	}
	txtBnds, _ := text.BoundString(s)
	bounds := image.Rect(
		txtBnds.Min.X.Floor(),
		txtBnds.Min.Y.Floor(),
		txtBnds.Max.X.Ceil(),
		txtBnds.Max.Y.Ceil(),
	)
	text.Dst = image.NewRGBA(bounds)
	text.DrawString(s)
	return text.Dst
}

func displayWords(wordList []string, face font.Face, bounds image.Rectangle) func(draw.Image) image.Rectangle {
	var textImages []image.Image
	for _, w := range wordList {
		img := drawText(w, face)
		textImages = append(textImages, img)
	}
	searchBar := func(drw draw.Image) image.Rectangle {
		newR := bounds
		draw.Draw(drw, newR, &image.Uniform{color.RGBA{0, 150, 100, 255}}, image.ZP, draw.Over)
		y := face.Metrics().Height.Ceil() * 2
		for i, img := range textImages {
			x1 := img.Bounds().Dx()
			fmt.Println(x1, y)
			fontR := image.Rect(MIN_X, (y * i), (x1 + MIN_X), (y * (i + 1)))
			// padded := fRect.Inset(-2)
			draw.Draw(drw, fontR, img, image.ZP, draw.Over)
		}
		return newR
	}
	return searchBar
}

func Search(env gui.Env, fontFaces map[string]font.Face, words <-chan string) {

	for {
		select {
		case lookup := <-words:
			// fmt.Println(lookup)
			splitWords := strings.Split(lookup, " ")
			var list []string
			for _, wd := range splitWords {
				word := strings.Trim(wd, " ,.!?';:“”’\"()")
				if !isCommon(word) {
					list = append(list, word)
				}
			}
			wordCorner := image.Rect(900, 0, 1200, 450)
			env.Draw() <- displayWords(list, fontFaces["regular"], wordCorner)
		case _, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			// switch e := e.(type) {
			// case win.MoDown:
			// 	// fmt.Println(e.X, e.Y)
			// case win.MoUp:
			// }
		}
	}
}
