package ifelse

import "github.com/qlova/uct/compiler"

var End = compiler.Statement {
	Name: compiler.Translatable {
		compiler.English: "end",
	},
	 
	OnScan: func(c *compiler.Compiler) {
		c.LoseScope()
	},
} 
