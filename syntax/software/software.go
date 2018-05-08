package software 

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/global"
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

var Exit = compiler.Statement {
	Name: compiler.Translatable{
		compiler.English: "exit",
	},
	 
	OnScan: func(c *compiler.Compiler) {
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
		
		
		c.Int(int64(len(c.GlobalScope.Variables)))
		c.Make()
		c.NameList("GLOBAL")
		c.PushList("GLOBAL")
		
		//Do the globals.
		for _, variable := range c.GlobalScope.Variables {
			var data = variable.Data.(global.Data)
			c.LoadCache(data.Cache, data.FileName, data.Line)
			var t = c.ScanExpression()
			c.Int(int64(data.Index)-1)
			t.Base.Attach(c)
		}
		
		c.DropList()
	},
}

var End = compiler.Statement {
	Name: compiler.NoTranslation(symbols.CodeBlockEnd),
	 
	OnScan: func(c *compiler.Compiler) {
		c.LoseScope()
	},
}
