package fixed 

import "github.com/qlova/uct/compiler"

var Name = compiler.Translatable{
	compiler.English: "fixed",
}

var Statement = compiler.Statement {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) {
		
		var line = c.Scanners[len(c.Scanners)-1].Line
		
		var name = c.Scan()
		c.Expecting("=")
		var cache = c.NewCache("", "\n")
		
		//TODO Do some sanity checking on the cache...
		
		c.RegisterExpression(compiler.Expression{
			Name: compiler.NoTranslation(name),
							 
			OnScan: func(c *compiler.Compiler) compiler.Type {
				
				c.LoadCache(cache, name, line)
				
				return c.ScanExpression()
			},
		})
	},
}
