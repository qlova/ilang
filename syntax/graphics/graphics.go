package graphics

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/software"

import "github.com/qlova/ilang/types/text"

var Name = compiler.Translatable {
		compiler.English: "graphics",
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.Back()
		c.Native("py", "pyglet.app.run()")
	},
}

var Graphics = compiler.Statement {
	Name: Name,
	
	OnScan:  func(c *compiler.Compiler) {
		if _, index := c.GetFlag(software.Flag); index == -1 {
			c.SwapOutput()
			
			c.Native("go", `
import "github.com/hajimehoshi/ebiten"
import "github.com/hajimehoshi/ebiten/text"
import "image/color"
import "encoding/base64"
import "github.com/golang/freetype/truetype"
import "golang.org/x/image/font"
import "fmt"`)

			c.SwapOutput()
			
						c.Native("py", `
import pyglet
import sys
import math
import queue
from pyglet.gl import *
from pyglet.window import key

window = pyglet.window.Window()
grate_label = pyglet.text.Label("", font_name="Times New Roman",font_size=16,x=0, y=0, anchor_y="center", anchor_x="center")

runtime = Runtime()
			
@window.event
def on_draw():
	glClearColor(0, 0, 0, 1)
	glClear(GL_COLOR_BUFFER_BIT)

	glLoadIdentity()
	glTranslatef(window.width/2, window.height/2, 0)

	glEnable(GL_BLEND)
	glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)
			
	grate_graphics(runtime)
`)

			c.Native("go", `
var runtime = new(Runtime)
var grate_screen *ebiten.Image

var grate_font font.Face
var grate_font_drawer font.Drawer

func init() {
	data, err := base64.StdEncoding.DecodeString(`+"`"+Font+"`"+`)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	tt, err := truetype.Parse(data)
	if err != nil {
		println(err.Error())
	}
	
	const dpi = 72
	grate_font = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	grate_font_drawer = font.Drawer{Face: grate_font}
}

// update is called every frame (1/60 [s]).
func update(screen *ebiten.Image) error {
	grate_screen = screen

    // Write your game's logical update.
    grate_graphics(runtime)

    return nil
}

func display_text() {
	var w, h = grate_screen.Size()
	var message = string(runtime.PullList().Bytes)
	
	//Measure.
	var length = grate_font_drawer.MeasureString(message)
	
	text.Draw(grate_screen, message, grate_font, w/2 - length.Round()/2, h/2, color.White)
}

func main() {
    // Call ebiten.Run to start your game loop.
    ebiten.Run(update, 800, 600, 1, "")
}
				
`)
		}
		
		c.Expecting(symbols.CodeBlockBegin)
		c.Code("grate_graphics")
		c.GainScope()
		c.SetFlag(Flag)
	},
}

var Display = compiler.Statement {
	Name: compiler.Translatable {
		compiler.English: "display",
	},
	
	OnScan: func(c *compiler.Compiler) {
		c.Expecting(symbols.FunctionCallBegin)
		var arg = c.ScanExpression()
		c.Expecting(symbols.FunctionCallEnd)
		
		if !arg.Equals(text.Type) {
			c.RaiseError(errors.ExpectingType(text.Type, arg))
		}
		
		c.Native("py", "grate_label.text = runtime.Lists.pop().decode('utf8')")
		c.Native("py", "grate_label.draw()")
		
		c.Native("go", "display_text()")
	},
}
