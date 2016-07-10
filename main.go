package main

import (
	"github.com/jroimartin/gocui"
	"log"
	"time"
)

const (
	delta    = 1
	interval = 150 * time.Millisecond
)

var (
	games   []Life
	curGame = -1
	idxGame = 0
)

func main() {
	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetLayout(layout)
	newView(g)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}
}
