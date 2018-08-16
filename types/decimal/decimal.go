package decimal

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/types/text"
import "github.com/qlova/ilang/syntax/symbols"

import "math/big"
import "fmt"
import "strings"

var Name = compiler.Translatable{
	compiler.English: "decimal",
	compiler.Maori: "tauāira",
}

var Type = compiler.Type {
	Name: Name,
	
	Base: compiler.INT,
}

var DefaultExponent = big.NewInt(1000000)

type Data struct {
	Exponent *big.Int
	Precision int64
}

func (d Data) Name(l compiler.Language) string {
	return fmt.Sprint(d.Precision)
}

func (self Data) Equals(d compiler.Data) bool {
	
	if (d.(Data).Precision == self.Precision) {
		return true
	}
	
	return false
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
		
		//Decimal numbers. TODO parsing errors.
		if strings.Contains(c.Token(), ".") {
		
			
		
			parts := strings.Split(c.Token(), ".")
			
			var precision = "1"+strings.Repeat("0", len(parts[1])) //Default
			
			var result = big.NewInt(10)
			result.Exp(result, big.NewInt(int64((len(precision)-1)-len(parts[1]))), nil)
			
			var p = big.NewInt(0)
			p.SetString(precision, 10)
			
			var big_a = big.NewInt(0)
			big_a.SetString(parts[0], 10)
			big_a.Mul(big_a, p)
			
			var big_b = big.NewInt(0)
			big_b.SetString(parts[1], 10)
			big_b.Mul(big_b, result)
			
			big_b.Add(big_a, big_b)
			
			c.BigInt(big_b)
			
			var t = Type.With(Data{
				Exponent: p,
				Precision: int64(len(parts[1])),
			})
			return &t
		}
		
		return nil
	},
}

var Method = compiler.Function {
	Name: Name,
	
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
	Type.Shunt = func(c *compiler.Compiler, symbol string, a, b compiler.Type) *compiler.Type  {
		

		if symbol == symbols.Divide {
			if a.Equals(b) {
				c.Swap()
				c.BigInt(a.Data.(Data).Exponent)
				c.Mul()
				c.Swap()
				c.Div()
				
				return &a
			}
		}
		
		if symbol == symbols.Less {
			if b.Equals(number.Type) {
				c.BigInt(a.Data.(Data).Exponent)
				c.Mul()
				b = a
			}
			
			if a.Equals(b) {
				c.More()
				
				return &a
			}
		}
		
		return nil
	}
	
	Type.Cast = func(c *compiler.Compiler, a, b compiler.Type) bool {
		if b.Equals(text.Type) {
			
			c.Int(10)
			c.BigInt(a.Data.(Data).Exponent)
			c.Int(a.Data.(Data).Precision)
			c.Call(&Dtoa)
			
			return true
		}
		return false			
	}
	
	Type.EmbeddedStatement = func(c *compiler.Compiler, list compiler.Type) {
		statement(c, true)
	}
	
	number.Type.Casts = append(number.Type.Casts, 
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(Type) {
				return true
			}
			return false			
		},
	)
	
	Method.Returns = []compiler.Type{Type.With(Data{Exponent: big.NewInt(1000000), Precision: 6})}
}
