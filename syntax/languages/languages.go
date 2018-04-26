package languages

import "github.com/qlova/uct/compiler"

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		switch c.Token() {
			case "Māori":
				
				c.Language = compiler.Maori
				
				return true
				
			default:
				return false
		}
	},
}
