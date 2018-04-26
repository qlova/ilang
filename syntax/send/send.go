package send

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/uct/compiler"

import "github.com/qlova/ilang/syntax/print"

var Name = compiler.Translatable{
	compiler.English: "send",
}

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		c.Expecting(symbols.FunctionCallBegin)
		
		c.List()
		c.Open()
		
		for {
			print.PrintType(c)
			
			switch c.Scan() {
				case symbols.ArgumentSeperator:
					continue
				case symbols.FunctionCallEnd:
					return
				default:
					c.Expected(symbols.ArgumentSeperator, symbols.FunctionCallEnd)
			}
		}
	},
}

