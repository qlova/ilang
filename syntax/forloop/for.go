package forloop

import "github.com/qlova/ilang/types/list"
import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/types/thing"
import "github.com/qlova/uct/compiler"

import "github.com/qlova/ilang/syntax/errors"

var Name = compiler.Translatable {
	compiler.English: "for",
}

var In = compiler.Translatable {
	compiler.English: "in",
}

var To = compiler.Translatable {
	compiler.English: "to",
}

var Each = compiler.Translatable {
	compiler.English: "each",
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		
		var flag = c.GetScope().Flags[Name[c.Language]]
		
		delete(c.GetScope().Flags, Name[c.Language])
		
		c.Push(GetIteratorName(c))
		c.Int(1)
		c.Add()
		c.Name(GetIteratorName(c))
		c.Redo()
		
		//Remove loop
		if flag.Bool {
			c.Int(0)
			c.Name(GetIteratorName(c))
			c.Loop()
				c.PushList(GetIteratorName(c)+"_remove")
				c.Size()
				c.Push(GetIteratorName(c))
				c.Same()
				c.Used()
				c.If()
					c.Done()
				c.No()
				
				c.PushList(GetIteratorName(c)+"_remove")
				c.Push(GetIteratorName(c))
				c.Get()
				c.Used()
				
				c.Size()
				c.Int(1)
				c.Sub()
				c.Get()
				
				c.Set()
				
				c.Pop()
				c.Drop()
				
				c.Push(GetIteratorName(c))
				c.Int(1)
				c.Add()
				c.Name(GetIteratorName(c))
			
			c.Redo()
		}
		
		if c.GetVariable(flag.Data).Defined {
			c.NameList(flag.Data)
		} else {
			c.Used()
		}
		
		c.No()
	},
}

var Remove = compiler.Statement {
	Name: compiler.Translatable {
		compiler.English: "remove",
	},
	
	OnScan: func(c *compiler.Compiler) {
		if flag, i := c.GetFlag(Flag); flag.Defined {
			c.DeleteFlag(flag)

			c.PushList(GetIteratorName(c)+"_remove")
			c.Push(GetIteratorName(c))
			c.Put()
			c.NameList(GetIteratorName(c)+"_remove")

			flag.Bool = true
			c.Scope[i].Flags[Name[c.Language]] = flag

		} else {
			c.RaiseError(compiler.Translatable {
				compiler.English: "You can only do this within a for loop!",
			})
		}
	},
}

var OverFlag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		delete(c.GetScope().Flags, Name[c.Language])
		
		c.Copy()
		c.Push(GetIteratorName(c))
		c.Less()
			c.If()
			
				c.Push(GetIteratorName(c))
				c.Int(1)
				c.Add()
				c.Name(GetIteratorName(c))
				
			c.Or()
				
				c.Copy()	
				c.Push(GetIteratorName(c))
				c.More()
				
				c.If()
					c.Push(GetIteratorName(c))
					c.Int(1)
					c.Sub()
					c.Name(GetIteratorName(c))
				c.Or()
					c.Done()
				c.No()
			c.No()
		c.Redo()
		c.Drop()
		c.No()
	},
}


var End = compiler.Statement {
	Name: compiler.Translatable {
		compiler.English: "end",
	},
	 
	OnScan: func(c *compiler.Compiler) {
		c.LoseScope()
	},
}

func GetIteratorName(c *compiler.Compiler) string {
	var depth = 0
	for _, scope := range c.Scope {
		if scope.Flags[Name[c.Language]].Defined {
			depth ++
		}
	}
	switch depth {
		case 0:
			return "i"
		
		case 1:
			return "j"
		
		case 2:
			return "k"
		
		case 3:
			return "l"
		
		default:
			c.RaiseError(compiler.Translatable{
				compiler.English: "Too many nested for loops! Try putting this one inside a function...",
			})
	}
	return ""
}

func ScanIn(c *compiler.Compiler) {
	var value = c.Token()
	if errors.IsInvalidName(value) {
		c.RaiseError(errors.InvalidName(value))
	}
	
	c.Expecting(In[c.Language])
	
	var variable = c.Peek()
	var t = c.ScanExpression()
	
	if !t.Equals(list.Type) {
		c.RaiseError(compiler.Translatable{
			compiler.English: "Cannot iterate over a "+t.Name[compiler.English]+" type!",
		})
	}
	
	//Proper step through array.
	var Step = 1
	if t.Data != nil && value != "each" {
		Step = t.Data.(*list.Data).Step
		if Step == 0 {
			Step = 1
		}
	}

	c.Int(1)
	c.If()
	c.List()
	c.PullList(GetIteratorName(c)+"_remove")
	c.Int(0)
	c.Pull(GetIteratorName(c))
	c.Loop()
		c.Size()
		c.Push(GetIteratorName(c))
		c.Same()
		c.If()
			c.Done()
		c.No()
		
		c.GainScope()
		c.SetVariable(GetIteratorName(c), number.Type)
		
		if value != "each" {
			c.Push(GetIteratorName(c))
			
			if t.Data == nil {
				c.Unimplemented()
			} else if !thing.NotThing(t.Data.(*list.Data).SubType) {
				
				c.Pull(value+"_pointer")
				c.CopyList()
				c.PullList(value)
				
				c.SetVariable(value, t.Data.(*list.Data).SubType)
				
			} else if t.Data.(*list.Data).SubType.Base == compiler.INT {
				
				c.Get()
				c.Pull(value)
				c.SetVariable(value, t.Data.(*list.Data).SubType)
				
				
			} else {
				c.Unimplemented()
			}
		}
		
		var f = Flag
		f.Value = 1
		f.Data = variable
		
		c.SetFlag(f)
		
		
}

var Statement = compiler.Statement {
	Name: Name,
	
	OnScan: func(c *compiler.Compiler) {
		switch c.Scan() {
			case Each[c.Language]:
				
				//for each 0 to 10
				if c.Peek() != In[c.Language] {
					c.Int(1)
					c.If()
					
					var start = c.ScanExpression()
					if !start.Equals(number.Type) {
						c.RaiseError(errors.ExpectingType(number.Type, start))
					}
					
					c.Pull(GetIteratorName(c))
					
					c.Expecting(To[c.Language])
					
					var end = c.ScanExpression()
					if !end.Equals(number.Type) {
						c.RaiseError(errors.ExpectingType(number.Type, end))
					}
					
					c.List()
					c.PullList(GetIteratorName(c)+"_remove")
					c.Loop()
						c.GainScope()
						c.SetVariable(GetIteratorName(c), number.Type)
						c.SetFlag(OverFlag)
					
					return
				} else {
				
					ScanIn(c)
				}
				
			default:

				ScanIn(c)
		}
	},
}
