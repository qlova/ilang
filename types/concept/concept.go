package concept

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"

type Data struct {
	Statement compiler.Statement
	Expression compiler.Expression
}

func (Data) Name(l compiler.Language) string {
	return ""
}

func (Data) Equals(d compiler.Data) bool {
	return false
}

var Type = compiler.Type {
	Name: compiler.Translatable{
		compiler.English: "concept",
	},
	
	Base: compiler.NULL,
}


var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
		if c.GetVariable(c.Token()).Equals(Type) {
			c.GetVariable(c.Token()).Type.Data.(Data).Statement.OnScan(c)
			return true
		}
		return false
	},
}

var Expression = compiler.Expression {
	
	Detect: func(c *compiler.Compiler) *compiler.Type {
		if c.GetVariable(c.Token()).Equals(Type) {
			
			if c.Peek() != symbols.FunctionCallBegin && c.Peek() != symbols.Index {
				var r = c.GetVariable(c.Token()).Type
				return &r
			}
			
			var r = c.GetVariable(c.Token()).Type.Data.(Data).Expression.OnScan(c)
			return &r
		}
		return nil
	},
}
