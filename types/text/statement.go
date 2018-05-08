package text

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/symbols"

import "github.com/qlova/ilang/types/number"

func init() {
	Type.EmbeddedStatement = func(c *compiler.Compiler, list compiler.Type) {
		statement(c, true)
	}
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		return statement(c, false)
	},
}

func statement(c *compiler.Compiler, embed bool) bool {
	if embed || c.GetVariable(c.Token()).Type.Equals(Type) {
			
			var name = c.Token()
			
			switch c.Scan() {
				
				case symbols.Equals:
					if embed {
						//Garbage collect! >:)
						c.Get()
						c.Copy()
						c.If()
							c.Flip()
							c.HeapList()
						c.Or()
							c.Drop()
						c.No()
					}
					
					var t = c.ScanExpression()
					if !t.Equals(Type) {
						c.RaiseError(errors.AssignmentMismatch(t, Type))
					}
					
					if embed {
						c.Int(0)
						c.HeapList()
						c.Set()
					} else {
						
						c.NameList(name)
					
					}
					
				case symbols.Plus:
					if embed {
						c.Unimplemented()
					}
					
					c.Expecting(symbols.Equals)
					
					c.PushList(name)
					
					var t = c.ScanExpression()
					
					//string += 1 (increments the numeric value encoded in string)
					if t.Equals(number.Type) {
						c.Int(10)
						c.Call(&Atoi)
						c.Add()
						c.Int(10)
						c.Call(&Itoa)
						c.NameList(name)
						
						return true
					}
					
					if !t.Equals(Type) {
						c.RaiseError(errors.AssignmentMismatch(t, Type))
					}
					
					c.Call(&Join)
					
					c.NameList(name)

				default:
					c.Unexpected(c.Token())
			}
			
			return true
		}
		return false
}
 

