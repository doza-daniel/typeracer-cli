package game

import (
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/rivo/tview"
)

func newText(raw string) *text {
	return &text{
		raw:   raw,
		words: strings.Split(raw, " "),
		once:  &sync.Once{},
	}
}

type text struct {
	raw         string
	words       []string
	currentWord int

	once  *sync.Once
	start time.Time
}

func (t *text) acceptFun(typingField *tview.InputField) func(string, rune) bool {
	return func(current string, input rune) bool {
		t.once.Do(func() {
			t.start = time.Now()
		})

		if !unicode.IsGraphic(input) {
			return false
		}

		switch {
		case t.isLastWord() && current == t.current():
			fallthrough
		case unicode.IsSpace(input) && current[:len(current)-1] == t.current():
			typingField.SetText("")
			t.next()
			return false
		default:
			return true
		}
	}
}

func (t *text) current() string {
	return t.words[t.currentWord]
}

func (t *text) next() {
	t.currentWord++
}

func (t *text) isLastWord() bool {
	return t.currentWord == len(t.words)-1
}

func (t *text) isDone() bool {
	return t.currentWord == len(t.words)
}

func (t *text) calculateWPM(end time.Time) float64 {
	duration := end.Sub(t.start)
	return float64(len(t.raw)) / (5.0 * duration.Minutes())
}

func (t *text) color(input string) string {
	var prefix string = strings.Join(t.words[:t.currentWord], " ")
	var suffix string = strings.Join(t.words[t.currentWord:], " ")

	if len(prefix) > 0 && len(suffix) > 0 {
		prefix += " "
	}

	greenLen := matchingSubstringLen(input, suffix)
	redEnd := min(len(input), len(suffix))

	return fmt.Sprintf(
		"[green]%s%s[red]%s[white]%s",
		prefix,
		suffix[:greenLen],
		suffix[greenLen:redEnd],
		suffix[redEnd:],
	)
}

func matchingSubstringLen(s1, s2 string) int {
	var i int
	for i = 0; i < len(s1) && i < len(s2) && s1[i] == s2[i]; i++ {
	}
	return i
}

func min(a, b int, nums ...int) int {
	nums = append(nums, a)

	min := b
	for _, x := range nums {
		if x < min {
			min = x
		}
	}

	return min
}
