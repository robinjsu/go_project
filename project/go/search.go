package main

import (
	"image"
	"image/draw"
	"strings"

	"github.com/faiface/gui"
	"github.com/faiface/gui/win"
	"golang.org/x/image/font"
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

// TODO: need to refactor so that scopes better suited to the context! having trouble
// with accessing the words and images as they are rendered, need way to access them later
// TODO: a bit more string clean-up to do
// TODO: display words to side
// TODO pull definitions

// set images with their bounds outside of the callback, so that they can be stored in the same context

func displayWords(wordList []string, face font.Face) []imageObj {
	var images []imageObj
	y := face.Metrics().Height.Ceil() * 2

	for i, w := range wordList {
		img, format := drawText(w, face)
		x1 := img.Bounds().Dx()
		fontR := image.Rect(MIN_X, (y * i), (x1 + MIN_X), (y * (i + 1)))
		images = append(images, imageObj{text: format, img: img, placement: fontR})
	}
	return images
}

func drawSearchBar(images []imageObj, bounds *image.Rectangle) func(draw.Image) image.Rectangle {
	searchBar := func(drw draw.Image) image.Rectangle {
		newR := *bounds
		draw.Draw(drw, newR, &image.Uniform{TEAL}, image.ZP, draw.Over)
		for _, obj := range images {
			draw.Draw(drw, obj.placement, obj.img, image.ZP, draw.Over)
		}
		return newR
	}
	return searchBar
}

func highlightWord(images []imageObj, p image.Point, drawDst image.Rectangle, define chan<- string) (func(draw.Image) image.Rectangle, string) {
	var target image.Rectangle
	var lookup string
	for _, img := range images {
		if p.In(img.placement) {
			lookup = img.text.txt
			target = img.placement
		}
	}
	highlight := func(drw draw.Image) image.Rectangle {
		draw.Draw(drw, drawDst, image.Transparent, image.ZP, draw.Over)
		draw.Draw(drw, target, &image.Uniform{LIGHT_GRAY}, image.ZP, draw.Over)
		return drawDst
	}
	return highlight, lookup
}

func splitWds(lookup string) []string {
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

func Search(env gui.Env, fontFaces map[string]font.Face, words <-chan string, define chan<- string) {
	var list []string
	var display []imageObj
	for {
		select {
		case lookup := <-words:
			list = splitWds(lookup)
			display = displayWords(list, fontFaces["regular"])
			env.Draw() <- drawSearchBar(display, &wordCorner)
		case e, ok := <-env.Events():
			if !ok {
				close(env.Draw())
				return
			}
			switch e := e.(type) {
			case win.MoDown:
				if image.Pt(e.X, e.Y).In(wordCorner) {
					env.Draw() <- drawSearchBar(display, &wordCorner)
					highlight, target := highlightWord(display, image.Pt(e.X, e.Y), wordCorner, define)
					define <- target
					env.Draw() <- highlight
				}
				// case win.MoUp:
			}
		}
	}
}
