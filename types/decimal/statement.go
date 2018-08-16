package decimal

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/symbols"

import "github.com/qlova/ilang/types/number"

func statement(c *compiler.Compiler, embed bool) bool {
	if embed || c.GetVariable(c.Token()).Type.Equals(Type) {
			
			var decimal = c.GetVariable(c.Token()).Type
			
			var name = c.Token()
			
			switch c.Scan() {
				
				case symbols.Equals:
					var t = c.ScanExpression()
					if !t.Equals(decimal) {
						c.RaiseError(errors.AssignmentMismatch(t, decimal))
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
								
								if t.Equals(number.Type) {
									
									var exponent = DefaultExponent
									if decimal.Data != nil {
										exponent = decimal.Data.(Data).Exponent
									}
									
									c.BigInt(exponent)
									c.Mul()
								} else {
								
									c.RaiseError(errors.AssignmentMismatch(t, decimal))
								}
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
							if !t.Equals(decimal) {
								
								
								
								c.RaiseError(errors.AssignmentMismatch(t, decimal))
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
							if !t.Equals(decimal) {
								c.RaiseError(errors.AssignmentMismatch(t, decimal))
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
							if !t.Equals(decimal) {
								c.RaiseError(errors.AssignmentMismatch(t, decimal))
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
 
