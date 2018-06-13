package graphics

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/grate"
import "github.com/qlova/ilang/syntax/update"

import "github.com/qlova/ilang/syntax/global"

import "github.com/qlova/ilang/types/text"

var Name = compiler.Translatable {
		compiler.English: "graphics",
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.Back()
		c.Main()
		
		global.Init(c)

		c.Native("py", "global_runtime = runtime")
		c.Native("py", "pyglet.clock.schedule_interval(update, 1/30.0)")
		c.Native("py", "pyglet.app.run()")
		
		
		
		// Call ebiten.Run to start your game loop.
		c.Native("go", `global_runtime = runtime; ebiten.Run(grate_update, 800, 600, 1, "")`)
		
		c.Exit()
	},
}

var Graphics = compiler.Statement {
	Name: Name,
	
	OnScan:  func(c *compiler.Compiler) {
		if !c.GlobalFlagExists(grate.Flag) {
			grate.Init(c)
			c.SetGlobalFlag(grate.Flag)
		}
		
		if !c.GlobalFlagExists(update.Flag) {
			update.Init(c)
			c.SetGlobalFlag(update.Flag)
		}
		
		c.Native("go", `
var global_runtime *Runtime
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

func display_text() {
	var w, h = grate_screen.Size()
	var message = string(global_runtime.PullList().Bytes)
	
	//Measure.
	var length = grate_font_drawer.MeasureString(message)
	
	text.Draw(grate_screen, message, grate_font, w/2 - length.Round()/2, h/2, color.White)
}
				
`)
		
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
		
		c.Native("py", "grate_label.text = bytearray(runtime.Lists.pop()).decode('utf8')")
		c.Native("py", "grate_label.draw()")
		
		c.Native("go", "display_text()")
	},
}
