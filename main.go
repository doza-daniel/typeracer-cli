package main

import (
	"flag"
	"log"
	"play/db"
	"play/game"
)

type mockCorpus struct {
	txt string
}

func (*mockCorpus) GetTextAt(int64) game.Text {
	return game.Text{
		Content: "The quick brown fox jumps over the lazy dog.",
		Source:  "Mock_Source",
		Type:    "Mock_Type",
		Author:  "Mock_Author",
	}
}
func (*mockCorpus) Size() int64 {
	return 1
}

func main() {
	dbPath := flag.String("db", "", "")
	flag.Parse()

	var corpus game.Corpus
	if dbPath == nil || *dbPath == "" {
		corpus = &mockCorpus{}
	} else {
		var err error
		corpus, err = db.New(*dbPath)
		if err != nil {
			log.Fatal("could not open the database: %v", err)
		}
	}

	g := game.New(corpus)
	log.Fatal(g.Run())
}
