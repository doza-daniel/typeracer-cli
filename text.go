package main

import (
	"fmt"
	"strings"
)

func newText(orig string) *text {
	return &text{words: strings.Split(orig, " ")}
}

type text struct {
	currentWord int
	words       []string
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
