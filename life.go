package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
	"math/rand"
	"strings"
	"time"
)

const delta = 1
const interval = 150 * time.Millisecond

type Board struct {
	board [][]bool
	w, h  int
}

type Life struct {
	board      Board
	name       string
	generation int
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

	if err := g.SetKeybinding("", 'r', gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if err := nextView(g, true); err != nil {
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
		fmt.Fprintln(v, "n New View")
		//        fmt.Fprintln(v, "Tab: Next View")
		//        fmt.Fprintln(v, "← ↑ → ↓: Move View")
		//        fmt.Fprintln(v, "Backspace: Delete View")
		//        fmt.Fprintln(v, "t: Set view on top")
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
	board := make([][]bool, h)
	for i := range board {
		board[i] = make([]bool, w)
	}
	return Board{w: w, h: h, board: board}
}

func NewLife(name string, w int, h int) Life {
	board := NewBoard(w, h)
	board.Random()

	return Life{name: name, board: board}
}

// print the most recent board
func (l *Life) String() string {
	return l.board.String()
}

func (l *Life) start(g *gocui.Gui) error {
	for {
		select {
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
			nb.board[y][x] = cb.NextAlive(x, y, w, h)
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
	return b.board[y][x]
}

// return whether a cell will be alive on the next iteration
// w & h are the width and height of the new board
func (b *Board) NextAlive(x, y, w, h int) bool {
	neighbors := b.Neighbors(x, y, w, h)

	// currently alive cell
	if b.Alive(x, y) {
		return neighbors >= 2 && neighbors <= 3
	} else {
		// reproduce
		if neighbors == 3 {
			return true
		}
	}
	return b.Alive(x, y)
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
func (b *Board) Neighbors(x, y, w, h int) int {
	count := 0
	lpos := saneModInt((x - 1), w) // cell to the left
	rpos := saneModInt((x + 1), w) // cell to the right
	apos := saneModInt((y - 1), h) // cell above
	bpos := saneModInt((y + 1), h) // cell below
	// fmt.Printf("x: %d, y: %d, b.w: %d, b.h: %d, %d %d %d %d", x, y, b.w, b.h, lpos, rpos, apos, bpos)
	// above left
	if b.Alive(lpos, apos) {
		count += 1
	}
	// above
	if b.Alive(x, apos) {
		count += 1
	}
	// above right
	if b.Alive(rpos, apos) {
		count += 1
	}
	// left
	if b.Alive(lpos, y) {
		count += 1
	}
	// right
	if b.Alive(rpos, y) {
		count += 1
	}
	// below left
	if b.Alive(lpos, bpos) {
		count += 1
	}
	// below
	if b.Alive(x, bpos) {
		count += 1
	}
	// below right
	if b.Alive(rpos, bpos) {
		count += 1
	}

	return count
}

func (b *Board) Random() {
	for y := range b.board {
		for x := range b.board[y] {
			// ~ 1/3rd of spaces will be filled
			b.board[y][x] = weightedRandBool(3)
		}
	}
}

func (b *Board) String() string {
	icons := make([][]rune, b.h)
	lines := make([]string, b.h)

	for y := range b.board {
		icons[y] = make([]rune, b.w)
		for x := range b.board[y] {
			if b.board[y][x] {
				icons[y][x] = '*'
			} else {
				icons[y][x] = ' '
			}
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

func nextView(g *gocui.Gui, disableCurrent bool) error {
	next := curGame + 1
	if next > len(games)-1 {
		next = 0
	}

	nv, err := g.View(games[next].name)
	if err != nil {
		return err
	}
	if err := g.SetCurrentView(games[next].name); err != nil {
		return err
	}
	nv.BgColor = gocui.ColorRed

	if disableCurrent && len(games) > 1 {
		cv, err := g.View(games[curGame].name)
		if err != nil {
			return err
		}
		cv.BgColor = g.BgColor
	}

	curGame = next
	return nil
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
