package loop

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/uct/compiler"

var Name = compiler.Translatable{
	compiler.English: "loop",
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.Redo()
	},
}

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		c.Loop()
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

