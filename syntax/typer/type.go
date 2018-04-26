/*
 * Oh god this is a messy module, I'm sorry.
 * At least it works ;)
 * 
 */

package typer

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/convert"
import "github.com/qlova/ilang/syntax/content"
import "github.com/qlova/ilang/syntax/concept"

import "github.com/qlova/ilang/types/thing"
import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/types/text"
import "github.com/qlova/ilang/types/list"

import "github.com/qlova/uct/compiler"

var Name = compiler.Translatable {
	compiler.English: "type",
}

var Statement = compiler.Statement {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) {
		var name = c.Scan()
		if errors.IsInvalidName(name) {
			c.RaiseError(errors.InvalidName(name))
		}
		
		var t = thing.Type
		t.Name[c.Language] = name
		
		var data = thing.Data{Offsets:make(map[string]int), Map:make(map[string]int)}
		
		c.Expecting(symbols.CodeBlockBegin)
		
		c.Code(name)
		
		c.Pull("pointer")
		
		c.GainScope()
		
		var i = 0
		var offset = 0
		
		//We scan the member variables and then the concpets and conversions.
		var MethodMode = false
		
		for {
			var element = c.Scan()
			if element == "\n" {
				continue
			}
			if element == "}" {
				break
			}
			

			//Ugly
			if element == convert.Name[c.Language] {
				if !MethodMode {
					c.LoseScope()
					c.Back()
				}
				
				c.Scan()
				t.Data = data
				Absorb(c, t, true)
				convert.New(c, &t)
				Decay(c)
				
				MethodMode = true

				continue
				
			} else if element == concept.Name[c.Language] {	
				
				if !MethodMode {
					c.LoseScope()
					c.Back()
				}
				
				c.Scan()
				t.Data = data
				Absorb(c, t, true)
				c.NextToken = t.Name[c.Language]+"_dot_"+c.Token()
				
				concept.Statement.OnScan(c)
				Decay(c)
				
				data.Concepts = append(data.Concepts, c.Functions[len(c.Functions)-1])
				
				MethodMode = true

				continue
			
			} else if element == content.Name[c.Language] {
				if !MethodMode {
					c.LoseScope()
					c.Back()
				}
				
				c.Scan()
				t.Data = data
				Absorb(c, t, true)
				content.New(c, &t)
				Decay(c)
				
				MethodMode = true

				continue
			}
			
			
			if MethodMode {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Cannot place member variables at the end of the type!",
				})
			}
			
			if errors.IsInvalidName(element) {
				c.RaiseError(errors.InvalidName(element))
			}
			
			if c.Peek() == "\n" {
				data.Elements = append(data.Elements, number.Type)
				data.Offsets[element] = offset
				data.Map[element] = i
				
				i += 1
				offset += 1
				continue
			}
			
			c.Expecting("=")
			
			
			//Lets try with more complex types.
			//Need to detect if this will be a subtype....
			var CheckIfThing = c.Peek()
			
			if getType := c.GetType(CheckIfThing); getType != nil && getType.Data != nil {
				if thingdata, ok := getType.Data.(thing.Data); ok {
					
					c.Scan()
					
					c.Push("pointer")
					c.Int(int64(offset))
					c.Add()
					
					//Oh great... it's a literal...
					if c.Peek() == symbols.CodeBlockBegin {
						//urgrgss
						c.Scan()

						thing.ScanEmbeddedLiteral(c, *getType)
					} else {
						
						c.Expecting("(")
						c.Expecting(")")
						
						c.CallRaw(getType.Name[c.Language])
					}
					
					data.Elements = append(data.Elements, *getType)
					data.Offsets[element] = offset
					data.Map[element] = i
					
					i += 1
					offset += thingdata.Size
					continue
				}
			}
			
			
			c.Int(int64(offset))
			c.Push("pointer")
			c.Add()
			
			var defaulting = c.ScanExpression()
			
			if defaulting.Equals(text.Type) {

				c.Int(0)
				c.HeapList()
				c.Set()
			
			} else if defaulting.Equals(number.Type) {
				
				c.Set()
			
			}  else if _, ok := defaulting.Data.(*list.Data); ok {
				
				c.Int(0)
				c.HeapList()
				c.Set()
				
			} else if defaulting.Equals(list.Type) {
				
				c.RaiseError(compiler.Translatable{
					compiler.English: "List assignment is ambigious, please use list.subtype()",
				})
				
			} else {
				c.Unimplemented()
			}
			
			
			
			data.Elements = append(data.Elements, defaulting)
			data.Offsets[element] = offset
			data.Map[element] = i
			
			i += 1
			offset += 1
		}
		
		//No methods!
		if !MethodMode {
			c.LoseScope()
			c.Back()
		}
		
		data.Size = offset
		t.Data = data
		
		c.RegisterExpression(compiler.Expression{
			Name: t.Name,
			
			OnScan: func(c *compiler.Compiler) compiler.Type {
				
				//Call methods on empty types globally!
				// eg. Math.add(a, b)
				if c.Peek() == symbols.Index {
					c.Scan()
					c.Scan()
					
					c.Int(0)
					//TODO make this a function.
					for _, method := range t.Data.(thing.Data).Concepts {
						if method.Name[c.Language] == t.Name[c.Language]+"_dot_"+c.Token() {
							var ret = concept.ScanCall(c, method)
							
							if ret == nil {
								c.RaiseError(compiler.Translatable{
									compiler.English: "Cannot use the concept "+method.Name[c.Language]+" inside a expression, no return values!",
								})
							}
							
							return *ret
						}
					}
					
					c.RaiseError(compiler.Translatable{
						compiler.English: "Cannot use the concept "+c.Token()+" does not exist!",
					})
				}
				
				
				c.Expecting(symbols.FunctionCallBegin)
				
				
				//Content-packed calls.
				if c.Peek() != symbols.FunctionCallEnd {
					var content = c.ScanExpression()
					c.Expecting(symbols.FunctionCallEnd)
					
					c.Int(int64(offset+1))
					c.Make()
					c.Int(0)
					
					var cast_name = content.Name[c.Language]+"_"+convert.Name[c.Language]+"_"+t.Name[c.Language]
					c.CallRaw(cast_name)
					
					c.Int(0)
					return t
				}
				
				c.Expecting(symbols.FunctionCallEnd)
				
				c.Int(int64(offset))
				c.Make()
				c.Int(0)
				c.CallRaw(name)
				
				//Pull pointer
				c.Int(0)
				
				return t
			},
		})
		
		
		 t.Shunts = compiler.Shunt {
			symbols.Index: func (c *compiler.Compiler, b compiler.Type) compiler.Type {
				
				var index = b.Name[c.Language]
				
				if _, ok := t.Data.(thing.Data).Map[index]; !ok {
					
					//Maybe it is a concept?
					
					for _, method := range t.Data.(thing.Data).Concepts {
						println(method.Name[c.Language])
						if method.Name[c.Language] == t.Name[c.Language]+"_dot_"+index {
							concept.ScanCall(c, method)
							
							if len(method.Returns) == 0 {
								c.RaiseError(compiler.Translatable{
									compiler.English: "Cannot use the concept "+method.Name[c.Language]+" inside a expression, no return values!",
								})
							}
							
							return method.Returns[0]
						}
					}
					
					c.RaiseError(errors.NoSuchElement(index, t))
				}
				
				var subtype = t.Data.(thing.Data).Elements[t.Data.(thing.Data).Map[index]]
				var offset = t.Data.(thing.Data).Offsets[index]
				
				c.Int(int64(offset))
				c.Add()
				
				//If it is a thing.
				if !NotThing(subtype) {
					
					return subtype
				}
				
				c.Get()
				c.DropList()
				
				if !subtype.Equals(number.Type) {
					
					if subtype.Equals(text.Type) {
						
						c.HeapList()
						
					} else {
						c.Unimplemented()
					}
				}
				
				return subtype
			},
		}
		
		t.Collect = func(c *compiler.Compiler) {
			for element, offset := range data.Offsets {
				var subtype = data.Elements[data.Map[element]]
				
				if subtype.Base == compiler.LIST && NotThing(subtype) {
					
					c.Int(int64(offset))
					c.Get()
					c.Flip()
					c.HeapList()
					
				} else if !NotThing(subtype) {
					
					c.Int(int64(offset))
					thing.Collect(c, subtype)
					
				}
			}
		}
		
		var Embed = thing.Embed(t)
		
		//eg. List[0] = Thing() or thing.Thing = Thing()
		t.EmbeddedStatement = func(c *compiler.Compiler, t compiler.Type) {
			if c.Peek() == "=" {
				c.Scan()
				
				if expr := c.ScanExpression(); !expr.Equals(t) {
					c.RaiseError(errors.ExpectingType(t, expr))
				}
				
				//Ok we are going to scan a standalone thing.
				//We can chuck some optimisations later to deal with simple Type() and Type{} cases.
				
				c.Call(Embed)
				
			} else {
				statement(c, t, true)
			}
		}
		
		
		Absorb(c, t, false)
		c.RegisterType(t)
	},
}

func NotThing(t compiler.Type) bool {
	data := t.Data
	
	if data == nil {
		return true
	}
	
	if _, ok := data.(thing.Data); ok {
		return false
	}
	
	return true
}
