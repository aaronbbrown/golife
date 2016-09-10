package main

import (
	"bytes"
	"errors"
)

type Board struct {
	board [][]Cell
	w, h  int
}

func NewBoard(w, h int) Board {
	board := make([][]Cell, h)
	for i := range board {
		board[i] = make([]Cell, w)
	}
	return Board{w: w, h: h, board: board}
}

// return whether a cell is currently alive
func (b *Board) Alive(x, y int) bool {
	// protect from out of bounds errors
	if y >= len(b.board) || x >= len(b.board[y]) {
		return false
	}
	return b.board[y][x].alive
}

// return whether a cell will be alive on the next iteration
// w & h are the width and height of the new board
func (b *Board) NextCell(cell *Cell, x, y, w, h int) error {
	neighbors := b.Neighbors(x, y, w, h)
	count := len(neighbors)

	// currently alive cell
	c, err := b.CellAt(x, y)
	if err != nil {
		return err
	}

	*cell = *c

	if c.alive {
		cell.alive = count >= 2 && count <= 3
	} else {
		// reproduce
		if count == 3 || count == 6 {
			cell.alive = true
			cell.SetNextShape(neighbors)
		}
	}
	return nil
}

func (b *Board) CellAt(x, y int) (*Cell, error) {
	// protect from out of bounds errors
	if y >= len(b.board) || x >= len(b.board[y]) {
		var c Cell
		return &c, errors.New("out of bounds")
	}
	return &b.board[y][x], nil
}

// returns the number living neighbors a cell has
// w & h are the width and height of the NEXT board
func (b *Board) Neighbors(x, y, w, h int) []*Cell {
	neighbors := make([]*Cell, 0, 8)
	lpos := saneModInt((x - 1), w) // cell to the left
	rpos := saneModInt((x + 1), w) // cell to the right
	apos := saneModInt((y - 1), h) // cell above
	bpos := saneModInt((y + 1), h) // cell below

	// above left
	if b.Alive(lpos, apos) {
		c, err := b.CellAt(lpos, apos)
		if err == nil {
			neighbors = append(neighbors, c)
		}
	}
	// above
	if b.Alive(x, apos) {
		c, err := b.CellAt(x, apos)
		if err == nil {
			neighbors = append(neighbors, c)
		}
	}
	// above right
	if b.Alive(rpos, apos) {
		c, err := b.CellAt(rpos, apos)
		if err == nil {
			neighbors = append(neighbors, c)
		}
	}

	// left
	if b.Alive(lpos, y) {
		c, err := b.CellAt(lpos, y)
		if err == nil {
			neighbors = append(neighbors, c)
		}
	}
	// right
	if b.Alive(rpos, y) {
		c, err := b.CellAt(rpos, y)
		if err == nil {
			neighbors = append(neighbors, c)
		}
	}
	// below left
	if b.Alive(lpos, bpos) {
		c, err := b.CellAt(lpos, bpos)
		if err == nil {
			neighbors = append(neighbors, c)
		}
	}
	// below
	if b.Alive(x, bpos) {
		c, err := b.CellAt(x, bpos)
		if err == nil {
			neighbors = append(neighbors, c)
		}
	}
	// below right
	if b.Alive(rpos, bpos) {
		c, err := b.CellAt(rpos, bpos)
		if err == nil {
			neighbors = append(neighbors, c)
		}
	}

	return neighbors
}

func (b *Board) Random() {
	for y := range b.board {
		for x := range b.board[y] {
			b.board[y][x].Random()
		}
	}
}

func (b *Board) String() string {
	var buf bytes.Buffer

	for y := range b.board {
		for x := range b.board[y] {
			c, _ := b.CellAt(x, y)
			buf.WriteString(c.String())
		}
		buf.WriteByte('\n')
	}

	return buf.String()
}
