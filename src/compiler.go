package main

import "text/scanner"
import (
	"os"
	"fmt"
	"io"
	"strings"
	"strconv"
)

const (
	Name = scanner.Ident
	
)


//Deal with scoping for variables.
type Variables map[string]Type

type Compiler struct {
	Scope []Variables
	
	Output io.Writer
	Lib io.Writer
	
	Header bool
	
	Scanner *scanner.Scanner
	Scanners []*scanner.Scanner //Multiple files.
	NextToken string
	
	DefinedTypes map[string]Type
	DefinedFunctions map[string]Function
	CurrentFunction Function
	
	//Flags for compiling.
	SoftwareBlockExists bool
	
	GUIExists bool
	GUIMainExists bool
	
	InSwitchCase bool
	SwitchExpression string
	FirstCase bool
	
	InIssueBlock bool
	FirstIssue bool
	
	Fork bool
	
	InOperatorFunction bool
	
	LastDefinedType Type
	LastDefinedFunction Function
	LastDefinedFunctionName string
	
	ExpressionType Type
	
	Unique int
	
	InPackageDir bool
}

func (ic *Compiler) Tmp(mod string) string {
	ic.Unique++
	return "i_"+mod+fmt.Sprint(ic.Unique)
}

//This function increases the scope of the compiler for example when it reaches an if statement block.
func (ic *Compiler) asm(asm ...interface{}) string {
	var tabs = strings.Repeat("\t", len(ic.Scope)-1)
	var s = ""
	if len(asm) > 0 {
		if strings.Contains(fmt.Sprint(asm[0]), "%") && len(asm) > 1 {
			s = fmt.Sprintf(tabs+fmt.Sprint(asm[0]), asm[1:]...)
		} else {
			s = fmt.Sprint(tabs+fmt.Sprint(asm...))
		}
	}
	s = strings.Replace(s, "\n", "\n"+tabs, -1)
	return s
}

//This function increases the scope of the compiler for example when it reaches an if statement block.
func (ic *Compiler) Library(asm ...interface{}) {
	fmt.Fprintln(ic.Lib, ic.asm(asm...))
}

//This function increases the scope of the compiler for example when it reaches an if statement block.
func (ic *Compiler) Assembly(asm ...interface{}) {
	fmt.Fprintln(ic.Output, ic.asm(asm...))
}


//This function increases the scope of the compiler for example when it reaches an if statement block.
func (c *Compiler) GainScope() {
	c.Scope = append(c.Scope, make(map[string]Type))
}

//This will return the value of a scopped flag.
func (c *Compiler) GetFlag(sort Type) bool {
	for i:=len(c.Scope)-1; i>=0; i-- {
		if _, ok := c.Scope[i][sort.Name]; ok {
			return true
		}
	}
	return false
}

//This will return the type of the variable. UNDEFINED for undefined variables.
func (ic *Compiler) GetVariable(name string) Type {
	for i:=len(ic.Scope)-1; i>=0; i-- {
		if v, ok := ic.Scope[i][name]; ok {
			return v
		}
	}
	
	//Allow table values to be indexed in a method.
	if ic.GetFlag(InMethod) {
		if _, ok := ic.LastDefinedType.Detail.Table[name]; ok {
			var value = ic.IndexUserType(ic.LastDefinedType.Name, name)
			ic.AssembleVar(name, value)
			return ic.ExpressionType
		}
	}
	
	return Undefined
}

//Set the type of a variable, this is akin to creating or assigning a variable.
func (c *Compiler) SetVariable(name string, sort Type) {
	c.Scope[len(c.Scope)-1][name] = sort
}

//Set the type of a variable, this is akin to creating or assigning a variable.
func (c *Compiler) SetFlag(flag Type) {
	c.Scope[len(c.Scope)-1][flag.Name] = flag
}

func (ic *Compiler) Peek() string {
	return string(ic.Scanner.Peek())
}

func (c *Compiler) Scan(verify rune) string {
	if c.NextToken != "" {
		var text = c.NextToken
		if verify > 0  && rune(text[0]) != verify {
			text = strconv.Quote(text)
			c.RaiseError("Unexpected "+text+", expecting "+string(verify))
		}
		c.NextToken = ""
		return text
	} else {
		tok := c.Scanner.Scan()
		if verify > 0 && tok != verify {
			if verify > 9 {
				c.Expecting(string(verify))
			}
			c.RaiseError("Unexpected "+c.Scanner.TokenText())
		}
			
		if tok == scanner.EOF {
			if len(c.Scanners) > 0 {
				c.Scanner = c.Scanners[len(c.Scanners)-1]
				c.Scanners = c.Scanners[:len(c.Scanners)-1]
				
				return c.Scan(verify)
			} else {
				
				//Final cleanup and tasks.
				for _, t := range c.DefinedTypes {
					c.Collect(t)
				}
				
				fmt.Fprintf(c.Lib, `DATA i_newline "\n"`+"\n")
				c.LoadFunction("strings.equal")
				
				if !c.SoftwareBlockExists && c.GUIExists && c.GUIMainExists {
					c.Assembly("SOFTWARE")
					c.GainScope()
						c.Assembly("SHARE gui_main")
						c.Assembly("RUN gui")
						c.Assembly("LOOP")
						c.Assembly("SHARE i_newline")
						c.Assembly("RUN grab")
						c.Assembly("IF ERROR")
						c.GainScope()
							c.Assembly("EXIT")
						c.LoseScope()
						c.Assembly("END")
						c.Assembly("REPEAT")
					c.LoseScope()
					c.Assembly("EXIT")
					c.LoadFunction("grab")
					c.LoadFunction("gui")
					c.LoadFunction("output_m_pipe")
					c.LoadFunction("reada_m_pipe")
				}
				
				os.Exit(0)
			}
		}
		return c.Scanner.TokenText()
	}
}

func (c *Compiler) Expecting(token string) {
	if c.Scanner.TokenText() != token {
		c.RaiseError("Expecting "+token+" found "+strconv.Quote(c.Scanner.TokenText()))
	}
}

func (ic *Compiler) LoseScope() {

	//Erm garbage collection???
	for name, variable := range ic.Scope[len(ic.Scope)-1] {
		if ic.Scope[len(ic.Scope)-1][name+"."] != Protected { //Protected variables
			if ic.GetFlag(InMethod) && name == ic.LastDefinedType.Name {
				continue
			}
			
			if variable.IsUser() != Undefined {
				ic.Assembly("SHARE ", name)
				ic.Assembly("RUN collect_m_", variable.Name)
			}
		}
	}

	if len(ic.Scope) == 0 {
		ic.RaiseError()
	}
	ic.Scope = ic.Scope[:len(ic.Scope)-1]
	
}

func (c *Compiler) RaiseError(message ...interface{}) {
	if len(message) == 0 {
		fmt.Fprintf(os.Stderr, "[%v] %v\n", c.Scanner.Pos(), "Unexpected "+strconv.Quote(c.Scanner.TokenText()))
	} else {
		fmt.Fprintf(os.Stderr, "[%v] %v\n", c.Scanner.Pos(), fmt.Sprint(message...))
	}
	//panic("DEBUG TRACEBACK")
	os.Exit(1)
}

func (ic *Compiler) RunFunction(name string) string {
	f, ok := ic.DefinedFunctions[name]
	if !ok {
		ic.RaiseError(name, " does not exist!")
	}
	
	ic.LoadFunction(name)
	
	if f.Import != "" {
		ic.LoadFunction(f.Import)
	}
	
	if f.Inline {
		return f.Data
	} else if ic.Fork {
		return "FORK "+name
	} else {
		return "RUN "+name
	}
}

func (ic *Compiler) LoadFunction(name string) {
	f, ok := ic.DefinedFunctions[name]
	if !ok {
		ic.RaiseError(name, " does not exist!")
	}
	if !f.Inline && !f.Loaded {
		fmt.Fprintf(ic.Lib, f.Data)
		f.Loaded = true
		ic.DefinedFunctions[name] = f
	}
	if f.Import != "" {
		ic.LoadFunction(f.Import)
	}
}

func NewCompiler(input io.Reader) Compiler {
	var s scanner.Scanner
	s.Init(input)
	s.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
	
	c := Compiler{
		Scanner: &s,
		DefinedFunctions: make(map[string]Function),
		DefinedTypes: make(map[string]Type),
		LastDefinedType: Something,
	}
	
	c.Builtin()
	
	return c
}

//This will run through the compliation process.
func (ic *Compiler) Compile() {
	ic.GainScope()
	ic.SetVariable("error", Number)
	
	ic.Assembly(".import ilang")
	
	ic.Header = true
	
	for {
		token := ic.Scan(0)
		
		//These are all the tokens in ilang.
		switch token {
			case "\n", ";":
			
			//Inline assembly.
			case ".":
				cmd := ic.Scan(0)
				asm := strings.ToUpper(cmd)
				
				var data bool
				if cmd == "data" {
					data = true
				}
				
				//Are we in a block of code?
				var block = false
				
				var peeking = ic.Scan(0) 
				if peeking  == "{" {
					block = true
					ic.Scan('\n')
					asm = ""
				} else {
					ic.NextToken = peeking 
				}
				
				for {
					var token = ic.Scan(0)
					if data {
						ic.SetVariable(token, Text)
						data = false
					}
					if token == "\n" {
						if block {
							asm = strings.ToUpper(cmd)+" "+asm
						}
						if ic.Header {
							ic.Library(asm)
						} else {
							ic.Assembly(asm)
						}
						if !block {
							break
						} else {
							asm = ""
						}
					} else {
						if asm == "" {
							asm = token
						} else {
							asm += " "+token
						}
					}
					
					if block && token == "}" {
						break
					}
				}
				
			case "!":
				ic.Assembly("ADD ERROR 0 0")
				
			case "function":
				ic.Header = false
				ic.ScanFunction()
			
			case "method":
				ic.ScanMethod()
			
			case "type":
				ic.Header = false
				ic.ScanType()
			
			case "gui":
				ic.Header = false
				ic.ScanGui()
				
			case "new":
				ic.ScanNew()
			
			case "fork":
				name := ic.Scan(Name)
				ic.Scan('(')
				ic.Fork = true
				ic.ScanFunctionCall(name)
				ic.Scan(')')
			
			case "const":
				var name = ic.Scan(Name)
				ic.Scan('=')
				var value = ic.ScanExpression()
				if ic.ExpressionType.Push != "PUSH" {
					ic.RaiseError("Constant must be a numerical value! (",ic.ExpressionType.Name,")")
				} 
				ic.Assembly(".const %v %v", name, value)
				ic.SetVariable(name, ic.ExpressionType)
			
			case "import":
				pkg := ic.Scan(Name)
				ic.Scan('\n')
				
				file, err := os.Open(pkg+".i")
				if err != nil {
					if file, err = os.Open(pkg+"/"+pkg+".i"); err != nil {
						ic.RaiseError("Cannot import "+pkg+", does not exist!")
					} else {
						os.Chdir("./"+pkg)
						ic.InPackageDir = true
					}
				}
				ic.Scanners = append(ic.Scanners, ic.Scanner)
				
				ic.Scanner = &scanner.Scanner{}
				ic.Scanner.Init(file)
				ic.Scanner.Position.Filename = pkg+".i"
				ic.Scanner.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
				
			case "software":
				ic.Header = false
				ic.Scan('{')
				ic.Assembly("SOFTWARE")
				ic.GainScope()
				ic.SetFlag(Software)
				ic.SoftwareBlockExists = true
				
				if ic.GUIExists && ic.GUIMainExists {
					ic.Assembly("SHARE gui_main")
					ic.Assembly("RUN gui")
					ic.LoadFunction("gui")
					ic.LoadFunction("output_m_pipe")
				}
			
			case "return":
				if !ic.CurrentFunction.Exists {
					ic.RaiseError("Cannot return, not in a function!")
				}
				
				if len(ic.CurrentFunction.Returns) == 1 {
					r := ic.ScanExpression()
					if ic.ExpressionType != ic.CurrentFunction.Returns[0] {
						ic.RaiseError("Cannot return '",ic.ExpressionType.Name,
							"', expecting ",ic.CurrentFunction.Returns[0].Name)
					}
					
					ic.Assembly("%v %v", ic.ExpressionType.Push, r)
					
				}
				if len(ic.Scope) > 2 {
					//TODO garbage collection.
					//ic.CollectGarbage()
					ic.Assembly("RETURN")
				}
				
			case "switch":
				var expression = ic.ScanExpression()
				ic.Scan('{')
				ic.GainScope()
				ic.SetFlag(Type{Name: "flag_switch", Push: expression})
				
				for {
					token := ic.Scan(0)
					if token != "\n" {
						ic.Expecting("case")
						break
					}
				}
				expression = ic.ScanExpression()
				var condition = ic.Tmp("case")
				ic.Assembly("VAR ", condition)
				ic.Assembly("SEQ %v %v %v", condition, expression, ic.GetVariable("flag_switch").Push)
				ic.Assembly("IF ",condition)
				ic.GainScope()
			
			case "default":
				if ic.GetVariable("flag_switch") == Undefined {
					ic.RaiseError("'default' must be within a 'switch' block!")
				}
				ic.LoseScope()
				ic.Assembly("ELSE")
				ic.GainScope()
			
			case "case":
				if ic.GetVariable("flag_switch") == Undefined {
					ic.RaiseError("'case' must be within a 'switch' block!")
				}
			
				var expression = ic.ScanExpression()
				var condition = ic.Tmp("case")
				nesting, ok := ic.Scope[len(ic.Scope)-2]["flag_nesting"]
				if !ok {
					nesting.Int = 0
				}
				
				ic.LoseScope()
				
				ic.Assembly("ELSE")
				ic.SetVariable("flag_nesting", Type{Int:nesting.Int+1})
				
				ic.Assembly("VAR ", condition)
				ic.Assembly("SEQ %v %v %v", condition, expression, ic.GetVariable("flag_switch").Push)
				ic.Assembly("IF ",condition)
				ic.GainScope()
			
			case "issues":
				ic.Scan('{')
				ic.Assembly("IF ERROR")
				ic.GainScope()
				ic.Assembly("VAR issue")
				ic.Assembly("ADD issue ERROR 0")
				ic.Assembly("ADD ERROR 0 0")
				ic.SetFlag(Issues)
				
				var token string
				for {
					token = ic.Scan(0)
					if token != "\n" {
						if token != "issue" {
							ic.NextToken = token
						}
						break
					}
				}
				if token == "issue" {
				
					var expression = ic.ScanExpression()
					var condition = ic.Tmp("issue")
					ic.Assembly("VAR ", condition)
					ic.Assembly("SEQ %v %v %v", condition, expression, "issue")
					ic.Assembly("IF ",condition)
					ic.GainScope()
					ic.SetFlag(Issue)
				}
			
			case "issue":
				if !ic.GetFlag(Issues) {
					ic.RaiseError("'issue' must be within a 'issues' block!")
				}
			
				var expression = ic.ScanExpression()
				var condition = ic.Tmp("issue")
				nesting, ok := ic.Scope[len(ic.Scope)-2]["flag_nesting"]
				if !ok {
					nesting.Int = 0
				}
				
				ic.LoseScope()
				
				ic.Assembly("ELSE")
				ic.SetVariable("flag_nesting", Type{Int:nesting.Int+1})
				
				ic.Assembly("VAR ", condition)
				ic.Assembly("SEQ %v %v %v", condition, expression, "issue")
				ic.Assembly("IF ",condition)
				ic.GainScope()
				ic.SetFlag(Issue)
				
			case "loop":
				ic.Assembly("LOOP")
				ic.GainScope()
				ic.Scan('{')
				ic.SetFlag(Loop)
			
			case "break":
				//TODO garbage collection.
				//ic.CollectGarbage()
				ic.Assembly("BREAK")
			
			case "for":
				ic.ScanForLoop()
			
			case "var":
				name := ic.Scan(Name)
				token := ic.Scan(0)
				if token == "=" {
					ic.AssembleVar(name, ic.ScanExpression())
				} else if token == "is" {
					ic.AssembleVar(name, ic.ScanConstructor())
				} else {
					ic.RaiseError()
				}
			
			case "print":
				ic.Scan('(')
				arg := ic.ScanExpression()
				ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
				ic.Assembly(ic.RunFunction("text_m_"+ic.ExpressionType.Name))
				ic.Assembly("STDOUT")
				
				for {
					token := ic.Scan(0)
					if token != "," {
						if token != ")" {
							ic.RaiseError()
						}
						break
					}
					arg := ic.ScanExpression()
					ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
					ic.Assembly(ic.RunFunction("text_m_"+ic.ExpressionType.Name))
					ic.Assembly("STDOUT")
				}
				
				ic.Assembly("SHARE i_newline")
				ic.Assembly("STDOUT")
			
			case "if":
				var expression = ic.ScanExpression()
				ic.Assembly("IF ", expression)
				ic.GainScope()
				
			case "else":
				nesting, ok := ic.Scope[len(ic.Scope)-2]["flag_nesting"]
				if !ok {
					nesting.Int = 0
				}
				ic.LoseScope()
				ic.Assembly("ELSE")
				ic.GainScope()
				ic.SetVariable("flag_nesting", Type{Int:nesting.Int})
			case "elseif":
				nesting, ok := ic.Scope[len(ic.Scope)-2]["flag_nesting"]
				if !ok {
					nesting.Int = 0
				}
				ic.LoseScope()
				ic.Assembly("ELSE")
				var expression = ic.ScanExpression()
				ic.Assembly("IF ", expression)
				ic.GainScope()
				ic.SetVariable("flag_nesting", Type{Int:nesting.Int+1})
				
				
			
			case "end":
			
				nesting, ok := ic.Scope[len(ic.Scope)-2]["flag_nesting"]
				if ok {
					for i:=0; i < nesting.Int; i++ {
						ic.Assembly("END")
					}
				}
			
				loopBefore := ic.GetFlag(ForLoop)
				ic.LoseScope()
				loopAfter := ic.GetFlag(ForLoop)
				if loopBefore != loopAfter {
					ic.Assembly("REPEAT")
				} else {
					ic.Assembly("END")
				}
				
			case "}":
			
				nesting, ok := ic.Scope[len(ic.Scope)-2]["flag_nesting"]
				if ok {
				
					if ic.GetVariable("flag_switch") != Undefined {
						ic.LoseScope()
					}
					
					if ic.GetFlag(Issue) {
						ic.LoseScope()
					}
				
					for i:=0; i < nesting.Int+1; i++ {
						ic.Assembly("END")
					}
				}
			
				//TODO optimsise this with a scope saving checker object.
				softwarebefore := ic.GetFlag(Software)
				functionbefore := ic.GetFlag(InFunction)
				issuesbefore := ic.GetFlag(Issues)
				loopbefore := ic.GetFlag(Loop)
				
				newbefore := ic.GetFlag(New)
				
				ic.LoseScope()
				
				newafter := ic.GetFlag(New)
				
				softwareafter := ic.GetFlag(Software)
				functionafter := ic.GetFlag(InFunction)
				issuesafter := ic.GetFlag(Issues)
				loopafter := ic.GetFlag(Loop)
				
				if softwarebefore != softwareafter {
					ic.Assembly("EXIT")
				}
				if newbefore != newafter {
					ic.Assembly("SHARE ", ic.LastDefinedType.Name)
				}
				if functionbefore != functionafter {
					ic.Assembly("RETURN")
				}
				if issuesbefore != issuesafter {
					ic.Assembly("END")
				}
				if loopbefore != loopafter {
					ic.Assembly("REPEAT")
				}
				
			default:
				
				if t := ic.GetVariable(token); t != Undefined {
					switch t {
						case Number:
							var name = token
							if name == "error" {
								name = "ERROR"
							}
							token = ic.Scan(0)
							switch token {
								case "=":
									value := ic.ScanExpression()
									if ic.ExpressionType.Push != "PUSH" {
										ic.RaiseError("Only numeric values can assigned to ",name,".")
									}
									if _, ok := ic.LastDefinedType.Detail.Table[name]; ic.GetFlag(InMethod) && ok {
										ic.SetUserType(ic.LastDefinedType.Name, name, value)	
									} else {	
										ic.Assembly("ADD %v %v %v", name, 0, value)
									}
								default:
									ic.RaiseError()
							}
						case Array, Text:
							var name = token
							token = ic.Scan(0)
							switch token {
								case "&":
									value := ic.ScanExpression()
									if ic.ExpressionType.Push != "PUSH" {
										ic.RaiseError("Only numeric values can be added to arrays.")
									}
									ic.Assembly("PLACE ", name)
									ic.Assembly("PUT ", value)
								case "[":
									var index = ic.ScanExpression()
									ic.Scan(']')
									ic.Scan('=')
									var value = ic.ScanExpression()
									
									ic.Set(name, index, value)
								case "=":
									value := ic.ScanExpression()
									if ic.ExpressionType != t {
										ic.RaiseError("Only ",t.Name," values can be assigned to ",name,".")
									}
									
									if _, ok := ic.LastDefinedType.Detail.Table[name]; ic.GetFlag(InMethod) && ok {
										ic.SetUserType(ic.LastDefinedType.Name, name, value)	
									} else {									
										ic.Assembly("PLACE ", value)
										ic.Assembly("RENAME ", name)
									}
								default:
									ic.RaiseError()
							}
						case Func:
							var name = token
							token = ic.Scan(0)
							switch token {
								case "(":
									ic.Scan(')')
									ic.Assembly("EXE ", name)
								case "=":
									value := ic.ScanExpression()
									if ic.ExpressionType != Func {
										ic.RaiseError("Only ",Func.Name," values can be assigned to ",name,".")
									}
									ic.Assembly("PLACE ", value)
									ic.Assembly("RELOAD ", name)
								default:
									ic.RaiseError()
							}
						case User:
							if !ic.GetFlag(InMethod) {
								ic.RaiseError()
							}
							var name = token
							ic.Scan('=')
							var value = ic.ScanExpression()
							ic.SetUserType(ic.LastDefinedType.Name, name, value)	
						
						case t.IsUser():
							var name = token
							token = ic.Scan(0)
							if token == "." {
								var index = ic.Scan(Name)
								ic.Scan('=')
								var value = ic.ScanExpression()
								ic.SetUserType(name, index, value)
							} else {
								ic.RaiseError()
							}
							
						default:
							ic.RaiseError()
					}
				} else if _, ok := ic.DefinedFunctions[token]; ok {
					var check = ic.Scan(0)
					if check == "(" {
						ic.ScanFunctionCall(token)
						ic.Scan(')')
					} else if check == "@" {
						var variable = ic.expression()
						ic.Assembly("%v %v", ic.ExpressionType.Push, variable)
						ic.Scan('(')
						ic.ScanFunctionCall(token+"_m_"+ic.ExpressionType.Name)
						ic.Scan(')')
					} else {
						ic.RaiseError()
					}
				} else {
					ic.RaiseError()
				}
		}
	}
}
