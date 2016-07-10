package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"time"
)

type Life struct {
	board      Board
	name       string
	generation int
	close      chan bool
}

func NewLife(name string, w int, h int) Life {
	board := NewBoard(w, h)
	board.Random()

	return Life{name: name,
		board: board,
		close: make(chan bool)}
}

// print the most recent board
func (l *Life) String() string {
	return l.board.String()
}

func (l *Life) start(g *gocui.Gui) error {
	for {
		select {
		case <-l.close:
			return nil
		case <-time.After(interval):
			g.Execute(func(g *gocui.Gui) error {
				v, err := g.View(l.name)
				if err != nil {
					return err
				}

				sizeX, sizeY := v.Size()
				l.generation++
				v.Title = fmt.Sprintf("%s - Generation %d", l.name, l.generation)
				l.Step(sizeX, sizeY)
				v.Clear()
				fmt.Fprint(v, l.String())
				return nil
			})
		}
	}
	return nil
}

// make the next iteration of the board
func (l *Life) Step(w, h int) {
	nb := NewBoard(w, h)
	cb := l.board
	for y := range nb.board {
		for x := range nb.board[y] {
			cb.NextCell(&nb.board[y][x], x, y, w, h)
		}
	}

	l.board = nb
}
