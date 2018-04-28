package number

import "github.com/qlova/uct/compiler"

var fact = compiler.Function {
	Name: compiler.NoTranslation("number_fact"),
} 

func init() {
	fact.Compile = func (c *compiler.Compiler) {
		c.Pull("acc")
		c.Pull("n")
		
		c.Push("n")
		c.Int(0)
		c.Same()
		c.If()
			c.Push("acc")
			c.Back()
		c.No()
		
		c.Push("n")
		c.Int(1)
		c.Sub()
		c.Push("n")
		c.Push("acc")
		c.Mul()
		
		c.Call(&fact)
	}
}

var Factorial = compiler.Function {
	Name: compiler.NoTranslation("number_factorial"),
	
	Compile: func (c *compiler.Compiler) {
		c.Pull("n")
		
		
		
		c.Push("n")
		c.Int(0)
		c.More()
		c.If()
			c.Int(0)
			c.Back()
		c.No()
		
		c.Push("n")
		c.Int(1)
		c.Call(&fact)
	},
} 

