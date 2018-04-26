package content

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

import "github.com/qlova/ilang/syntax/convert"

var Name = compiler.Translatable {
	compiler.English: "content",
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
	
	var cast_name = b.Name[c.Language]+"_"+convert.Name[c.Language]+"_"+t.Name[c.Language]
	
	c.Code(cast_name)
		
		c.Pull("pointer")
		c.CopyList()
		c.PullList("thing")
		
				c.Pull(name)
	
	
	
	c.Expecting(symbols.CodeBlockBegin)
	c.GainScope()
	c.SetFlag(Flag)
	c.SetVariable(name, *b)
	for {
		c.ScanStatement()
		
		if _, ok := c.GetFlag(Flag); ok == -1  {
			break
		}
	}
	
	if len(f.Returns) != 0 {
		c.RaiseError(compiler.Translatable{
			compiler.English: "This should not return anything!",
		})
	}
}
