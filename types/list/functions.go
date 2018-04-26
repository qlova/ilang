package list

import "github.com/qlova/uct/compiler"

var Copy = compiler.Function {
	Name: compiler.NoTranslation("list_copy"),
		
	Compile: func (c *compiler.Compiler) {
		c.Size()
		c.PullList("list")
		c.Make()
		c.PullList("copy")
		
		c.Int(0)
		c.Pull("i")
		c.Loop()
			c.PushList("list")
			c.Push("i")
			c.Size()
			c.Int(1)
			c.Sub()
			c.Same()
			c.If()
				c.Done()
			c.No()

			c.Push("i")
			c.Copy()
			c.Get()
			c.Used()
			
			c.PushList("copy")
			c.Set()
			c.Used()
			
			c.Push("i")
			c.Int(1)
			c.Add()
			c.Name("i")
		c.Redo()
		
		c.PushList("copy")
	},
}  
