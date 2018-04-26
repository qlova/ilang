package text

import "github.com/qlova/uct/compiler"

var Itoa = compiler.Function {
	Name: compiler.NoTranslation("text_itoa"),
		
	Compile: func (c *compiler.Compiler) {
		c.Pull("base")
		c.Pull("number")
		c.List()
		
		c.Int(0)
		c.Push("number")
		c.Same()
		c.If()
			c.Int('0')
			c.Put()
			c.Back()
		c.No()
		
		c.Int(1)
		c.Pull("exponent")
		
		c.Int(0)
		c.Push("number")
		c.Less()
		c.If()
			c.Int('-')
			c.Put()
			c.Push("number")
			c.Flip()
			c.Name("number")
		c.No()
		
		//What is the highest power to 10 which fits in num.
		c.Loop()
			c.Push("number")
			c.Push("exponent")
			c.More()
			c.If()
				c.Push("base")
				c.Push("exponent")
				
				c.Div()
				c.Name("exponent")
				c.Done()
			c.No()
			
			c.Push("base")
			c.Push("exponent")
			c.Mul()
			c.Name("exponent")
		c.Redo()
		
		//Find each digit.
		c.Loop()
		
			//if exponent <= 0
			c.Int(0)
			c.Push("number")
			c.Less()
			c.Int(0)
			c.Push("number")
			c.Same()
			c.Add()
			
			c.If()
				c.Done()
			c.No()
			
			c.Push("exponent")
			c.Push("number")
			c.Div()
			c.Copy()
			c.Push("exponent")
			c.Mul()
			c.Push("number")
			c.Sub()
			c.Name("number")
			
			c.Int('0')
			c.Add()
			c.Put()
			
			c.Push("exponent")
			c.Push("base")
			c.Div()
			c.Name("exponent")
		c.Redo()
	},
}  
