package grate

import "github.com/qlova/uct/compiler"

var Flag = compiler.Flag{}

func Init(c *compiler.Compiler) {
	c.SwapOutput()
	
	c.Native("go", `
import "github.com/hajimehoshi/ebiten"
import "github.com/hajimehoshi/ebiten/text"
import "image/color"
import "encoding/base64"
import "github.com/golang/freetype/truetype"
import "golang.org/x/image/font"
import "fmt"`)
	
	c.Native("py", `
import pyglet
import sys
import math
import queue
from pyglet.gl import *
from pyglet.window import key

grate_window = pyglet.window.Window()
grate_label = pyglet.text.Label("", font_name="Times New Roman",font_size=16,x=0, y=0, anchor_y="center", anchor_x="center")
	
`)
	c.SwapOutput()

	c.Native("py", `	
global_runtime = Runtime()
			
@grate_window.event
def on_draw():
	glClearColor(0, 0, 0, 1)
	glClear(GL_COLOR_BUFFER_BIT)

	glLoadIdentity()
	glTranslatef(grate_window.width/2, grate_window.height/2, 0)

	glEnable(GL_BLEND)
	glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA)
			
	grate_graphics(global_runtime)
`)
}
