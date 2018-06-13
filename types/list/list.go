package list

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/types/text"
import "github.com/qlova/ilang/types/thing"

import "fmt"

var unique int
func Tmp() string {
	unique++
	return fmt.Sprint("list", unique)
}

type Data struct {
	Step int
	SubType compiler.Type
}

func (d *Data) Name(l compiler.Language) string {
	return d.SubType.Name[l]
}

func (d *Data) Equals(b compiler.Data) bool {
	return d.SubType.Equals(b.(*Data).SubType)
}

//Add type t that is sitting on the top of the stack to list l that is sitting next on the stack.
func (d *Data) Add(c *compiler.Compiler, t compiler.Type) {
	
	if d.SubType.Name[c.Language] == "" {
		d.SubType = t
	}
	
	if t.Equals(number.Type) {
		
		c.Put()

		return
	}
	
	if t.Equals(text.Type) {

		c.Int(0)
		c.HeapList()
		c.Put()
	
		return
	}
	
	c.Unimplemented()
}

var Type = compiler.Type {
	Name: compiler.Translatable{
		compiler.English: "list",
	},
	
	Base: compiler.LIST,
	
	Casts: []func(*compiler.Compiler, compiler.Type) bool {
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(text.Type) {
				return true
			}
			return false
		},
	},
}

func Is(t compiler.Type) bool {
	if t.Name[0] == Type.Name[0] {
		return true
	}
	return false
}

func SubType(t compiler.Type) compiler.Type {
	return t.Data.(*Data).SubType
}


func Of(t compiler.Type) compiler.Type {
	var r = Type
	
	r.Data = &Data{
		Step: 1,
		SubType: t,
	}
	
	return r
}

func Shunt(data *Data, c *compiler.Compiler) *compiler.Type {
	switch c.Token() {

		case symbols.SelectMethod:
			
			c.Scan()
			
			if c.Token() == "size" {
				
				c.Expecting(symbols.FunctionCallBegin)
				c.Expecting(symbols.FunctionCallEnd)
				
				c.Size()
				c.Used()
				
				return &number.Type
			} else if c.Token() == "copy" {
				
				c.Expecting(symbols.FunctionCallBegin)
				c.Expecting(symbols.FunctionCallEnd)
				
				c.Call(&Copy)
				
				var l = Type
				l.Data = data
				
				AddShunts(&l)
				
				return &l
			} else {
				c.Unimplemented()
			}
	}
	
	c.Unexpected(c.Token())
	return nil
}

var Expression = compiler.Expression {
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		if c.Token() == Type.Name[c.Language] {
			if c.Peek() == "." {
				c.Scan()
				
				var sub = c.Scan()
				if t := c.GetType(sub); t == nil {
					c.RaiseError(compiler.Translatable{
						compiler.English: "Cannot create list."+sub+", type "+sub+" does not exist!",
					})
				} else {
					
					c.Expecting("(")
					c.Expecting(")")
					
					var list = Type
					
					list.Data = &Data {
						Step: 1,
						SubType: *t,
					}
					
					c.List()
					
					if CheckIndex(c, &list) {
						return &list.Data.(*Data).SubType
					}
						
					AddShunts(&list)
						
					return &list
				}
			} 
			return nil
		}

		//Shunt here.
		if c.GetVariable(c.Token()).Type.Equals(Type) {
			
			var name = c.Token()
			var list = c.GetVariable(c.Token()).Type
			
			switch c.Peek() {
				case symbols.SelectMethod:
					
					c.PushList(name)
					
					c.Scan()
					return Shunt(list.Data.(*Data), c)
			}
			
			c.PushList(name)
			
			if CheckIndex(c, &list) {
				return &list.Data.(*Data).SubType
			}
			
			AddShunts(&list)
			
			return &list
		}
		
		//List literal.
		if c.Token() == symbols.ListBegin {
			
			if c.Peek() == symbols.ListEnd  {
				c.Scan()
				c.List()
				
				return &Type
			} else {
				
				c.List()
				
				var list = Type
				list.Data = &Data{}
				
				for {
					list.Data.(*Data).Add(c, c.ScanExpression())
					
					if c.Peek() == symbols.ListEnd {
						break
					} else {
						c.Expecting(symbols.ListSeperator)
					}
				}
				
				c.Expecting(symbols.ListEnd)
				
				if CheckIndex(c, &list) {
					return &list.Data.(*Data).SubType
				}
				
				AddShunts(&list)
				
				return &list
			}
			
			return nil
		} else {
			return nil
		}
	},
}

func statement(c *compiler.Compiler, list compiler.Type, embed bool) bool {
	if embed || c.GetVariable(c.Token()).Type.Equals(Type) {
		
		var name = c.Token()
		
		switch c.Scan() {
			
			//Beware of subtype woes.
			case symbols.Equals:
				var t = c.ScanExpression()
				if !t.Equals(Type) {
					c.RaiseError(errors.AssignmentMismatch(t, Type))
				}
				
				if embed {
					c.Int(0)
					c.HeapList()
					c.Set()
				} else {
					c.NameList(c.Token())
				}
			
			// list += value
			case symbols.Plus:
				c.Expecting(symbols.Equals)
				
				if embed {
					c.Get()
					c.HeapList()
				}
				
				//Could implement lookahead here...
				var t = c.ScanExpression()
				
				//Auto update type.
				if list.Data == nil {
					list.Data = &Data {
						Step: 1,
						SubType: t,
					}

					c.UpdateVariable(name, list)
				}
				
				if list.Data.(*Data).SubType.Equals(number.Type) {
					
					c.PushList(name)
					c.Put()
					c.NameList(name)
					
				} else {
					
					if t.Base != compiler.INT {
						
						//uh oh, things are so painful.
						if _, ok := t.Base.(thing.Base); ok {
							
							for i := 0; i <  t.Data.(thing.Data).Size; i++ {
								c.Copy() //pointer
								c.Int(int64(i))
								c.Add()
								c.Get()
								
								if embed {
									c.SwapList()
								} else {
									c.PushList(name)
								}
								c.Put()
								if embed {
									c.SwapList()
								} else {
									c.DropList()
								}
							}
							
							c.Drop()
							c.DropList()
							
							if !embed {
								list.Data = &Data {
									Step: t.Data.(thing.Data).Size,
									SubType: t,
								}
								c.UpdateVariable(name, list)
							}
							
							return true
						}
						
						c.Unimplemented()
					}
					
					c.Unimplemented()
				}
			
			//list[0]
			case symbols.ListBegin:
				
				if !embed {
					c.PushList(name)
				} else {
					c.Get()
					c.HeapList()
				}
				
				if !c.ScanExpression().Equals(number.Type) {
					c.RaiseError(compiler.Translatable{
						compiler.English: "Lists can only be indexed with the number type.",
					})
				}
				
				c.Expecting("]")
				
				if int64(list.Data.(*Data).Step) == 0 {
					list.Data.(*Data).Step = 1
				}
				c.Size()
				c.If()
					c.Copy()
					c.Int(int64(list.Data.(*Data).Step))
					c.Mul()
					c.Size()
					c.Mod()
					
					if list.Data.(*Data).SubType.EmbeddedStatement == nil {
						println(list.Data.(*Data).SubType.String())
						c.Unimplemented()
					}
					
					list.Data.(*Data).SubType.EmbeddedStatement(c, list.Data.(*Data).SubType)
				c.No()
				
				c.DropList()
				c.Drop()
					
				//c.DropList()
			
			default:
				c.Unexpected(c.Token())
		}
		
		return true
	}
	return false
}

var Statement = compiler.Statement {
	
	Detect: func(c *compiler.Compiler) bool {
		return statement(c, c.GetVariable(c.Token()).Type, false)
	},
}

func CheckIndex(c *compiler.Compiler, t *compiler.Type) bool {
	if c.Peek() == symbols.IndexBegin {
		c.Expecting(symbols.IndexBegin)
		var b = c.ScanExpression()
		c.Expecting(symbols.IndexEnd)
		
		if !b.Equals(number.Type) {
			c.RaiseError(compiler.Translatable{
				compiler.English: "Lists can only be indexed with the number type.",
			})
		}

		
		var subtype = t.Data.(*Data).SubType
		
		if subtype.Equals(number.Type) || subtype.Equals(text.Type) {
			
			c.Size()
			c.If()
				c.Size()
				c.Mod()
				c.Get()
				c.Used()
				
				if subtype.Equals(text.Type) {
					c.HeapList()
				}
			c.Or()
				//Return a blank version of the type.
				
				if subtype.Equals(text.Type) {
					c.List()
				} else {
					c.Int(0)
				}
			
			c.No()
		
		} else if data, ok := subtype.Data.(thing.Data); ok {
			
			c.Size()
			c.If()
				c.Int(int64(data.Size))
				c.Mul()
				c.Size()
				c.Mod()
			c.Or()
				
				//TODO Return a blank version of the type.
				c.List()
			
			c.No()
			
		} else {
			c.Unimplemented()
		}
		
		return true
	}
	
	return false
}

func AddShunts(list *compiler.Type) {
	list.Shunts = compiler.Shunt {
		symbols.Plus: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
			if !list.Equals(t) {
				c.RaiseError(errors.Single(*list, symbols.Plus,t))
			}
			
			c.Call(&text.Join)
			
			return t
		},
		
		symbols.SelectMethod: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		
			if t.Name[c.Language] == "size" {
				
				c.Expecting(symbols.FunctionCallBegin)
				c.Expecting(symbols.FunctionCallEnd)
				
				c.Size()
				c.Used()
				
				return number.Type
			} else if t.Name[c.Language] == "copy" {
				
				c.Expecting(symbols.FunctionCallBegin)
				c.Expecting(symbols.FunctionCallEnd)
				
				c.Call(&Copy)
				
				return *list
			} else {
				c.Unimplemented()
			}
			return *list
		},
	}
}

func init() {
	
	Type.Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) *compiler.Type {
		if symbol == symbols.Plus {
			if a.Equals(b) {
				
				c.Call(&text.Join)

				return &a
			}
		}
		return nil
	}
	
	Type.EmbeddedStatement = func(c *compiler.Compiler, list compiler.Type) {
		statement(c, list, true)
	}
	
	text.Type.Casts = append(text.Type.Casts, 
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(Type) {
				return true
			}
			return false			
		},
	)
	
	number.Type.Casts = append(number.Type.Casts, 
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(Type) {
				c.Make()
				return true
			}
			return false			
		},
	)
}
