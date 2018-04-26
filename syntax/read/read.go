package read

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"

import "github.com/qlova/ilang/types/letter"
import "github.com/qlova/ilang/types/text"

var Expression = compiler.Expression {
	Name: compiler.Translatable{
		compiler.English: "read",
		compiler.Maori: "rīti",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Expecting(symbols.FunctionCallBegin)
		
		if c.ScanIf(symbols.FunctionCallEnd) {
			c.Int('\n')
			c.List()
			c.Open()
			c.Flip()
			c.Read()
			return text.Type
		}
		
		var argument = c.ScanExpression()
		c.Expecting(symbols.FunctionCallEnd)
		
		switch {
			case letter.Type.Equals(argument):
				
				c.List()
				c.Open()
				c.Flip()
				c.Read()
				
				return text.Type
			
			default:
				c.Unimplemented()
				
		}
		
		return compiler.Type{}
	},
}
