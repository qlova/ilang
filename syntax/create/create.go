package create

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/types/connection"
import "github.com/qlova/ilang/types/text"

var Expression = compiler.Expression {
	Name: compiler.Translatable{
		compiler.English: "create",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Expecting(symbols.FunctionCallBegin)
		var argument = c.ScanExpression()
		c.Expecting(symbols.FunctionCallEnd)
		
		switch {
			case text.Type.Equals(argument):
				
				c.Open()
				c.CopyPipe()
				c.List()
				c.Send()

				return connection.Type
			
			default:
				c.Unimplemented()
				
		}
		
		return compiler.Type{}
	},
}

var Statement = compiler.Statement {
	Name: Expression.Name,
	
	OnScan: func(c *compiler.Compiler) {
		Expression.OnScan(c)
		c.FreePipe()
	},
}
