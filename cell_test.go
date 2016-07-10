package main

import (
	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewCell(t *testing.T) {
	cell := NewCell()
	assert.Equal(t, cell.alive, false, "Cell should be dead")
	assert.Equal(t, cell.color, color.FgMagenta, "Cell should be magenta")
	assert.Equal(t, cell.shape, '+', "Cell should be +")
}

func TestCellString(t *testing.T) {
	cell := NewCell()
	cell.alive = true
	assert.NotEqual(t, cell.String(), " ", "Cell should not be blank if alive")

	cell.alive = false
	assert.Equal(t, cell.String(), " ", "Dead Cell should be ' '")
}
