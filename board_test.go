package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func BenchmarkNewBoard(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewBoard(200, 200)
	}
}

func BenchmarkNeighbors(b *testing.B) {
	board := NewBoard(200, 200)
	board.Random()
	for i := 0; i < b.N; i++ {
		board.Neighbors(100, 100, 100, 100)
	}
}

func BenchmarkNextCell(b *testing.B) {
	board := NewBoard(200, 200)
	board.Random()
	cell, err := board.CellAt(100, 100)
	for i := 0; i < b.N; i++ {
		board.NextCell(cell, 100, 100, 200, 200)
	}
	assert.Equal(b, err, nil)
}

func BenchmarkCellAt(b *testing.B) {
	board := NewBoard(200, 200)
	board.Random()

	for i := 0; i < b.N; i++ {
		board.CellAt(100, 100)
	}
}
