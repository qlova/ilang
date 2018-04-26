package error 

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/number"

var Name = compiler.Translatable{
	compiler.English: "error",
}

var Expression = compiler.Expression {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Push("ERROR")
		return number.Type
	},
}

var Statement = compiler.Statement {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) {
		c.NextToken = "ERROR"
		c.Scan()
		
		if !c.GetVariable("ERROR").Defined {
			c.SetVariable("ERROR", number.Type)
		}
		
		number.Statement.Detect(c)
	},
}
	
