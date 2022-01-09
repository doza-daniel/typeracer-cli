package main

import (
	"unicode"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

func main() {
	txt := newText("Hello World!")

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

	typingField.SetAcceptanceFunc(func(s string, r rune) bool {
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
			textView.SetText("Congratulations!")
			app.SetRoot(textFlex, true)
		} else {
			textView.SetText(txt.color(s))
		}
	})

	if err := app.Run(); err != nil {
		panic(err)
	}
}