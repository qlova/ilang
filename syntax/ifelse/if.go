package ifelse

import "github.com/qlova/uct/compiler"

var If = struct { 
	Name compiler.Translatable 
	Flag compiler.Flag
	Statement compiler.Statement 
}{

	Name: compiler.Translatable{
			compiler.English: "if",
	},
}

func init() { 
	
	If.Statement = compiler.Statement {
		Name: If.Name,
		
		OnScan: func(c *compiler.Compiler) {
			var test = c.ScanExpression()
			if test.Base == compiler.INT {
				
				c.If()
				c.GainScope()
				
				c.SetFlag(If.Flag)
				
				
			} else {
				c.Unimplemented()
			}
		},
	}
	
	If.Flag = compiler.Flag {
		Name: If.Name,
		
		OnLost: func(c *compiler.Compiler) {
			if flag, ok := c.GetScope().Flags[If.Flag.Name[c.Language]]; ok {
				for i := 0; i < flag.Value+1; i++ {
					c.No()
				}
			}
		},
	}
}
