package update

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/grate"
import "github.com/qlova/uct/compiler"

var Name = compiler.Translatable {
	compiler.English: "update",
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.Native("py", "for key in keyboard_buffer:")
		c.Native("py", "	keyboard_buffer[key] = False")
		
		c.Back()
		//c.SetBool("update_exists", true)
	},
}

func Init(c *compiler.Compiler) {
	c.Native("go", `
// update is called every frame (1/60 [s]).
func update(screen *ebiten.Image) error {
	grate_screen = screen

    // Write your game's logical update.
    grate_graphics(global_runtime)

    return nil
}
`)
}

var Statement = compiler.Statement {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) {
		if len(c.Scope) != 0 {
			c.RaiseError(compiler.Translatable{
				compiler.English: "Update must be placed at the top level!",
			})
		}
		
		if !c.GlobalFlagExists(grate.Flag) {
			grate.Init(c)
			c.SetGlobalFlag(grate.Flag)
		}
		c.SwapOutput()
		c.Code("i_placeholder_i")
		c.Send()
		c.Back()
		c.SwapOutput()
		
		c.Native("go", `
// update is called every frame (1/60 [s]).
func grate_update(screen *ebiten.Image) error {
	grate_screen = screen
	
	//TODO add this dynamically.
	Last_Keyboard_Buffer = Keyboard_Buffer
	Update_Keyboard_Buffer()

	// Write your game's logical update.
	grate_graphics(global_runtime)
	
	update(global_runtime)
		
	Stdout.Flush()

	return nil
}
		`)
		
		c.Code("update")
		c.Native("py", "runtime = global_runtime")
		
		c.Expecting(symbols.CodeBlockBegin)
		c.GainScope()
		c.SetFlag(Flag)
		
		c.SetGlobalFlag(Flag)
	},
}
