package convert

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

var Name = compiler.Translatable {
	compiler.English: "convert",
}

var Expression = compiler.Expression {
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		var name = c.Token()
		
		var t = c.GetType(name)
		if t == nil {
			return nil
		}
		
		
		if f := c.GetFunction(name); f != nil && f.Inline != nil {
			c.Expecting("(")
			c.Expecting(")")
			f.Inline(c)
			return &f.Returns[0]
		} else {
			c.UndefinedError(name)
		}
		
		if c.GetVariable(name).Defined {
			return nil
		}
		
		if c.Peek() != symbols.FunctionCallBegin {
			return nil
		}
		
		c.Expecting(symbols.FunctionCallBegin)
		
		if c.Peek() == symbols.FunctionCallEnd {
			c.Scan()
			if f := c.GetFunction(name); f != nil && f.Inline != nil {
				f.Inline(c)
				return t
			} else {
				c.UndefinedError(name)
			}
		}
			
		var cast = c.ScanExpression()
		
		c.Cast(cast, *t)
		
		c.Expecting(symbols.FunctionCallEnd)
		
		return t
	},
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if t := Expression.Detect(c); t == nil {
			return false
		} else {
			c.DropType(*t)
			return true
		}
	},
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.Back()
	},
}


//This is currently tied to typer/type.go and the 'thing' type.
func New(c *compiler.Compiler, t *compiler.Type) {
	var name = c.Token()
	var b = c.GetType(name)
	if b == nil {
		c.RaiseError(errors.UnknownType(name))
	}
	
	var f compiler.Function
	c.CurrentFunction = &f
	
	var cast_name = t.Name[c.Language]+"_"+Name[c.Language]+"_"+b.Name[c.Language]
	
	c.Code(cast_name)
	c.Pull("pointer")
	c.CopyList()
	c.PullList("thing")
	
	c.Expecting(symbols.CodeBlockBegin)
	c.GainScope()
	c.SetFlag(Flag)
	for {
		c.ScanStatement()
		
		if _, ok := c.GetFlag(Flag); ok == -1  {
			break
		}
	}
	
	if len(f.Returns) == 0 {
		c.RaiseError(compiler.Translatable{
			compiler.English: "This conversion must return a "+b.Name[compiler.English]+" type!",
		})
	}
	if !f.Returns[0].Equals(*b) {
		c.RaiseError(compiler.Translatable{
			compiler.English: "This conversion must return a "+b.Name[compiler.English]+" type, not a "+f.Returns[0].Name[compiler.English]+" type.",
		})
	}
	
	t.Casts = append(t.Casts, 
		func(c *compiler.Compiler, t compiler.Type) bool {
			if t.Equals(*b) {
				
				c.CallRaw(cast_name)
				
				if t.Base == compiler.LIST {
					c.SwapList()
				}
				c.DropList()
				
				return true
			}
			return false			
		},
	)
}
