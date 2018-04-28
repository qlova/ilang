package maths

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

var Abs = compiler.Expression {
	Name: compiler.Translatable{
		compiler.English: "abs",
	},

	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Expecting(symbols.FunctionCallBegin)
		var value = c.ScanExpression()
		c.Expecting(symbols.FunctionCallEnd)
		
		if value.Base != compiler.INT {
			c.RaiseError(errors.MustBeNumeric(value))
		}
		
		c.Copy()
		c.Int(0)
		c.More()
		c.If()
			c.Flip()
		c.No()
		
		return value
	},
}
