package array

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"

import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/types/list"

//Arrays have an immutable size.

var Type = compiler.Type {
	Name: compiler.Translatable{
		compiler.English: "array",
	},
	
	Base: compiler.LIST,
	
}

func init() {
	Type.Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) *compiler.Type {
		switch symbol {
			
			case symbols.IndexBegin:
				c.Expecting(symbols.IndexEnd)

				SubType(a).Base.Detach(c)
				
				var result = SubType(a)
				return &result

		}
		return nil
	}
}

func Is(t compiler.Type) bool {
	if t.Equals(Type) {
		return true
	}
	return false
}

func SubType(t compiler.Type) compiler.Type {
	return list.SubType(t)
}

func statement(c *compiler.Compiler, name string, t compiler.Type, embed bool) {
	switch c.Scan() {
		case symbols.Equals:
			
			c.ExpectingType(t)
			
			if embed {
				c.HeapList()
				c.Set()
			} else {
				c.NameList(name)
			}
			
		case symbols.IndexBegin:
			
			if embed {
				c.Unimplemented()
			} else {
				c.PushList(name)
			}
			
			c.ExpectingType(number.Type)
			c.Expecting(symbols.IndexEnd)

			c.ScanEmbeddedStatement(SubType(t))
			
		default:
			
			c.Unimplemented()
	}
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if c.GetVariable(c.Token()).Equals(Type) {
			statement(c, c.Token(), c.GetVariable(c.Token()).Type, false)
			return true
		}
		return false
	},
}

//Constructor for the array, array(), array(number), array(list)
var Expression = compiler.Expression {
	Name: Type.Name,
	
	Detect: func(c *compiler.Compiler) *compiler.Type {
		if c.GetVariable(c.Token()).Equals(Type) {
			var result = c.GetVariable(c.Token()).Type
			return &result
		}
		return nil
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		if c.Peek() == symbols.FunctionCallBegin {
			c.Scan()
			
			if c.Peek() == symbols.FunctionCallEnd {
				c.Scan()
				c.List()
				
				var result = Type
				result.Data = &list.Data{
					Step: 1,
					SubType: number.Type,
				}
				
				return result
			}
			
			var arg = c.ScanExpression()
			c.Expecting(")")
			
			if arg.Equals(number.Type) {
				
				c.Make()
				
				
				var result = Type
				result.Data = &list.Data{
					Step: 1,
					SubType: number.Type,
				}
				
				return result
			
			} else if list.Is(arg) {
				
				c.Call(&list.Copy)
				
				var result = Type
				result.Data = arg.Data
				
				return result
				
			} else {
				c.Unimplemented()
			}
				
		} else {
			c.Unimplemented()
		}
		return Type
	},
}
