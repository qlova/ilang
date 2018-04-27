package error 

import "github.com/qlova/uct/compiler"

import "github.com/qlova/ilang/syntax/symbols"

import "github.com/qlova/ilang/types/number"

var Name = compiler.Translatable{
	compiler.English: "error",
}

var Expression = compiler.Expression {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Push("ERROR")
		
		if c.ScanIf(symbols.Index) {
			c.Scan()
			
			for code, name := range Codes {
				if name[c.Language] == c.Token() {
					
					c.Expecting(symbols.FunctionCallBegin)
					c.Expecting(symbols.FunctionCallEnd)
					
					c.Int(int64(code))
					c.Same()
					return number.Type
				}
			}
			
			c.RaiseError(compiler.Translatable{
				compiler.English: "Invalid Error "+c.Token(),
			})
		}
		
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
	
