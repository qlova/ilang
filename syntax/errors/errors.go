package errors

import "github.com/qlova/uct/compiler"
import "github.com/qlova/ilang/syntax/symbols"

import "strings"

type Base struct {}
func (Base) Push(c *compiler.Compiler, data string) {}
func (Base) Pull(c *compiler.Compiler, data string) {}
func (Base) Drop(c *compiler.Compiler) {}
func (Base) Attach(c *compiler.Compiler) {}
func (Base) Detach(c *compiler.Compiler) {}

func (Base) Free(c *compiler.Compiler) {
	c.Int(0)
	c.Name("ERROR")
}

var Name = compiler.Translatable{
	compiler.English: "errors",
}

var Statement = compiler.Statement{
	Name: Name,

	OnScan: func(c *compiler.Compiler) {
		c.Expecting(symbols.CodeBlockBegin)
		
		c.Push("ERROR")
		c.If()
		
		c.GainScope()
		c.SetFlag(Flag)
		
		if !c.GetVariable("error").Defined {
			c.SetVariable("error", compiler.Type{
				Base: Base{},
			})
		}
	},
}

var Flag = compiler.Flag {
	Name: Name,
	
	OnLost: func(c *compiler.Compiler) {
		c.No()
	},
}

var End = compiler.Statement {
	Name: compiler.NoTranslation(symbols.CodeBlockEnd),
	 
	OnScan: func(c *compiler.Compiler) {
		c.LoseScope()
	},
}

func AssignmentMismatch(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Cannot assign a "+a.String()+" value to a variable of type "+b.String(),
	}
}

func ExpectingType(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Expecting a value of type "+a.String()+" instead got a value of type "+b.String(),
	}
}

func UnknownType(a string) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Unknown Type: "+a,
	}
}

func MustBeNumeric(a compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: a.String()+" is not a numeric type!",
	}
}

func Single(a compiler.Type, symbol string, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "The relationship "+a.String()+symbol+b.String()+" is not defined!",
	}
}

func IsInvalidName(name string) bool {
	return strings.Contains(name, "_")
}


func InvalidName(name string) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Invalid name: "+name,
	}
}

func Inconsistent(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "The usage here of the '"+a.String()+"' type is inconsistent with the\n use of the '"+b.String()+"' type before this!", 
	}
}

func NoSuchElement(a string, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "No such element '"+a+"' in type '"+b.String()+"'!", 
	}
}
