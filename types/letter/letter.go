package letter

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/types/number"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"
import "strconv"

var Type = compiler.Type {
	Name: compiler.Translatable{
		compiler.English: "letter",
	},
	
	Base: compiler.INT,
	
	Casts: []func(*compiler.Compiler, compiler.Type) bool {
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(number.Type) {
				return true
			}
			return false
		},
	},
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
		
		if c.Token()[0] == '\'' {
			
			text, err := strconv.Unquote(c.Token())
			if err != nil {
				c.RaiseError(compiler.Translatable{
					compiler.English: "Invalid Letter!",
				})
			}
			
			c.Int(int64(text[0]))
			
			return &Type
		} else {
			return nil
		}
	},
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if c.GetVariable(c.Token()).Type.Equals(Type) {
			
			var name = c.Token()
			
			switch c.Scan() {
				
				case symbols.Equals:
					var t = c.ScanExpression()
					if !t.Equals(Type) {
						c.RaiseError(errors.AssignmentMismatch(t, Type))
					}
					
					c.Name(name)
				
				default:
					c.Unexpected(name+c.Token())
			}
			
			return true
		}
		return false
	},
}

func init() {
	number.Type.Casts = append(number.Type.Casts, 
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(Type) {
				return true
			}
			return false			
		},
	)
}
