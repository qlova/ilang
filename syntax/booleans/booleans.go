package booleans

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/syntax/errors"

var True = compiler.Expression {
	Name:  compiler.Translatable{
		compiler.English: "true",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		
		c.Int(1)
		
		return number.Type
	},
}

var False = compiler.Expression {
	Name:  compiler.Translatable{
		compiler.English: "false",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		
		c.Int(0)
		
		return number.Type
	},
}

var Not = compiler.Expression {
	Name:  compiler.Translatable{
		compiler.English: "not",
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		
		c.Expecting("(")
		var value = c.ScanExpression()
		if !value.Equals(number.Type) {
			c.RaiseError(errors.ExpectingType(number.Type, value))
		}
		c.Expecting(")")
		
		c.Int(0)
		c.Div()
		
		return number.Type
	},
}
