package expression

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/number"

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

var Expression = compiler.Expression {
	Name: compiler.NoTranslation(symbols.SubExpressionStart),
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		var expression = c.ScanExpression()
		c.Expecting(")")
		return expression
	},
}

var NumberOf = compiler.Expression{
	Name: compiler.NoTranslation(symbols.NumberOf),
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		var list = c.Shunt(c.Expression(), 5)
		
		if list.Base != compiler.LIST {
			c.RaiseError(compiler.Translatable{
				compiler.English: "Can only get the number of elements for a list type!",
			})
		}
		c.Size()
		c.DropType(list)
		
		return number.Type
	},
}

var Negative = compiler.Expression {
	Name: compiler.NoTranslation(symbols.Negative),
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		var value = c.ScanExpression()
		
		if value.Base != compiler.INT {
			c.RaiseError(errors.MustBeNumeric(value))
		}
		
		c.Flip()
		return value
	},
}
