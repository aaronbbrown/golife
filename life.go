package main

import (
	"errors"
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"math/rand"
	"strings"
	"time"
)

const (
	delta    = 1
	interval = 150 * time.Millisecond
	icon     = '✼'
)

const (
	_ = iota
	Red
	Blue
	Green
	_ = iota
	Star
	Hash
	Circle
)

type Board struct {
	board [][]Cell
	w, h  int
}

type Life struct {
	board      Board
	name       string
	generation int
	close      chan bool
}

type Cell struct {
	alive bool
	color int
	shape int
}

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

func (c *Cell) Random() {
	colors := []int{Red, Green, Blue}
	shapes := []int{Star, Hash, Circle}

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
	}

	return ' '
}

/*
func (c *Cell) Next(neighbors []Cell) *Cell {

}
*/

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
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", 'n', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return newView(g)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", 'w', gocui.ModNone, closeView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if err := nextView(g); err != nil {
				return err
			}
			return ontop(g, v)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, -delta, 0)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'h', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, -delta, 0)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, delta, 0)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'l', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, delta, 0)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, 0, delta)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'j', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, 0, delta)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, 0, -delta)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'k', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return moveView(g, v, 0, -delta)
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 't', gocui.ModNone, ontop); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '+', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return resizeView(g, v, delta, 0)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '_', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return resizeView(g, v, -delta, 0)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '=', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return resizeView(g, v, 0, delta)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", '-', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return resizeView(g, v, 0, -delta)
		}); err != nil {
		return err
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	return nil
}

func layout(g *gocui.Gui) error {
	maxX, _ := g.Size()
	v, err := g.SetView("legend", maxX-25, 0, maxX-1, 8)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "KEYBINDINGS")
		fmt.Fprintln(v, "n: New View")
		fmt.Fprintln(v, "Tab: Next View")
		fmt.Fprintln(v, "← ↑ → ↓: Move View")
		fmt.Fprintln(v, "w: Delete View")
		fmt.Fprintln(v, "^C or q: Exit")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

// weighted coin toss.  odds is the change
// that the toss will appear true
func weightedRandBool(odds int) bool {
	return rand.Int()%odds != 0
}

func NewBoard(w, h int) Board {
	board := make([][]Cell, h)
	for i := range board {
		board[i] = make([]Cell, w)
	}
	return Board{w: w, h: h, board: board}
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
			err := cb.NextCell(&nb.board[y][x], x, y, w, h)
			if err != nil {
				log.Panicf("Couldn't get Cell at %d %d.  This is unexpected", x, y)
			}
		}
	}

	l.board = nb
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

	cell.Copy(c)

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

func (b *Board) CellAt(x, y int) (Cell, error) {
	// protect from out of bounds errors
	if y >= len(b.board) || x >= len(b.board[y]) {
		var c Cell
		return c, errors.New("out of bounds")
	}
	return b.board[y][x], nil
}

// a sane modulo operator that works like every other damn language
// see https://github.com/golang/go/issues/448
// https://groups.google.com/forum/#!topic/golang-nuts/xj7CV857vAg
func saneModInt(x, y int) int {
	result := x % y
	if result < 0 {
		result += y
	}
	return result
}

// returns the number living neighbors a cell has
// w & h are the width and height of the NEXT board
func (b *Board) Neighbors(x, y, w, h int) []Cell {
	neighbors := make([]Cell, 0)
	lpos := saneModInt((x - 1), w) // cell to the left
	rpos := saneModInt((x + 1), w) // cell to the right
	apos := saneModInt((y - 1), h) // cell above
	bpos := saneModInt((y + 1), h) // cell below
	// fmt.Printf("x: %d, y: %d, b.w: %d, b.h: %d, %d %d %d %d", x, y, b.w, b.h, lpos, rpos, apos, bpos)
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
	icons := make([][]rune, b.h)
	lines := make([]string, b.h)

	for y := range b.board {
		icons[y] = make([]rune, b.w)
		for x := range b.board[y] {
			c, _ := b.CellAt(x, y)
			icons[y][x] = c.Rune()
		}
		lines[y] = string(icons[y])
	}

	return strings.Join(lines, "\n")
}

func newView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := fmt.Sprintf("v%v", idxGame)
	v, err := g.SetView(name, 0, 0, maxX/5*4, maxY/5*4)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if err := g.SetCurrentView(name); err != nil {
		return err
	}

	w, h := v.Size()

	v.Title = name
	game := NewLife(name, w, h)
	games = append(games, game)
	go game.start(g)
	curGame = len(games) - 1
	idxGame += 1
	return nil
}

func moveView(g *gocui.Gui, v *gocui.View, dx, dy int) error {
	name := v.Name()
	x0, y0, x1, y1, err := g.ViewPosition(name)
	if err != nil {
		return err
	}
	if _, err := g.SetView(name, x0+dx, y0+dy, x1+dx, y1+dy); err != nil {
		return err
	}
	return nil
}

func nextView(g *gocui.Gui) error {
	next := curGame + 1
	if next > len(games)-1 {
		next = 0
	}

	if err := g.SetCurrentView(games[next].name); err != nil {
		return err
	}
	curGame = next
	return nil
}

func closeView(g *gocui.Gui, v *gocui.View) error {
	games[curGame].close <- true
	if err := g.DeleteView(games[curGame].name); err != nil {
		return err
	}
	// delete the game
	games = append(games[:curGame], games[curGame+1:]...)
	return nextView(g)
}

func ontop(g *gocui.Gui, v *gocui.View) error {
	_, err := g.SetViewOnTop(games[curGame].name)
	return err
}

func resizeView(g *gocui.Gui, v *gocui.View, xdelta, ydelta int) error {
	x0, y0, x1, y1, err := g.ViewPosition(games[curGame].name)

	x1 += xdelta
	y1 += ydelta
	_, err = g.SetView(games[curGame].name, x0, y0, x1, y1)
	if err != nil {
		return err
	}

	return nil
}

// go doesn't have a Int.max function
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
