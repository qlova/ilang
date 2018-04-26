package thing

import "github.com/qlova/uct/compiler"

func Embed(t compiler.Type) *compiler.Function {
	
	if NotThing(t) {
		panic("Illegal call to thing.Embed(t)")
	}
	
	return &compiler.Function {
		Name: compiler.NoTranslation("embed_"+t.String()),

		//The list we are embedding in will be on the stack.
		Compile: func(c *compiler.Compiler) {
			c.Pull("thing_pointer")
			c.PullList("thing")
			
			c.Pull("pointer")
			
			for i := 0; i <  t.Data.(Data).Size; i++ {
				c.Push("pointer")
				c.Int(int64(i))
				c.Add()
				
				c.Push("thing_pointer") //pointer
				c.Int(int64(i))
				c.Add()
				c.PushList("thing")
				c.Get()
				c.DropList()
				
				c.Set()
			}
		},
	}
}
