package ifelse

import "github.com/qlova/uct/compiler"

var Else = struct { 
	Name compiler.Translatable 
	Flag compiler.Flag
	Statement compiler.Statement 
}{

	Name: compiler.Translatable{
			compiler.English: "else",
	},
}

func init() {	
	Else.Statement = compiler.Statement {
		Name: Else.Name,
		
		OnScan: func(c *compiler.Compiler) {
			
			if flag, ok := c.GetScope().Flags[If.Flag.Name[c.Language]]; ok {
				delete(c.GetScope().Flags, If.Flag.Name[c.Language]) 
				
				c.LoseScope()
				c.Or()
				//c.SetFlag(Else.Flag)
				
				//Need to remember how many if's there has been.
				if c.Peek() == If.Flag.Name[c.Language] {
					c.Scan()
					
					If.Statement.OnScan(c)
					flag.Value++
				} else {
					c.GainScope()
				}
				c.SetFlag(flag)

			} else {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Else statement must be used after an if statement!",
				})
			}
			
		},
	}
}
