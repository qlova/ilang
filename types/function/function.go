package function

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/concept"

var Name = compiler.Translatable{
	compiler.English: "function",
}


type Data struct {
	Arguments []compiler.Type
}

func (Data) Name(l compiler.Language) string {
	return ""
}

func (Data) Equals(d compiler.Data) bool {
	return false
}

var Type = compiler.Type {
	Name: Name,
	
	Base: compiler.PIPE,
}

var Statement = compiler.Statement {
	Detect: func(c *compiler.Compiler) bool {
				
		//Shunt here.
		if c.GetVariable(c.Token()).Type.Equals(Type) {
			var name = c.Token()
			var function = c.GetVariable(c.Token()).Type
			
			c.Expecting(symbols.FunctionCallBegin)
			
			if function.Data != nil && len(function.Data.(Data).Arguments) > 0 {
				var data = function.Data.(Data)
				for i:=0; i < len(data.Arguments); i++ {
					c.ExpectingType(data.Arguments[i])
					if i < len(data.Arguments)-1 {
						c.Expecting(symbols.ArgumentSeperator)
					}
				}
			}
			
			if c.Peek() != symbols.FunctionCallEnd {
				c.Unimplemented()
			}
			
			c.Expecting(symbols.FunctionCallEnd)
			
			c.PushPipe(name)
			c.List()
			c.Send()

			return true
		}
		
		return false
	},
}	

var Expression = compiler.Expression {
	Detect: func(c *compiler.Compiler) *compiler.Type {
		
		//Shunt here.
		if c.GetVariable(c.Token()).Type.Equals(Type) {
			c.PushPipe(c.Token())
			
			var t = c.GetVariable(c.Token()).Type
			return &t
		}
		
		for _, f := range c.Functions {
			if f.Name[c.Language] == c.Token() {
				
				var name = c.Token()
				
				if c.Peek() == symbols.FunctionCallBegin {
					var r = concept.ScanCall(c, f)
					
					if r == nil {
						c.RaiseError(compiler.Translatable{
							compiler.English: "Cannnot use the function "+f.Name[compiler.English]+" in an expression, no return values!",
						})
					}
					
					return r
				}
				
				c.Wrap(name)
				return &Type
			}
		}
		
		return nil
	},
}
