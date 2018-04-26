package typer

import "github.com/qlova/uct/compiler"

import "github.com/qlova/ilang/types/thing"
import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/types/text"
import "github.com/qlova/ilang/types/list"

import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/symbols"

import "github.com/qlova/ilang/syntax/concept"

import "fmt"

//Decay the scope of the type. eg. Lose the member variables from the top level of scope.
func Decay(c *compiler.Compiler) {
	c.Expressions = c.Expressions[:len(c.Expressions)-1]
	c.Statements = c.Statements[:len(c.Statements)-1]
}

var tmp = 0
func Tmp() string {
	tmp++
	return fmt.Sprint(tmp)
}

func statement(c *compiler.Compiler, t compiler.Type, global bool) bool {
	// s/global/embed
	
	if global || (!global && c.GetVariable(c.Token()).Type.Equals(t)) {
		
		var name string
		if !global {
			c.PushType(t, c.Token())
			name = c.Token()
		}
		
		
		if global || (!global && c.Peek() == symbols.Index) {
			
			if !global {
				c.Scan()
				c.Scan()
			}
			
			for _, method := range t.Data.(thing.Data).Concepts {
				if method.Name[c.Language] == t.Name[c.Language]+"_dot_"+c.Token() {
					
					Absorb(c, t, true)
					var ret = concept.ScanCall(c, method)
					Decay(c)
					
					c.DropType(*ret)
					
					return true
				}
			}
			
			var subtype = thing.ScanStatement(c, t)
			
			if subtype.EmbeddedStatement != nil {
			
				if global {
					c.Push("pointer")
					c.Add()
				}
				
				subtype.EmbeddedStatement(c, subtype)
				
				if !global {
					c.DropList()
				}
				
				return true
			} else {
				c.Unimplemented()
			}
			
		} else if !global {
			
			
			switch c.Scan() {
				
				case symbols.Equals:
					var b = c.ScanExpression()
					if !b.Equals(t) {
						c.RaiseError(errors.AssignmentMismatch(b, t))
					}
					
					c.Name(name)

				default:
					c.Unexpected(c.Token())
			}
		}
		
		if !global {
			c.DropType(t)
		}
		
		return true
	}
	return false
}

func Absorb(c *compiler.Compiler, t compiler.Type, global bool) {
	c.RegisterExpression(compiler.Expression{		
		Detect: func(c *compiler.Compiler) *compiler.Type {
			
			for _, method := range t.Data.(thing.Data).Concepts {
				if method.Name[c.Language] == t.Name[c.Language]+"_dot_"+c.Token() {
					
					c.Push("pointer")
					
					var ret = concept.ScanCall(c, method)
					
					if ret == nil {
						c.RaiseError(compiler.Translatable{
							compiler.English: "Cannot use the concept "+method.Name[c.Language]+" inside a expression, no return values!",
						})
					}
					
					return ret
				}
			}
			
			if global || (!global && c.GetVariable(c.Token()).Type.Equals(t)) {
				
				if !global {
					c.PushType(t ,c.Token())
				} else {
					c.PushList("thing")
				}
				
				
				if global || (!global && c.Peek() == symbols.Index) {
					
					if !global {
						c.Scan()
						c.Scan()
					}
					
					var element = c.Token()
					if offset, ok := t.Data.(thing.Data).Offsets[element]; ok {
						
						var subtype = t.Data.(thing.Data).Elements[t.Data.(thing.Data).Map[element]]
						
						c.Int(int64(offset))
						if global {
							c.Push("pointer")
							c.Add()
						}
						
						
						
						//If it is a thing.
						if !NotThing(subtype) {
							return &subtype
						}
						
						c.Get()
						c.DropList()
						
						if !subtype.Equals(number.Type) {
							
							if subtype.Equals(text.Type) {
								
								c.HeapList()
								
							
							} else if _, ok := subtype.Data.(*list.Data); ok {
								
								c.HeapList()
								
								switch c.Peek() {
									case symbols.SelectMethod:
										
										c.Scan()
										return list.Shunt(subtype.Data.(*list.Data), c)
								}
								
								if list.CheckIndex(c, &subtype) {
									return &subtype.Data.(*list.Data).SubType
								}
								
							} else {
								c.Unimplemented()
							}
						}
						
						return &subtype
					} else if !ok {
						if global {
							

							return nil
						} else {
							
							//Maybe it is a concept?
							
							for _, method := range t.Data.(thing.Data).Concepts {
								if method.Name[c.Language] == t.Name[c.Language]+"_dot_"+element {
									var ret = concept.ScanCall(c, method)
									
									if ret == nil {
										c.RaiseError(compiler.Translatable{
											compiler.English: "Cannot use the concept "+method.Name[c.Language]+" inside a expression, no return values!",
										})
									}
									
									return ret
								}
							}
							
							c.RaiseError(errors.NoSuchElement(element, t))
						}
					}
					
					
				} else if !global { 
					return &t
					
				} else {
					c.Unimplemented()
				}
				
			}
			
			return nil
		},
	})

	//Work on garbage collection!
	c.RegisterStatement(compiler.Statement {
		Detect: func(c *compiler.Compiler) bool {
			return statement(c, t, global)
		},
	})
}
