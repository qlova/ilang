package software 

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/uct/compiler"

var Name = compiler.Translatable{
	compiler.English: "software",
	compiler.Maori: "taupānga",
	compiler.Chinese: "软件",
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.Exit()
	},
}
 
var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		c.Main()
		c.Expecting(symbols.CodeBlockBegin)

		c.GainScope()
		c.SetFlag(Flag)
	},
}

var End = compiler.Statement {
	Name: compiler.NoTranslation(symbols.CodeBlockEnd),
	 
	OnScan: func(c *compiler.Compiler) {
		c.LoseScope()
	},
}
