package text

import "github.com/qlova/uct/compiler"

var Atoi = compiler.Function {
	Name: compiler.NoTranslation("text_atoi"),
		
	Compile: func (c *compiler.Compiler) {
		c.Pull("base")
		
		//Result.
		c.Int(0)
		c.Pull("number")

		c.Int(1)
		c.Pull("exponent")
		
		c.Size()
		c.Int(1)
		c.Sub()
		c.Pull("i")
		
		c.Loop()
			c.Push("i")
			c.Int(0)
			c.More()
			c.If()
				c.Done()
			c.No()
			
			c.Push("i")
			c.Get()
			c.Pull("v")
			
			
			c.Int(57)
			c.Push("v")
			c.More()

			c.Int(46)
			c.Push("v")
			c.Less()
			
			c.Add()
			c.If()
				c.Int(0)
				c.Int(1)
				c.Name("ERROR")
				c.Back()
			c.No()

			c.Push("v")
			c.Int(48)
			c.Sub()
			c.Push("exponent")
			c.Mul()
			c.Push("number")
			c.Add()
			c.Name("number")
			
			c.Push("exponent")
			c.Push("base")
			c.Mul()
			c.Name("exponent")
			
			c.Push("i")
			c.Int(1)
			c.Sub()
			c.Name("i")
		c.Redo()
		
		c.Used()
		c.Push("number")
	},
}  
