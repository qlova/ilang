package connection

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"

var Name = compiler.Translatable{
	compiler.English: "connection",
}

var Type = compiler.Type {
	Name: Name,
	
	Base: compiler.PIPE,
	
	Shunts: compiler.Shunt{
		symbols.Index: func(c *compiler.Compiler, b compiler.Type) compiler.Type {
			switch b.Name[c.Language] {
				case "read":
					c.Unimplemented()
				default:
					c.RaiseError(compiler.Translatable{
						compiler.English: "No such method "+b.Name[c.Language],
					})
			}
			
			return compiler.Type{}
		},
	},
}

var Statement = compiler.Statement {
	
	OnScan: func(c *compiler.Compiler) {
		c.Expecting(symbols.Index)
			
		switch c.Scan() {
			case "write":
				c.Expecting(symbols.FunctionCallBegin)
				var arg = c.ScanExpression()
				c.Expecting(symbols.FunctionCallEnd)
				
				if arg.Base != compiler.LIST {
					c.RaiseError(compiler.Translatable{
						compiler.English: "Cannot write "+arg.Name[c.Language]+" type to connection.",
					})
				}
				
				c.Send()
			default:
				c.RaiseError(compiler.Translatable{
					compiler.English: "No such method "+c.Token(),
				})
		}
	},
}

func init() {
	Statement.Detect = func(c *compiler.Compiler) bool {
		
		if c.GetVariable(c.Token()).Type.Equals(Type) {
			c.PushPipe(c.Token())
			Statement.OnScan(c)
			return true
		}
		
		return false
	}
}
