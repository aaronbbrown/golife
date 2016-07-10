package main

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"log"
)

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

func newView(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	name := fmt.Sprintf("v%v", idxGame)
	v, err := g.SetView(name, 0, 0, maxX/10*9, maxY/10*9)
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
