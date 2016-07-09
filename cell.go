package main

import (
	"github.com/fatih/color"
	"math/rand"
)

type Cell struct {
	alive bool
	color color.Attribute
	shape int
}

var colors = [...]color.Attribute{
	color.FgMagenta,
	color.FgGreen,
	color.FgCyan,
	color.FgBlue,
	color.FgRed}

const (
	_ = iota
	Star
	Hash
	Circle
	Plus
)

func (c *Cell) Random() {
	shapes := []int{Star, Hash, Circle, Plus}

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

	switch c.shape {
	case Hash:
		return '#'
	case Circle:
		return '0'
	case Star:
		return '*'
	case Plus:
		return '+'
	}

	return ' '
}

func (c *Cell) String() string {
	colorFunc := color.New(c.color, color.Bold).SprintFunc()
	return colorFunc(string(c.Rune()))
}

func (c *Cell) SetNextShape(neighbors []Cell) {
	shapeCount := make(map[int]int)
	for _, n := range neighbors {
		shapeCount[n.shape]++
	}
	mostCommonShape := 0
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

}
