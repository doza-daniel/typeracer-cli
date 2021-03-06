package game

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
)

// Text ...
type Text struct {
	Content string
	Type    string
	Author  string
	Source  string
}

// Corpus ...
type Corpus interface {
	GetTextAt(int64) Text
	Size() int64
}

// Game ...
type Game struct {
	app    *tview.Application
	corpus Corpus

	txtDetails Text

	mainContainer *tview.Flex
}

// New ...
func New(corpus Corpus) *Game {
	return &Game{
		corpus: corpus,
		app:    tview.NewApplication(),
	}
}

// Run ...
func (game *Game) Run() error {
	game.app.SetInputCapture(func(e *tcell.EventKey) *tcell.EventKey {
		switch e.Key() {
		case tcell.KeyCtrlR:
			game.start()
			return nil
		default:
			return e
		}
	})

	game.start()
	return game.app.Run()
}

func (game *Game) start() {
	at := rand.Int63n(game.corpus.Size())
	text := game.corpus.GetTextAt(at)
	game.txtDetails = text
	txt := newText(text.Content)

	textView := game.textView(txt)
	typingField := game.typingInputField(txt, textView)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(textView, 0, 1, false)
	flex.AddItem(typingField, 3, 0, true)

	flex.SetBorder(true).SetTitle("Typing Test")

	game.mainContainer = flex

	game.app.SetRoot(game.mainContainer, true)
}

func (game *Game) textView(txt *text) *tview.TextView {
	textView := tview.NewTextView()
	textView.SetWordWrap(true)
	textView.SetDynamicColors(true)
	textView.SetText(txt.color(""))

	return textView
}

func (game *Game) typingInputField(txt *text, textView *tview.TextView) *tview.InputField {
	typingField := tview.NewInputField()

	typingField.SetAcceptanceFunc(txt.acceptFun(typingField))

	typingField.SetChangedFunc(func(s string) {
		if txt.isDone() {
			wpm := txt.calculateWPM(time.Now())
			game.resultPage(wpm)
		} else {
			textView.SetText(txt.color(s))
		}
	})

	typingField.SetBorder(true)
	typingField.SetTitle("type here")
	typingField.SetFieldBackgroundColor(tcell.ColorBlack)

	return typingField
}

func (game *Game) resultPage(wpm float64) {
	details := fmt.Sprintf(
		"Type: %s\nSource: %s\nAuthor: %s",
		game.txtDetails.Type,
		game.txtDetails.Source,
		game.txtDetails.Author,
	)

	txt := fmt.Sprintf(
		"\n\nCongratulations! Your typing speed is %.2f WPM\n\n\n%s\n\n\nPress ENTER to play again.",
		wpm,
		details,
	)

	textView := tview.NewTextView()
	textView.SetTextAlign(tview.AlignCenter)
	textView.SetText(txt)
	textView.SetDoneFunc(func(tcell.Key) {
		game.start()
	})

	game.mainContainer.Clear()
	game.mainContainer.AddItem(textView, 0, 1, true)
	game.app.SetRoot(game.mainContainer, true)
}
