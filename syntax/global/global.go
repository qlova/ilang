package global

import "github.com/qlova/uct/compiler"

type Data struct {
	Type compiler.Type
	Cache compiler.Cache
	Line int
	FileName string
	Index int
}

func (Data) Name(compiler.Language) string {
	return ""
}

func (Data) Equals(compiler.Data) bool {
	return false
}

var Type = compiler.Type{}


var Expression = compiler.Expression {
	Detect: func(c *compiler.Compiler) *compiler.Type {
		if global, ok := c.GlobalScope.Variables[c.Token()]; ok {
			c.PushList("GLOBAL")
			c.Int(int64(global.Data.(Data).Index-1))
			c.Get()
			c.DropList()
			
			if global.Data.(Data).Type.Base == compiler.LIST {
				c.HeapList()
			}
			
			var t = global.Data.(Data).Type 
			
			return &t
		}
		
		return nil
	},
}
