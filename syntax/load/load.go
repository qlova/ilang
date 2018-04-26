package load

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/text"

var Expression = compiler.Expression {
	Name: compiler.Translatable{
		compiler.English: "load",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Expecting("(")
		var argument = c.ScanExpression()
		c.Expecting(")")
		
		switch {
			case text.Type.Equals(argument):
				
				c.Load()

				return text.Type
			
			default:
				c.Unimplemented()
				
		}
		
		return compiler.Type{}
	},
}
