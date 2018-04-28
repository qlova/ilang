package metatype

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"

import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/types/text"

import "strconv"

type Data struct {
	Type compiler.Type
}

func (Data) Name(l compiler.Language) string {
	return ""
}

func (Data) Equals(d compiler.Data) bool {
	return false
}

var Type = compiler.Type {
	Name: compiler.Translatable{
		compiler.English: "type",
	},
	
	//TODO allow types to be passed and compared as strings. This may open up powerful types of reflection.
	Base: compiler.NULL,
}

func init() {
	Type.Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) *compiler.Type {
		
		if symbol == symbols.Equals {
			
			if !a.Equals(b) {
				return nil
			}
			
			//TODO use fixed type so if statements using this expression can be optimised away.
			if a.Data.(Data).Type.Equals(b.Data.(Data).Type) {
				c.Int(1)
			} else {
				c.Int(0)
			}
			
			return &number.Type
		}
		
		return nil
	}
	
	Type.Cast = func(c *compiler.Compiler, a, b compiler.Type) bool {
		if b.Equals(text.Type) {
			
			c.SwapOutput()
			c.Data("text_literal"+strconv.Itoa(text.Tmp), []byte(a.Data.(Data).Type.String()))
			c.SwapOutput()
			
			c.PushList("text_literal"+strconv.Itoa(text.Tmp))
			
			text.Tmp++
			
			return true
		}
		return false
	}
}


var Expression = compiler.Expression {
	Name: Type.Name,
	
	Detect: func(c *compiler.Compiler) *compiler.Type {
		if c.GetType(c.Token()) != nil {
			if c.Peek() != symbols.FunctionCallBegin && c.Peek() != symbols.Index {
				var result = Type.With(Data{
					Type: *c.GetType(c.Token()),
				})
				return &result
			}
		}
		return nil
	},
	
	OnScan: func(c *compiler.Compiler) compiler.Type {
		c.Expecting(symbols.FunctionCallBegin)
		var arg = c.ScanExpression()
		c.Expecting(symbols.FunctionCallEnd)
		
		c.DropType(arg)
		
		return Type.With(Data{
			Type: arg,
		})
	},
}
