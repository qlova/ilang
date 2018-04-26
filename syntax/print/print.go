package print

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/uct/compiler"

import "github.com/qlova/ilang/types/text"

var Name = compiler.Translatable{
	compiler.English: "print",
	compiler.Maori: "perehitia",
}

func PrintType(c *compiler.Compiler) {

	var t = c.ScanExpression() 
	
	c.Cast(t, text.Type)
	
	c.CopyPipe()
	c.Send()
}

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		c.Expecting(symbols.FunctionCallBegin)
		
		c.List()
		c.Open()
		
		for {
			PrintType(c)
			
			switch c.Scan() {
				case symbols.ArgumentSeperator:
					continue
				case symbols.FunctionCallEnd:
					//Print Newline
					c.List()
					c.Int('\n')
					c.Put()
					c.Send()
					
					return
				default:
					c.Expected(symbols.ArgumentSeperator, symbols.FunctionCallEnd)
			}
		}
	},
}

