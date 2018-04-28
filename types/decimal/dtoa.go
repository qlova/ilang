package decimal

import "github.com/qlova/ilang/types/text"
import "github.com/qlova/uct/compiler"

var Dtoa = compiler.Function {
	Name: compiler.NoTranslation("text_dtoa"),
	
	//Ported from old_version.
	Compile: func (c *compiler.Compiler) {
		c.Pull("precision")
		c.Pull("exponent")
		c.Pull("base")
		c.Pull("value")
		
		c.Int(0)
		c.Pull("negative")
		
		c.Push("value")
		c.Int(0)
		c.More()
		c.If()
			c.Push("value")
			c.Flip()
			c.Name("value")
			
			c.Int(1)
			c.Name("negative")
		c.No()
		
		c.Push("value")
		c.Push("exponent")
		c.Div()
		
		c.Push("negative")
		c.If()
			c.Flip()
		c.No()
		
		c.Push("base")
		c.Call(&text.Itoa)
		c.Int('.')
		c.Put()
		c.PullList("integer")
		
		c.Push("value")
		c.Push("exponent")
		c.Mod()
		c.Push("base")
		c.Call(&text.Itoa)
		c.Size() //Size is on the stack.
		c.PullList("fractional")
		
		
		c.PushList("integer") //Integer is on the stack
		c.List()
		c.Int(0)
		c.Pull("i")
		
		c.Flip()
		c.Push("precision")
		c.Add()
		c.Pull("amount")
		
		c.Loop()
		
			c.Push("i")
			c.Push("amount")
			c.Same()
			c.If()
				c.Done()
			c.No()
			
			c.Int('0')
			c.Put()
			
			c.Push("i")
			c.Int(1)
			c.Add()
			c.Name("i")
		c.Redo()
		
		c.Call(&text.Join)
		c.PushList("fractional")
		c.Call(&text.Join)
	},
}  
