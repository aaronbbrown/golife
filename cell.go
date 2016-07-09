package main

import (
	"github.com/fatih/color"
	"math/rand"
)

type Cell struct {
	alive bool
	color color.Attribute
	shape rune
}

const (
	mutationRate = 10
)

var colors = [...]color.Attribute{
	color.FgMagenta,
	color.FgGreen,
	color.FgCyan,
	color.FgBlue,
	color.FgRed}

var shapes = [...]rune{'+', '*', '0', '#'}

func NewCell() *Cell {
	return &Cell{alive: false, color: colors[0], shape: shapes[0]}
}

func (c *Cell) Random() {
	c.alive = weightedRandBool(2)
	c.color = colors[rand.Intn(len(colors))]
	c.shape = shapes[rand.Intn(len(shapes))]
}

func (c *Cell) Copy(src Cell) {
	c.alive = src.alive
	c.color = src.color
	c.shape = src.shape
}

func (c *Cell) Rune() rune {
	if c.alive == false {
		return ' '
	}
	return c.shape
}

func (c *Cell) String() string {
	colorFunc := color.New(c.color, color.Bold).SprintFunc()
	return colorFunc(string(c.Rune()))
}

func (c *Cell) SetNextShape(neighbors []Cell) {
	shapeCount := make(map[rune]int)
	for _, n := range neighbors {
		shapeCount[n.shape]++
	}
	var mostCommonShape rune
	max := 0
	for shape, count := range shapeCount {
		if count > max {
			max = count
			mostCommonShape = shape
		}
	}
	if max == 1 {
		// pick a random one from the parents
		c.shape = neighbors[rand.Intn(len(neighbors))].shape
	} else {
		c.shape = mostCommonShape
	}

	// TODO DRY this up.
	colorCount := make(map[color.Attribute]int)
	for _, n := range neighbors {
		colorCount[n.color]++
	}
	var mostCommonColor color.Attribute
	max = 0
	for color, count := range colorCount {
		if count > max {
			max = count
			mostCommonColor = color
		}
	}
	if max == 1 {
		// pick a random one from the parents
		c.color = neighbors[rand.Intn(len(neighbors))].color
	} else {
		c.color = mostCommonColor
	}

	c.mutate()
}

// Throw in a random mutation every once in a while
func (c *Cell) mutate() bool {
	var result bool
	if rand.Intn(mutationRate) == 0 {
		c.color = colors[rand.Intn(len(colors))]
	}
	if rand.Intn(mutationRate) == 0 {
		c.shape = shapes[rand.Intn(len(shapes))]
	}
	return result
}
