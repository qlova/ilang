package binary

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/types/text"

var Name = compiler.Translatable{
	compiler.English: "binary",
}

var Expression = compiler.Expression {
	Detect: func(c *compiler.Compiler) *compiler.Type {
		if c.Token() == Name[c.Language] {
			c.Expecting("(")
			var value = c.ScanExpression()
			
			if !value.Equals(number.Type) {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Only numbers can be converted into binary!",
				})
			}
			
			c.Expecting(")")
			
			c.Int(2)
			c.Call(&text.Itoa)
			
			return &text.Type
		}
		return nil
	},
}
