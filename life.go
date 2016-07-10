package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"time"
)

const (
	boardCount = 2
)

type Life struct {
	boards     [boardCount]Board
	board      *Board
	name       string
	generation int
	close      chan bool
}

func NewLife(name string, w int, h int) Life {
	l := Life{name: name, close: make(chan bool)}
	for i := 0; i < len(l.boards); i++ {
		l.boards[i] = NewBoard(w, h)
	}

	l.board = l.CurrentBoard()
	l.board.Random()

	return l
}

func (l *Life) CurrentBoard() *Board {
	return &l.boards[l.generation%len(l.boards)]
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
				l.Step(sizeX, sizeY)
				v.Title = fmt.Sprintf("%s - Generation %d", l.name, l.generation)
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
	l.generation++
	nb := l.CurrentBoard()
	cb := l.board
	for y := range nb.board {
		for x := range nb.board[y] {
			cb.NextCell(&nb.board[y][x], x, y, w, h)
		}
	}

	l.board = nb
}
