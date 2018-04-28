package text

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/letter"
import "github.com/qlova/ilang/types/number"

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

import "strconv"

var Type = compiler.Type {
	Name: compiler.Translatable{
		compiler.English: "text",
	},
	
	Base: compiler.LIST,
}

var Shunts = compiler.Shunt {
	symbols.Times: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(number.Type) {
			c.RaiseError(errors.Single(Type, symbols.Plus,t))
		}
		
		c.Size()
		c.Mul()
		c.Flip()
		c.List()
		c.Loop()
			c.Copy()
			c.Int(0)
			c.Same()
			c.If()
				c.Done()
			c.No()
		
			c.Copy()
			c.SwapList()
			c.Size()
			c.Mod()
			c.Get()
			c.SwapList()
			
			c.Put()
		
			c.Int(1)
			c.Add()
		c.Redo()
		
		c.SwapList()
		c.DropList()
		
		return Type
	},
	
	symbols.Plus: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			c.RaiseError(errors.Single(Type, symbols.Plus,t))
		}
		
		c.Call(&Join)
		
		return Type
	},
}

//Could optimise this if there are two of the same string.
var Tmp = 0

func Shunt(c *compiler.Compiler) *compiler.Type {

	switch c.Token() {
		case symbols.IndexBegin:
			
			if !c.ScanExpression().Equals(number.Type) {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Text can only be indexed with the number type.",
				})
			}
			
			c.Expecting(symbols.IndexEnd)
			
			
			c.Size()
			c.Mod()
			c.Get()
			c.Used()
			
			return &letter.Type
		
		case symbols.Plus:
			
			if t := c.ScanExpression(); !t.Equals(Type) {
				c.RaiseError(errors.Single(Type, symbols.Plus,t))
			}
			
			c.Call(&Join)
			
			return &Type
		case symbols.SelectMethod:			
			if c.Scan() == "size" {

				c.Expecting(symbols.FunctionCallBegin)
				c.Expecting(symbols.FunctionCallEnd)

				c.Size()
				
				return &number.Type
				
			} else {
				
				c.Unimplemented()
			}
	}
	return nil
}

var Expression = compiler.Expression {
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		if c.GetVariable(c.Token()).Type.Equals(Type) {
			
			c.PushList(c.Token())
			
			switch c.Peek() {
				case symbols.SelectMethod, symbols.IndexBegin:
					c.Scan()

					return Shunt(c) 
			} 
			
			return &Type
		}
		
		//TODO this needs to be copied if mutated. Data is not really supposed to be mutated!
		if c.Token()[0] == '"' {
			
			text, err := strconv.Unquote(c.Token())
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Invalid String!",
				})
			}
			
			c.SwapOutput()
			c.Data("text_literal"+strconv.Itoa(Tmp), []byte(text))
			c.SwapOutput()
			
			c.PushList("text_literal"+strconv.Itoa(Tmp))
			
			Tmp++

			var t = Type
			t.Immutable = true
			
			
			
			switch c.Peek() {
				case symbols.SelectMethod, symbols.IndexBegin:
					c.Scan()

					return Shunt(c) 
			} 
			return &t
		} else {
			return nil
		}
	},
}

func init() {
	Type.Shunts = Shunts
	
	Type.Casts = append(Type.Casts, 
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(number.Type) {
				
				c.Int(10)
				c.Call(&Atoi)
				
				return true
			}
			return false			
		},
	)

	number.Type.Casts = append(number.Type.Casts, 
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(Type) {
				
				c.Int(10)
				c.Call(&Itoa)
				
				return true
			}
			return false			
		},
	)

	letter.Type.Casts = append(letter.Type.Casts, 
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(Type) {
				c.List()
				c.Put()
				return true
			}
			return false			
		},
	)
}

