package global

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/list"

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


func Init(c *compiler.Compiler) {
	if len(c.GlobalScope.Variables) > 0 {
		
		c.Int(int64(len(c.GlobalScope.Variables))+2)
		c.Make()
		c.NameList("GLOBAL")
		c.PushList("GLOBAL")
		
		//Do the globals.
		//TODO fix race condition that occurs when a global variable modifies another global variable.
		for _, variable := range c.GlobalScope.Variables {
			var data = variable.Data.(Data)
			c.LoadCache(data.Cache, data.FileName, data.Line)
			var t = c.ScanExpression()
			c.Int(int64(data.Index)-1)
			t.Base.Attach(c)
		}
		
		c.DropList()
	}
}

var Expression = compiler.Expression {
	Detect: func(c *compiler.Compiler) *compiler.Type {
		if global, ok := c.GlobalScope.Variables[c.Token()]; ok {
			c.PushList("GLOBAL")
			c.Int(int64(global.Data.(Data).Index+1))
			c.Get()
			c.DropList()
			
			if global.Data.(Data).Type.Base == compiler.LIST {
				c.HeapList()
			}
			
			var t = global.Data.(Data).Type 
			
			if t.Equals(list.Type) {
				if list.CheckIndex(c, &t) {
					return &t.Data.(*list.Data).SubType
				}
			}
			
			return &t
		}
		
		return nil
	},
}

//Global statements act on the embedded GLOBAL list in UCT.
var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool { 
		if global, ok := c.GlobalScope.Variables[c.Token()]; ok {
			c.PushList("GLOBAL")
			c.Int(int64(global.Data.(Data).Index+1))
			
			if global.Data.(Data).Type.EmbeddedStatement == nil {
				c.Unimplemented()
			}
			
			global.Data.(Data).Type.EmbeddedStatement(c, global.Data.(Data).Type)
			
			return true
		}
		
		return false
	},
}
