package breakfast

import "github.com/qlova/uct/compiler"

var Name = compiler.Translatable{
	compiler.English: "break",
}

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		c.Done()
	},
}
