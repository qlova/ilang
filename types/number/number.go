package number

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/errors"
import "github.com/qlova/ilang/syntax/symbols"

import "math/big"

var Name = compiler.Translatable{
	compiler.English: "number",
	compiler.Maori: "tau",
}

var Type = compiler.Type {
	Name: Name,
	
	Base: compiler.INT,
}

func init() {
	Type.Shunts = Shunt
}

var Expression = compiler.Expression {
	
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		//Shunt here.
		if c.GetVariable(c.Token()).Type.Equals(Type) {
			var name = c.Token()
			
			c.Push(name)
			
			var t = c.GetVariable(name).Type
			return &t
		}
		
		switch c.Token()[0] {
			
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
				var b big.Int
				var worked bool
				
				if c.Token()[0] == '0' {
					_, worked = b.SetString(c.Token(), 2)
					
				} else {
				
					_, worked = b.SetString(c.Token(), 10)
				
				}
				
				if !worked {
					if len(c.Token()) > 2 {
						_, worked = b.SetString(c.Token()[2:], 16)
					}
				}
				
				if worked {
					c.BigInt(&b)
					
					return &Type
				} else {
					return nil
				}
				
			default:
				return nil
		}
	},
}

var Shunt = compiler.Shunt {
	symbols.Equals: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Same()
		
		return Type
	},
	
	symbols.Plus: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Add()
		
		return Type
	},
	
	symbols.Minus: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Sub()
		
		return Type
	},
	
	symbols.Times: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Mul()
		
		return Type
	},
	
	symbols.Divide: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Div()
		
		return Type
	},
	
	symbols.Power: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Pow()
		
		return Type
	},
	
	symbols.Modulus: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Mod()
		
		return Type
	},
	
	symbols.And: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Mul()
		
		return Type
	},
	
	symbols.Or: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Add()
		
		return Type
	},
	
	symbols.More: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.Less()
		
		return Type
	},
	
	symbols.Less: func (c *compiler.Compiler, t compiler.Type) compiler.Type {
		if !t.Equals(Type) {
			errors.Single(Type, symbols.Plus, t)
		}
		
		c.More()
		
		return Type
	},
}

var Method = compiler.Function {
	Name: Name,
	Returns: []compiler.Type{Type},
	
	Inline: func(c *compiler.Compiler) {
		c.Int(0)
	},
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		return statement(c, false)
	},
}

func init() {
	Type.EmbeddedStatement = func(c *compiler.Compiler, list compiler.Type) {
		statement(c, true)
	}
}
