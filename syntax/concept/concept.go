package concept

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/ilang/syntax/errors"

import "github.com/qlova/ilang/types/list"
import "github.com/qlova/ilang/types/number"

import "github.com/qlova/uct/compiler"

//Format hacks.
import (
	"os"
	"runtime/debug"
	"fmt"
	"bytes"
	"strings"
)


var Name = compiler.Translatable{
	compiler.English: "concept",
}

//Urggy this is kinda hacky.
type Data struct {
	Children []*compiler.Function
}

func Overload(function *compiler.Function, args []compiler.Type) *compiler.Function {
	if len(function.Tokens) == 0 {
		return function
	}
	
	if function.Data == nil {
		function.Data = &Data{}
	}
	var data = function.Data.(*Data)

	function.Arguments = args
	
	here:
	for _, f := range data.Children {
		
		if len(f.Arguments) != len(args) {
			continue
		}
		
		for i := range f.Arguments {
			if !f.Arguments[i].Equals(args[i]) {
				continue here
			}
		}
		
		return f
	}
	
	var f compiler.Function
	f = *function
	f.Arguments = args
	f.Name = function.Name
	
	
	
	for i := range f.Name {
		
		var args_string string
		for j := range args {
			args_string += args[j].Name[i]
		}
		
		f.Name[i] = f.Name[i]+"_"+args_string
	}
	
	data.Children = append(data.Children, &f)
	
	return &f
}

var Return = compiler.Statement {
	Name: compiler.Translatable{
		compiler.English: "return",
	},
	
	OnScan: func(c *compiler.Compiler) {
		
		c.CollectAll()
		
		if c.Peek() == "\n" {
			c.Back()
			return
		}
		var t = c.ScanExpression()
		
		if len(c.CurrentFunction.Returns) == 0 {
			c.CurrentFunction.Returns = append(c.CurrentFunction.Returns, t)
		} else {
			if !c.CurrentFunction.Returns[len(c.CurrentFunction.Returns)-1].Equals(t) {
				c.RaiseError(errors.Inconsistent(t, c.CurrentFunction.Returns[len(c.CurrentFunction.Returns)-1]))
			}
		}
		if len(c.Scope) > 1 {
			c.Back()
		}
	},
}

func ScanCall(c *compiler.Compiler, f *compiler.Function) *compiler.Type {
	var args []compiler.Type
				
	c.Expecting(symbols.FunctionCallBegin)
	
	for i := 0; i < len(f.Tokens); i++ {
		if f.Variadic && i == len(f.Tokens)-1 {
			
			//Scan Variadic arguments!
			
			c.List()
			
			if c.Peek() == symbols.FunctionCallEnd {
				args = append(args, list.Type)
				break
			}
			
			var first = true
			for {
				variadic  := c.ScanExpression()
				if first {
					
					if !variadic.Equals(number.Type) {
						c.Unimplemented()
					}
					
					args = append(args, list.Of(variadic))
					first = false
				}
				
				//Type Checking!
				if !variadic.Equals(args[len(args)-1].Data.(*list.Data).SubType) {
					c.RaiseError(errors.Inconsistent(variadic, args[len(args)-1]))
				}
				
				c.Put()
				
				if c.Peek() == symbols.FunctionCallEnd {
					break
				}
				c.Expecting(symbols.ArgumentSeperator)
			}
			
			//Deal with arguments?
			break
		}
		args = append(args, c.ScanExpression())
		if i < len(f.Tokens)-1 {
			c.Expecting(symbols.ArgumentSeperator)
		}
	}
	
	c.Expecting(symbols.FunctionCallEnd)

	var overloaded = Overload(f, args)
	
	var old = c.CurrentFunction 
	//Hacky, way to get overloaded function modified by the return statement...
	c.CurrentFunction = overloaded
	
	c.Call(overloaded)
	c.CurrentFunction = old
	
	if len(overloaded.Returns) > 0 {
		return &overloaded.Returns[0]
	} else {
		return nil
	}
}

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		
		var name = c.Scan()
		if errors.IsInvalidName(name) {
			if !strings.Contains(name, "_dot_")  {
				c.RaiseError(errors.InvalidName(name))
			}
		}
		
		var f compiler.Function
		f.Name = compiler.Translatable{}
		f.Name[c.Language] = name
		
		c.Expecting(symbols.ArgumentListBegin)
		
		//Arguments...
		if c.Peek() != symbols.ArgumentListEnd {
			
			for {
				var arg = c.Scan()
				if errors.IsInvalidName(arg) {
					c.RaiseError(errors.InvalidName(arg))
				}
				
				f.Tokens = append(f.Tokens, arg)
			
				if c.Peek() != symbols.ArgumentListEnd {
					
					if c.Peek() == "." {
						c.Expecting(".")
						c.Expecting(".")
						c.Expecting(".")
						
						f.Variadic = true
						break
					}
					
					c.Expecting(symbols.ArgumentSeperator)
				} else {
					break
				}
			}
		}
		
		c.Expecting(symbols.ArgumentListEnd)

		c.Expecting(symbols.CodeBlockBegin)
		
		var filename = c.Scanners[len(c.Scanners)-1].Filename
		var line = c.Scanners[len(c.Scanners)-1].Line-1
		
		var cache = c.NewCache(symbols.CodeBlockBegin, symbols.CodeBlockEnd)
		//Need to be able to detect if cache is inlineable!
		
		
		
		f.Compile = func(c *compiler.Compiler) {
			
			//Pull arguments...
			if len(f.Arguments) > 0 {
				for i := len(f.Arguments)-1; i >= 0; i-- {
					c.PullType(f.Arguments[i], f.Tokens[i])
					c.SetVariable(f.Tokens[i], f.Arguments[i])
				}
			}
			
			if strings.Contains(name, "_dot_") {
				c.Pull("pointer")
			}
			
			var b bytes.Buffer
			c.StdErr = append(c.StdErr, &b)
			var linenumber = c.Scanners[len(c.Scanners)-1].Line
			
			defer func(c *compiler.Compiler, line int) {
				
				c.StdErr = c.StdErr[:len(c.StdErr)-1] // restoring the real stdout
				
				if r := recover(); r != nil {
					if r == "error" {
						c.Errors = true
						
						if os.Getenv("PANIC") == "1" { 
							panic("PANIC=1")
						}
						
						if len(f.Arguments) > 0 {
							c.Scanners[len(c.Scanners)-1].Line = line-c.LineOffset
							c.RaiseError(compiler.Translatable{
								compiler.English: "Cannot pass "+f.Arguments[0].Name[0]+" to "+f.Name[0]+
								" because\n\n"+b.String(),
							})
						}
						fmt.Println(b.String())
						panic("error")
						
					} else if r != "done" {
						fmt.Println(r, string(debug.Stack()))
					}
				}
			}(c, linenumber)

			

			c.CompileCache(cache, filename, line)
		}

		c.RegisterFunction(&f)
	},
	
	Detect: func(c *compiler.Compiler) bool {		
		for _, f := range c.Functions {
			if f.Name[c.Language] == c.Token() {
				
				ScanCall(c, f)
				
				//Deal with returns...
				
				return true
			}
		}
		
		return false
	},
}
