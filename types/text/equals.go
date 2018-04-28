package text

import "github.com/qlova/uct/compiler"

var Equals = compiler.Function {
	Name: compiler.NoTranslation("text_equals"),
	
	Compile: func (c *compiler.Compiler) {
		c.PullList("b")
		c.PullList("a")
		
		//Exit early.
		c.PushList("b")
		c.Size()
		c.DropList()
		
		c.PushList("a")
		c.Size()
		c.DropList()
		
		c.Same()
		c.Int(0)
		c.Div()
		
		c.If()
			c.Int(0)
			c.Back()
		c.No()
		
		//Test every member.
		c.Int(0)
		c.Pull("i")
		c.Loop()
			c.PushList("a")
			c.Size()
			c.DropList()
			c.Push("i")
			c.Same()
			c.If()
				c.Int(1)
				c.Back()
			c.No()
			
			c.PushList("a")
			c.Push("i")
			c.Get()
			c.DropList()
			
			c.PushList("b")
			c.Push("i")
			c.Get()
			c.DropList()
			
			c.Same()
			c.Int(0)
			c.Div()
			c.If()
				c.Int(0)
				c.Back()
			c.No()
			
			c.Push("i")
			c.Int(1)
			c.Add()
			c.Name("i")
		c.Redo()
	},
} 
