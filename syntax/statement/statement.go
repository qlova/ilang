package statement

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/global"
import "github.com/qlova/uct/compiler"

import "io/ioutil"

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if !c.GetVariable(c.Token()).Defined {
			
			var name = c.Token()
			
			//Don't reassign globals?
			for names := range c.GlobalScope.Variables {
				if names == name {
					return false
				}
			}
			
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
			
			//Global!
			if len(c.Scope) == 0 {
				if len(names) > 0 {
					c.Unimplemented()
				}
				
				var cache = c.NewCache("", "\n")
				
				c.LoadCache(cache, "statement.go", 0)
				
				output := c.Output
				c.Output = ioutil.Discard
				
				var t = c.ScanExpression()

				c.Output = output
				
				c.SetGlobal(name, global.Type.With(global.Data{
					Type: t,
					Cache: cache,
					Line: 0,
					FileName: "statement.go",
					Index: len(c.GlobalScope.Variables),
				}))
				

				return true
			}
			
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
