package text

import "github.com/qlova/uct/compiler"

var Join = compiler.Function {
	Name: compiler.NoTranslation("text_join"),
	
	Compile: func (c *compiler.Compiler) {
		c.PullList("b")
		c.PullList("a")
		
		//Create a list to store the result.
		c.List()
		c.PullList("r")
		
		//Add contents of a.
		c.Int(0)
		c.Pull("i")
		c.Loop()
			c.PushList("a")
			c.Size()
			c.Push("i")
			c.Same()
			c.If()
				c.DropList()
				c.Done()
			c.No()
			
			c.Push("i")
			c.Get()
			c.DropList()
			
			c.PushList("r")
			c.Put()
			c.NameList("r")
			
			c.Push("i")
			c.Int(1)
			c.Add()
			c.Name("i")
		c.Redo()
		
		//Add contents of b.
		c.Int(0)
		c.Name("i")
		c.Loop()
			c.PushList("b")
			c.Size()
			c.Push("i")
			c.Same()
			c.If()
				c.DropList()
				c.Done()
			c.No()
			
			c.Push("i")
			c.Get()
			c.DropList()
			
			c.PushList("r")
			c.Put()
			c.NameList("r")
			
			c.Push("i")
			c.Int(1)
			c.Add()
			c.Name("i")
		c.Redo()
		
		c.PushList("r")
	},
} 
