package main

import (
	"fmt"
	"net"

	"github.com/jroimartin/gocui"
	// Import the server package
)

type TUIClient struct {
	conn net.Conn
	gui  *gocui.Gui
}

func NewTUIClient(conn net.Conn) (*TUIClient, error) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, err
	}

	return &TUIClient{
		conn: conn,
		gui:  g,
	}, nil
}

func (c *TUIClient) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Chat view
	if v, err := g.SetView("chat", 0, 0, maxX-1, maxY-3); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Chat"
		v.Wrap = true
		v.Autoscroll = true
	}

	// Input view
	if v, err := g.SetView("input", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Input"
		v.Editable = true
		v.Wrap = true
		if _, err := g.SetCurrentView("input"); err != nil {
			return err
		}
	}

	return nil
}

func (c *TUIClient) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *TUIClient) sendMessage(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	msg := v.Buffer()
	if len(msg) == 0 {
		return nil
	}

	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		return err
	}

	v.Clear()
	v.SetCursor(0, 0)
	return nil
}

func (c *TUIClient) receiveMessages() {
	const (
		ColorBrightRed = "\033[91m"
		ColorReset     = "\033[0m"
	)

	buffer := make([]byte, 1024)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			c.gui.Update(func(g *gocui.Gui) error {
				v, err := g.View("chat")
				if err != nil {
					return err
				}
				fmt.Fprintln(v, ColorBrightRed+"Disconnected from server"+ColorReset)
				return nil
			})
			return
		}

		message := string(buffer[:n])
		c.gui.Update(func(g *gocui.Gui) error {
			v, err := g.View("chat")
			if err != nil {
				return err
			}
			fmt.Fprint(v, message)
			return nil
		})
	}
}

func (c *TUIClient) Run() error {
	defer c.gui.Close()
	defer c.conn.Close()

	c.gui.SetManagerFunc(c.layout)

	if err := c.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit); err != nil {
		return err
	}

	if err := c.gui.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, c.sendMessage); err != nil {
		return err
	}

	go c.receiveMessages()

	if err := c.gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
