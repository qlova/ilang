package statement

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/uct/compiler"

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if !c.GetVariable(c.Token()).Defined {
			
			var name = c.Token()
			if errors.IsInvalidName(name) {
				c.RaiseError(errors.InvalidName(name))
			}
			
			//Define the variable. Maybe we should explain why we are expecting an '='
			if c.Peek() != symbols.Equals &&  c.Peek() != symbols.ArgumentSeperator {
				
				//Nicer errors.
				switch c.Peek() {
					case symbols.LessThan, symbols.MoreThan, 
						symbols.IndexBegin, symbols.Plus, symbols.Minus:
						
							c.UndefinedError(c.Token())
				}
				
				return false
			}
			
			var names []string
			for c.Peek() == "," {
				c.Scan()
				names = append(names, c.Scan())
			}
			
			c.Expecting("=")
			
			
			var t = c.ScanExpression()
			
			c.PullType(t, name)
			
			c.SetVariable(name, t)
			
			if len(names) > 0 {
				c.Expecting(",")
			}
			
			for i, name := range names {
				var t = c.ScanExpression()
				
				println(t.String())
			
				c.PullType(t, name)
				
				c.SetVariable(name, t)
				
				if i != len(names)-1 {
					c.Expecting(",")
				}
			}
			
			return true
		}
		
		return false
	},
}
