package print

import "github.com/qlova/ilang/syntax/symbols"
import "github.com/qlova/uct/compiler"

import "github.com/qlova/ilang/types/text"
import "github.com/qlova/ilang/types/list"
import "github.com/qlova/ilang/types/array"
import "github.com/qlova/ilang/types/thing"
import "github.com/qlova/ilang/types/concept"
import "github.com/qlova/ilang/types/function"

import "fmt"
import "io/ioutil"

var Name = compiler.Translatable{
	compiler.English: "print",
	compiler.Maori: "perehitia",
}

func PrintType(c *compiler.Compiler, t compiler.Type) {
	
	if list.Is(t) || array.Is(t) {
		
		if thing.Is(list.SubType(t)) {
			c.Unimplemented()
		}
		
		c.CopyPipe()
		c.List()
		c.Int('[')
		c.Put()
		c.Send()
		
		c.Int(0)
		c.Loop()
			c.Copy()
			c.Size()
			c.Same() 
			c.If()
				c.Done()
			c.No()
		
			c.Copy()
			c.Get()
			
			if list.SubType(t).Base != compiler.INT && !list.SubType(t).Equals(text.Type) {
				c.Unimplemented()
			}
			
			//Add beginning quote.
			if  list.SubType(t).Equals(text.Type) {
				
				c.CopyPipe()
				c.List()
				c.Int('"')
				c.Put()
				c.Send()
				
				c.HeapList()
			}
			
			PrintType(c, list.SubType(t))
			
			if  list.SubType(t).Equals(text.Type) {
				c.CopyPipe()
				c.List()
				c.Int('"')
				c.Put()
				c.Send()
			}
			
			c.Copy()
			c.Size()
			c.Int(1)
			c.Sub()
			c.More()
			c.If()
				c.CopyPipe()
				c.List()
				c.Int(',')
				c.Put()
				c.Send()
			c.No()
			
			c.Int(1)
			c.Add()
		c.Redo()
		
		c.CopyPipe()
		c.List()
		c.Int(']')
		c.Put()
		c.Send()
		
		return
	}
	c.Cast(t, text.Type)
	
	c.CopyPipe()
	c.Send()
}


var Expression = compiler.Expression {
	Name: Name,
}

var Tmp int

func init() {
	Expression.OnScan = func(c *compiler.Compiler) compiler.Type {
		
		switch c.Peek() {
			
			//Testcase.
			case symbols.IndexBegin:
				c.Scan()
				
				var args []string
				var types  []*compiler.Type
				for {
					var v = c.Scan()
					if c.GetType(v) != nil {
						types = append(types, c.GetType(v))
						args = append(args, "")
					} else {
						types = append(types, nil)
						args = append(args, v)
					}
					if c.ScanIf(symbols.IndexEnd) {
						break
					}
					c.Expecting(symbols.ArgumentSeperator)
				}
				Tmp++
				
				
				//TODO cache this.
				
				var fargs []compiler.Type
				
				
				var out = c.Output 
				c.Output = ioutil.Discard
				
				var tmp = text.Tmp
				
				//Sift out any text literals. Swappy madness.
				for i := 0; i < len(args); i++ {
					if types[i] == nil {
							
							//DO some hacker level business here.
							var cache compiler.Cache
							cache.Write([]byte(args[i]))
							cache.Write([]byte{')'})
							
							c.LoadCache(cache, "print.go", 0)
							
							PrintType(c, c.ScanExpression())
							c.Expecting(")")
					}
				}
				text.Tmp = tmp
				
				c.SwapOutput()
				
				c.Code("print_"+fmt.Sprint(Tmp))
				
					c.List()
					c.Open()
				
					for i := 0; i < len(args); i++ {						
						if types[i] == nil {
							
							//DO some hacker level business here.
							var cache compiler.Cache
							cache.Write([]byte(args[i]))
							cache.Write([]byte{')'})
							
							c.LoadCache(cache, "print.go", 0)
							
							PrintType(c, c.ScanExpression())
							c.Expecting(")")
							
							
						} else {
							fargs = append(fargs, *types[i])
							PrintType(c, *types[i])
						}
					}
					
					c.List()
					c.Int('\n')
					c.Put()
					c.Send()
				
				c.Back()
				c.SwapOutput()
				
				c.Output = out
				
				c.Wrap("print_"+fmt.Sprint(Tmp))
				
				return function.Type.With(function.Data{
					Arguments: fargs,
				})
			
			default:
				return concept.Type.With(concept.Data{
					Statement: Statement,
					Expression: Expression,
				})
		}
		
		return compiler.Type{Fake: true}
	}
}

var Statement = compiler.Statement {
	Name: Name,
	 
	OnScan: func(c *compiler.Compiler) {
		c.Expecting(symbols.FunctionCallBegin)
		
		c.List()
		c.Open()
		
		if c.ScanIf(symbols.FunctionCallEnd) {
			//Print Newline
			c.List()
			c.Int('\n')
			c.Put()
			c.Send()
			return
		}
		
		for {
			PrintType(c, c.ScanExpression())
			
			switch c.Scan() {
				case symbols.ArgumentSeperator:
					continue
				case symbols.FunctionCallEnd:
					//Print Newline
					c.List()
					c.Int('\n')
					c.Put()
					c.Send()
					
					return
				default:
					c.Expected(symbols.ArgumentSeperator, symbols.FunctionCallEnd)
			}
		}
	},
}

