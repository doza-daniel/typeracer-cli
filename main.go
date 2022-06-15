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

func (*mockCorpus) GetTextAt(int64) string {
	return "The quick brown fox jumps over the lazy dog."
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
