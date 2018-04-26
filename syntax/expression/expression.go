package expression

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/number"

import "github.com/qlova/ilang/syntax/symbols"

var Expression = compiler.Expression {
	Name: compiler.NoTranslation(symbols.SubExpressionStart),
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		var expression = c.ScanExpression()
		c.Expecting(")")
		return expression
	},
}

var Negative = compiler.Expression {
	Name: compiler.NoTranslation(symbols.Negative),
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		var value = c.ScanExpression()
		
		if !value.Equals(number.Type) {
			c.RaiseError(compiler.Translatable{
				compiler.English: "Only numbers can be negative!",
			})
		}
		
		c.Flip()
		return number.Type
	},
}
