package decimal

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/symbols"

func statement(c *compiler.Compiler, embed bool) bool {
	if embed || c.GetVariable(c.Token()).Type.Equals(Type) {
			
			var name = c.Token()
			
			switch c.Scan() {
				
				case symbols.Equals:
					var t = c.ScanExpression()
					if !t.Equals(Type) {
						c.RaiseError(errors.AssignmentMismatch(t, Type))
					}
					
					if embed {
						c.Set()
					} else {
						c.Name(name)
					}
					
				case symbols.Plus:
					switch c.Scan() {
						case symbols.Plus:
							
							if embed {
								c.Get()
							} else {
								c.Push(name)
							}

							c.Int(1)
							c.Add()
							
							if embed {
								c.Set()
							} else {
								c.Name(name)
							}
						
						case symbols.Equals:
							
							if embed {
								c.Get()
							} else {
								c.Push(name)
							}
							
							var t = c.ScanExpression()
							if !t.Equals(Type) {
								c.RaiseError(errors.AssignmentMismatch(t, Type))
							}
							
							c.Add()
							
							if embed {
								c.Set()
							} else {
								c.Name(name)
							}
						
						default:
							c.Unexpected(c.Token())
					}
				
				case symbols.Minus:
					switch c.Scan() {
						case symbols.Minus:
							if embed {
								c.Get()
							} else {
								c.Push(name)
							}
							
							c.Int(1)
							c.Sub()
							
							if embed {
								c.Set()
							} else {
								c.Name(name)
							}
						
						case symbols.Equals:
							c.Push(name)
							
							var t = c.ScanExpression()
							if !t.Equals(Type) {
								c.RaiseError(errors.AssignmentMismatch(t, Type))
							}
							
							
							c.Sub()
							c.Name(name)
							
							if embed {
								c.Unimplemented()
							}
						
						default:
							c.Unexpected(c.Token())
					}
				
				case symbols.Times:
					switch c.Scan() {						
						case symbols.Equals:
							var t = c.ScanExpression()
							if !t.Equals(Type) {
								c.RaiseError(errors.AssignmentMismatch(t, Type))
							}
							
							c.Push(name)
							c.Mul()
							c.Name(name)
							
							if embed {
								c.Unimplemented()
							}
						
						default:
							c.Unexpected(c.Token())
					}
					
				case symbols.Divide:
					switch c.Scan() {						
						case symbols.Equals:
							c.Push(name)
							
							var t = c.ScanExpression()
							if !t.Equals(Type) {
								c.RaiseError(errors.AssignmentMismatch(t, Type))
							}
							
							
							c.Div()
							c.Name(name)
							
							if embed {
								c.Unimplemented()
							}
						
						default:
							c.Unexpected(c.Token())
					}
				
				default:
					c.Unexpected(c.Token())
			}
			
			return true
		}
		return false
}
 
