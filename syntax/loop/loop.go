package loop

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/number"

var Name = compiler.Translatable{
	compiler.English: "loop",
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.Redo()
		
		if _, ok := c.Scope[len(c.Scope)-2].Flags[Name[c.Language]]; ok {
			c.Copy()
			c.If()
				c.Int(1)
				c.Sub()
				c.Done()
			c.No()
			c.Drop()
		}
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

var Break = compiler.Statement {
	Name: compiler.Translatable{
		compiler.English: "break",
	},
	 
	OnScan: func(c *compiler.Compiler) {
		
		if c.Peek() != "\n" {
			c.ExpectingType(number.Type)
			
			if _, ok := c.Scope[len(c.Scope)-2].Flags[Name[c.Language]]; ok {
				c.Int(1)
				c.Sub()
			}
			c.Done()
			
			return
		}

		if _, ok := c.Scope[len(c.Scope)-2].Flags[Name[c.Language]]; ok {
			c.Int(0)
		}
		
		c.Done()
	},
}
