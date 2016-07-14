[![Build Status](https://travis-ci.org/aaronbbrown/golife.svg?branch=master)](https://travis-ci.org/aaronbbrown/golife)

# Overview

[![asciicast](https://asciinema.org/a/79690.png)](https://asciinema.org/a/79690)

This is [Conway's Game of Life](https://en.wikipedia.org/wiki/Conway%27s_Game_of_Life) ([Highlife rules](https://en.wikipedia.org/wiki/Highlife_(cellular_automaton))) implemented as a console application in Go.
Additional rules have been added that simulate a genome with a gene for color and shape.

The complete rules are:
* Survive when there are 2 or 3 living neighbors
* Die with less than 2 or more than 3 neighbors
* Born when there are exactly 3 or 6 neighbors (highlife)
* When a new cell is born, it takes on the color of the majority its parents
  or if there is no majority, a random color will be chosen from the set of colors
  present in the parents.
* When a new cell is born, it takes on the shape of the majority of its parents
  or if there is no majority, a random shape will be chosen from the set of shapes
  present in the parents.
* With each birth, there is a 10% chance it will take on a random color (mutation)
* With each birth, there is a 10% chance it will take on a random shape (mutation)

# Running

```
go build && ./life
```

Key mappings

| Key | Action |
| --- | ------ |
| n   | create a new board |
| w   | close a board |
| q / Ctrl-C   | quit the game |
| up, down, left, right | Move the window around |
| - | shrink window  vertically |
| = | grow window vertically |
| _ | shrink window horizontally |
| + | grow window horizontally |
| <tab> | cycle through windows |


# Testing

```
go test --bench=.
```
