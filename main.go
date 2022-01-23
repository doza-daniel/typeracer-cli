package main

import (
	"fmt"
	"play/db"
	"sync"
	"time"
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	texts, err := db.New("/home/daniel/Documents/texts.db")
	if err != nil {
		panic(err)
	}

	raw, err := texts.GetRandomText()
	if err != nil {
		panic(err)
	}

	txt := newText(raw)

	textView := tview.NewTextView()
	textView.SetDynamicColors(true)
	textView.SetText(txt.color(""))

	textFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	textFlex.SetBorder(true).SetTitle("TypeRacerCLI")
	textFlex.AddItem(textView, 0, 1, false)

	typingField := tview.NewInputField()
	typingField.SetFieldBackgroundColor(tcell.ColorBlack)

	typingFlex := tview.NewFlex().SetDirection(tview.FlexRow)
	typingFlex.SetBorder(true)
	typingFlex.AddItem(typingField, 0, 1, true)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(textFlex, 0, 1, false)
	flex.AddItem(typingFlex, 3, 1, true)

	app := tview.NewApplication()
	app.SetRoot(flex, true)

	once := &sync.Once{}
	var start time.Time

	typingField.SetAcceptanceFunc(func(s string, r rune) bool {
		once.Do(func() {
			start = time.Now()
		})
		if !unicode.IsGraphic(r) {
			return false
		}

		switch {
		case txt.isLastWord() && s == txt.current():
			fallthrough
		case unicode.IsSpace(r) && s[:len(s)-1] == txt.current():
			typingField.SetText("")
			txt.next()
			return false
		default:
			return true
		}
	})

	typingField.SetChangedFunc(func(s string) {
		if txt.currentWord == len(txt.words) {
			duration := time.Now().Sub(start)
			wpm := float64(len(raw)) / (5.0 * duration.Minutes())
			textView.SetText(fmt.Sprintf("Congratulations! Your typing speed is %.2f WPM", wpm))
			app.SetRoot(textFlex, true)
		} else {
			textView.SetText(txt.color(s))
		}
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}
