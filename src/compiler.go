package ilang

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
	Scope []Variables //This holds the current scope, it contains flags and variables.
	
	Output io.Writer //This is where the Compiler will output to.
	Lib io.Writer	//This is a Lib file which the compiler will write builtin functions to.
	
	Header bool 	//Are we still in the header of the file? or has a function been declared.
	Game bool		//Is this a game package?
	NewGame bool	//Does a NewGame constructor exist?
	UpdateGame bool //Does an update Game method exist?
	DrawGame bool	//Does a DrawGame method exist?
	
	Scanner *scanner.Scanner	//This is the current scanner.
	Scanners []*scanner.Scanner //Multiple files, multiple scanners.
	NextToken string			//You can overide the next token, this will be returned by the next call to scan.
	
	DefinedTypes map[string]Type			//A map of all the user created types.
	DefinedFunctions map[string]Function	//A map of all the user created functions.
	DefinedInterfaces map[string]Interface
	CurrentFunction Function				//If we are currently in a function, this is it.
	
	//Flags for compiling.
	SoftwareBlockExists bool //Does the software block exist?
	
	GUIExists bool		//Does a gui exist?
	GUIMainExists bool	//Does a main gui exist?
	
	Fork bool //Set this to true, in order to fork the next function call.
	
	//This is unused currently.
	InOperatorFunction bool
	
	LastDefinedType Type			//The latest user defined type.
	LastDefinedFunction Function	//The latest user defined function.
	LastDefinedFunctionName string  //The name of the latest user defined function.
	
	//Plugins
	Insertion []Plugin
	P int
	I int
	Lines int
	
	Plugins map[string][]Plugin
	
	//This is the expressiontype variable, it stores the type of the last scanned expression.
	ExpressionType Type
	
	//A counter for tmp variables so names do not clash.
	Unique int
	
	//These variables keep track of importing directories and packages.
	FileDepth int
	Dirs []string
	
	Stop bool //If this variable is set, the compiler will stop.
	
	Translation bool
	Translated bool
	Language string
}

//Return a string for a variable which will not clash with any other variables.
func (ic *Compiler) Tmp(mod string) string {
	ic.Unique++
	return "i_"+mod+fmt.Sprint(ic.Unique)
}

//This returns correctly formatted assembly.
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

//Assembly passed to this function will be output in the ilang.u library file.
func (ic *Compiler) Library(asm ...interface{}) {
	fmt.Fprintln(ic.Lib, ic.asm(asm...))
}

//Assembly passed to this function will be output to the file.
func (ic *Compiler) Assembly(asm ...interface{}) {
	fmt.Fprintln(ic.Output, ic.asm(asm...))
}


//This function increases the scope of the compiler for example when it reaches an if statement block.
func (c *Compiler) GainScope() {
	c.Scope = append(c.Scope, make(map[string]Type))
}

//Peek at the next character in the scanner.
func (c *Compiler) Peek() string {

	//Plugin injection.
	if len(c.Insertion) > 0 {
		if c.P != len(c.Insertion) {
			var text = c.Insertion[c.P].Tokens[c.I]
			return text
		}
	}

	return string(c.Scanner.Peek())
}

//Scan and return a token, can be called like:
/*
		var name = ic.Scan(Name) //Returns a name
		ic.Scan('(') 			//Expects a '(' char
		var token = ic.Scan(0) //Returns the string of the next token.
*/
//When an EOF is reached, Scan will stop the Compiler.
func (c *Compiler) Scan(verify rune) string {

	if c.NextToken != "" {
		var text = c.NextToken
		if verify > 0  && rune(text[0]) != verify {
			text = strconv.Quote(text)
			c.RaiseError("Unexpected "+text+", expecting "+string(verify))
		}
		c.NextToken = ""
		return text
	}

	//Plugin injection.
	if len(c.Insertion) > 0 {
		if c.P == len(c.Insertion) {
			c.P = 0
			c.I = 0
			c.Insertion = nil
			c.Lines = 0
		}
		
		if c.P != len(c.Insertion) {
			var text = c.Insertion[c.P].Tokens[c.I]
			if verify > 0  && rune(text[0]) != verify {
				text = strconv.Quote(text)
				c.RaiseError("Unexpected "+text+", expecting "+string(verify))
			}
			if text == "\n" {
				c.Lines++
			}
			
			if c.I == len(c.Insertion[c.P].Tokens)-1 {
				c.Lines = 0
				c.I = 0
				c.P++
			} else {
				c.I++
			}
		
			return text
		}
	}
	
	tok := c.Scanner.Scan()
	if verify > 0 && tok != verify {
		if verify > 9 {
			c.Expecting(string(verify))
		}
		c.RaiseError("Unexpected "+c.Scanner.TokenText())
	}
		
	if tok == scanner.EOF {
		if len(c.Scanners) > 0 {
			//var currentfile = c.Scanner.Filename
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
		
			//Create a software block.
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
			
			//Create a software block.
			if !c.SoftwareBlockExists && c.Game && !c.GUIExists {
				if !c.NewGame {
					c.Assembly("FUNCTION new_m_Game")
					c.GainScope()
					c.Assembly("ARRAY game")
					for range c.DefinedTypes["Game"].Detail.Elements {
						c.Assembly("PUT 0")
					}
					c.Assembly("SHARE game")
					c.LoseScope()
					c.Assembly("RETURN")
				}
					c.Assembly("FUNCTION Game")
					c.GainScope()
					c.Assembly("ARRAY game")
					for range c.DefinedTypes["Game"].Detail.Elements {
						c.Assembly("PUT 0")
					}
					c.Assembly("SHARE game")
					c.LoseScope()
					c.Assembly("RETURN")
				if !c.UpdateGame {
					c.Assembly("FUNCTION update_m_Game")
					c.Assembly("RETURN")
				}
				if !c.DrawGame {
					c.Assembly("FUNCTION draw_m_Game")
					c.Assembly("RETURN")
				}
			
				c.Assembly("SOFTWARE")
				c.GainScope()
				c.Assembly("RUN grate")
				c.LoseScope()
				c.Assembly("EXIT")
			}
			
			c.Stop = true
			return ""
		}
	}
	return c.Scanner.TokenText()
}

func (c *Compiler) Expecting(token string) {
	c.RaiseError("Expecting "+token)
}

func (ic *Compiler) LoseScope() {

	//Erm garbage collection???
	for name, variable := range ic.Scope[len(ic.Scope)-1] {
		if strings.Contains(name, "_") {
			var ok = false
			if ic.LastDefinedType.Detail != nil {
				_, ok = ic.LastDefinedType.Detail.Table[strings.Split(name, "_")[0]]
			}
			if variable == Unused && !(ic.GetFlag(InMethod) && ok ) {
				ic.RaiseError("unused variable! ", strings.Split(name, "_")[0])
			} 
		}
	
		if ic.Scope[len(ic.Scope)-1][name+"."] != Protected { //Protected variables
			if ic.GetFlag(InMethod) && name == ic.LastDefinedType.Name {
				continue
			}
			
			//Possible memory leak, TODO check up on this.
			if _, ok := ic.DefinedTypes[name]; ic.GetFlag(InMethod) && ok {
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

func (c *Compiler) TokenText() string {
	if len(c.Insertion) > 0 {
		return c.Insertion[c.P].Tokens[c.I-1]
	}
	return c.Scanner.TokenText()
}

func (c *Compiler) RaiseError(message ...interface{}) {
	pos := fmt.Sprint(c.Scanner.Pos())
	if len(c.Insertion) > 0 {
		pos = fmt.Sprintf("%v:%v:%v", c.Insertion[c.P].File, c.Insertion[c.P].Line+c.Lines, c.I) 
	}

	if len(message) == 0 {
		fmt.Fprintf(os.Stderr, "[%v] %v\n", pos, "Unexpected "+strconv.Quote(c.TokenText()))
	} else {
		fmt.Fprintf(os.Stderr, "[%v] %v\n", pos, fmt.Sprint(message...))
	}
	//panic("DEBUG TRACEBACK")
	os.Exit(1)
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
		DefinedInterfaces: make(map[string]Interface),
		DefinedTypes: make(map[string]Type),
		LastDefinedType: Something,
		Plugins: make(map[string][]Plugin),
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
		if ic.Stop {
			break
		}
		
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
				
				var braces = 0
				for {
					var token = ic.Scan(0)
					if strings.ContainsAny(token, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
						ic.GetVariable(token)
					}
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
					
					if block {
						if token == "}"  {
					 		if braces == 0 {
								break
							} else {
								braces--
							}
						}
						if token == "{" {
							braces++
						}
					}
				}
			case "plugin":
				ic.ScanPlugin()
			
			case "@":
				ic.Language = ic.Scan(Name)
				if ic.Language == "ch" {
					ic.Language = "zh-CN"
				}
				ic.Translation = true
				
			case "!":
				ic.Assembly("ADD ERROR 0 0")
				
			case "function":
				ic.Header = false
				ic.ScanFunction()
			
			case "method":
				ic.ScanMethod()
			
			case "interface":
				ic.ScanInterface()
			
			case "type":
				ic.Header = false
				ic.ScanType()
			
			case "gui":
				ic.Header = false
				ic.ScanGui()
				
			//case "new":
				//ic.ScanNew()
			
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
			
			//MESSY
			case "import":
				pkg := ic.Scan(Name)
				ic.Scan('\n')
				
				var filename = ""
				
				retry:
				file, err := os.Open(pkg+".i")
				if err != nil {
					if file, err = os.Open(pkg+"/"+pkg+".i"); err != nil {
						dir, _ := os.Getwd()
						if ic.FileDepth > 0 {
							ic.FileDepth--
							os.Chdir(ic.Dirs[len(ic.Dirs)-1])
							ic.Dirs = ic.Dirs[:len(ic.Dirs)-1]
							goto retry
						}
						
						ic.RaiseError("Cannot import "+pkg+", does not exist!", dir)
					} else {
						 filename = pkg+"/"+pkg+".i"
						 
						 dir, _ := os.Getwd()
						 
						ic.Dirs = append(ic.Dirs, dir)
						 
						os.Chdir("./"+pkg)
						ic.FileDepth++
					}
				} else {
					filename = pkg+".i"
				}
				ic.Scanners = append(ic.Scanners, ic.Scanner)
				
				ic.Scanner = &scanner.Scanner{}
				ic.Scanner.Init(file)
				ic.Scanner.Position.Filename = filename
				ic.Scanner.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
				
			case "software", "ソフトウェア", "программного", "软件":
				if token == "программного" {
					if ic.Scan(0) != "обеспечения" {
						ic.RaiseError("ожидая обеспечения")
					}
				}
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
			
			case "exit":
				if len(ic.Scope) > 2 {
					//TODO garbage collection.
					//ic.CollectGarbage()
					ic.Assembly("EXIT")
				}
			
			case "return":
				if !ic.CurrentFunction.Exists {
					ic.RaiseError("Cannot return, not in a function!")
				}
				
				if len(ic.CurrentFunction.Returns) == 1 {
					r := ic.ScanExpression()
					
					if ic.CurrentFunction.Returns[0] == User {
						ic.CurrentFunction.Returns[0] = ic.ExpressionType
					}
					
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
				ic.ScanSwitch()
			
			case "default":
				ic.ScanDefault()
			
			case "case":
				ic.ScanCase()
			
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
			
			case "delete":
				ic.Scan('(')
				var tok = ic.Scan(0)
				var arg string
				if tok != ")" {
					ic.NextToken = tok
					arg = ic.ScanExpression()
				}
				if ic.ExpressionType == Text {
					ic.Scan(')')
					ic.Assembly("SHARE ", arg)
					ic.Assembly("DELETE")
				} else if tok == ")" {
					if !ic.GetFlag(ForLoop) {
						ic.RaiseError("delete not in a for loop!")
					}
					//Delete things in a for loop.
					ic.Assembly("PLACE ", ic.GetVariable("i_for_delete").Name)
					ic.Assembly("PUT ", ic.GetVariable("i_for_id").Name)
				} else {
					ic.RaiseError("Invalid argument for delete.")
				}
				
			case "loop":
				ic.Assembly("LOOP")
				ic.GainScope()
				ic.NextToken = ic.Scan(0)
				if ic.NextToken != "{" {
					condition := ic.ScanExpression()
					ic.Assembly("SEQ ", condition, " 0 ", condition)
					ic.Assembly("IF ", condition)
					ic.Assembly("BREAK")
					ic.Assembly("END")
				}
				ic.Scan('{')
				ic.SetFlag(Loop)
			
			case "break":
				//TODO garbage collection.
				//ic.CollectGarbage()
				ic.Assembly("BREAK")
			
			case "for":
				ic.ScanForLoop()
			
			case "var", "ver", "变量":
				if len(ic.Scope) > 1 {
					name := ic.Scan(Name)
					token := ic.Scan(0)
					if token == "=" {
						ic.AssembleVar(name, ic.ScanExpression())
					} else if token == "is" {
						ic.AssembleVar(name, ic.ScanConstructor())
					} else if token == "has" {
						ic.AssembleVar(name, ic.ScanList())
					} else {
						ic.RaiseError("A variable should have a value assigned to it with an '=' sign.")
					}
				} else {
					ic.RaiseError("Global variables are not supported.")				
				}
		
			//This is the inbuilt print function. It takes multiple arguments of any type which has a text method.			
			case "print", "afdrukken", "印刷", "Распечатать", "打印":
				ic.Scan('(')
				arg := ic.ScanExpression()
				ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
				if ic.ExpressionType == Array {
					ic.LoadFunction("print_m_array")
					ic.LoadFunction("i_base_number")
					ic.Assembly("RUN print_m_array")
				} else {
					ic.Assembly(ic.RunFunction("text_m_"+ic.ExpressionType.Name))
					ic.Assembly("STDOUT")
				}
				
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
					if ic.ExpressionType == Array {
						ic.LoadFunction("print_m_array")
						ic.LoadFunction("i_base_number")
						ic.Assembly("RUN print_m_array")
					} else {
						ic.Assembly(ic.RunFunction("text_m_"+ic.ExpressionType.Name))
						ic.Assembly("STDOUT")
					}
				}
				
				ic.Assembly("SHARE i_newline")
				ic.Assembly("STDOUT")
			
			case "{":
				ic.Assembly("IF 1")
				ic.GainScope()
				ic.SetFlag(Block)
			
			case "if":
				var expression = ic.ScanExpression()
				if ic.ExpressionType != Number {
					ic.RaiseError("if statements must have numeric conditions!")
				}
				ic.Assembly("IF ", expression)
				ic.GainScope()
				
			case "else":
				nesting, ok := ic.Scope[len(ic.Scope)-1]["flag_nesting"]
				if !ok {
					nesting.Int = 0
				}
				ic.LoseScope()
				ic.Assembly("ELSE")
				ic.GainScope()
				ic.SetVariable("flag_nesting", Type{Int:nesting.Int})
			case "elseif":
				nesting, ok := ic.Scope[len(ic.Scope)-1]["flag_nesting"]
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
			
				nesting, ok := ic.Scope[len(ic.Scope)-1]["flag_nesting"]
				if ok {
					for i:=0; i < nesting.Int; i++ {
						ic.Assembly("END")
					}
				}
			
				var array = ic.GetVariable("i_for_array").Name
				var del = ic.GetVariable("i_for_delete").Name
			
				loopBefore := ic.GetScopedFlag(ForLoop)
				delBefore := ic.GetScopedFlag(Delete)
				ic.LoseScope()
				if loopBefore {
					ic.Assembly("REPEAT")
					
					if delBefore {
						ic.Assembly(`
	VAR ii_i8
	VAR ii_backup9
	LOOP
		VAR ii_in7
		ADD ii_i8 0 ii_backup9
		SGE ii_in7 ii_i8 #%v
		IF ii_in7
			BREAK
		END
		PLACE %v
		PUSH ii_i8
		GET i_v
		ADD ii_backup9 ii_i8 1

		VAR ii_operator11
		SUB ii_operator11 #%v 1
		PLACE %v
		PUSH ii_operator11
		GET ii_index12
		PLACE %v
		PUSH i_v
		SET ii_index12
		PLACE %v
		POP n
		ADD n 0 0
	REPEAT
						`, del, del, array, array, array, array)
					}
				}
				ic.Assembly("END")
				
			case "}":
			
				if ic.GetVariable("flag_switch") != Undefined {
					ic.LoseScope()
				}
					
				if ic.GetFlag(Issue) {
					ic.LoseScope()
				}
			
				nesting, ok := ic.Scope[len(ic.Scope)-1]["flag_nesting"]
				if ok {				
					for i:=0; i < nesting.Int+1; i++ {
						ic.Assembly("END")
					}
				}
			
				//TODO optimsise this with a scope saving checker object.
				softwarebefore := ic.GetFlag(Software)
				functionbefore := ic.GetFlag(InFunction)
				issuesbefore := ic.GetFlag(Issues)
				loopbefore := ic.GetScopedFlag(Loop)
				codeblock := ic.GetScopedFlag(Block)
				
				newbefore := ic.GetFlag(New)
				
				ic.LoseScope()
				
				newafter := ic.GetFlag(New)
				
				softwareafter := ic.GetFlag(Software)
				functionafter := ic.GetFlag(InFunction)
				issuesafter := ic.GetFlag(Issues)
				
				if softwarebefore != softwareafter {
					ic.Assembly("EXIT")
				}
				if newbefore != newafter {
					ic.Assembly("SHARE ", ic.LastDefinedType.Name)
				}
				if functionbefore != functionafter {
					if ic.InOperatorFunction {
						ic.InOperatorFunction = false
						ic.Assembly("SHARE c")
					}
					ic.Assembly("RETURN")
				}
				if issuesbefore != issuesafter || codeblock  {
					ic.Assembly("END")
				}
				if loopbefore {
					ic.Assembly("REPEAT")
				}
				
				
			default:
				
				//println(token)
				
				if t := ic.GetVariable(token); t != Undefined {
					switch t {
						case Number, Decimal, Letter:
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
									ic.ExpressionType = t
									ic.NextToken = token
									ic.Shunt(name)
									if ic.ExpressionType != Undefined {
										ic.RaiseError("blank expression!")
									}
									ic.ExpressionType = t
									if _, ok := ic.LastDefinedType.Detail.Table[name]; ic.GetFlag(InMethod) && ok {
										ic.SetUserType(ic.LastDefinedType.Name, name, name)
									}
							}
						case t.IsMatrix():
							var name = token
							token = ic.Scan(0)
							switch token {
								case "[":
									var x = ic.ScanExpression()
									ic.Scan(']')
									ic.Scan('[')
									var y = ic.ScanExpression()
									ic.Scan(']')
									ic.Scan('=')
									var value = ic.ScanExpression()
									
									ic.SetMatrix(name, x, y, value)
								
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
							}	
						
						case Array, Text, List, t.IsList(), t.IsArray():
							var name = token
							token = ic.Scan(0)
							
							switch token {
								case "&", "+":
									if !t.List && t != List {
										ic.ExpressionType = t
										ic.NextToken = token
										ic.Shunt(name)
										if ic.ExpressionType != Undefined {
											ic.RaiseError("blank expression!", ic.ExpressionType.Name)
										}
										continue
									}
									if token == "+" {
										ic.Scan('=')
									}
									
									value := ic.ScanExpression()
									
									if t == List && ic.ExpressionType.Push == "PUSH" {
										ic.SetVariable(name, Array)
										ic.Assembly("PLACE ", name)
										ic.Assembly("PUT ", value)
										continue
									}
									
									if t == List {
										list := ic.ExpressionType
										list.List = true
										list.User = false
										t = list
										ic.SetVariable(name, list)
										//println(name)
										if ic.GetFlag(InMethod) {
											ic.LastDefinedType.Detail.Elements[ic.LastDefinedType.Detail.Table[name]] = t
										}
									}
									
									//This appends elements to a list {..}
									if t.List {
										ic.PutList(t, name, value)
										
									} else {
									
										if ic.ExpressionType.Push != "PUSH" {
											ic.RaiseError("Only numeric values can be added to arrays.")
										}
										ic.Assembly("PLACE ", name)
										ic.Assembly("PUT ", value)
									}
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
								case "has":
									if ic.GetFlag(InMethod) {
										ic.SetUserType(ic.LastDefinedType.Name, name, ic.ScanList())
									} else {
										ic.AssembleVar(name, ic.ScanList())
									}
									
								default:
									ic.ExpressionType = t
									ic.NextToken = token
									ic.Shunt(name)
									if ic.ExpressionType != Undefined {
										ic.RaiseError("blank expression!")
									}
							}
						case Pipe:
							var name = token
							token = ic.Scan(0)
							switch token {
								case "(":
									argument := ic.ScanExpression()
									ic.Scan(')')
									if ic.ExpressionType != Text && ic.ExpressionType != Array {
										if ic.ExpressionType == Number {
											ic.Assembly("RELAY ", name)
											if argument != "" {
												ic.Assembly("PUSH ", argument)
											} else {
												ic.Assembly("PUSH 0")
											}
											ic.Assembly("IN")
											ic.Assembly("GRAB ", ic.Tmp("discard"))
											continue
										}
										ic.RaiseError("Only text and number values can be passed to a pipe call (outside of an expression).")
									}
									ic.Assembly("RELAY ", name)
									ic.Assembly("SHARE ", argument)
									ic.Assembly("OUT")
								case "=":
									value := ic.ScanExpression()
									if ic.ExpressionType != Pipe {
										ic.RaiseError("Only ",Func.Name," values can be assigned to ",name,".")
									}
									ic.Assembly("RELAY ", value)
									ic.Assembly("RELOAD ", name)
								default:
									ic.ExpressionType = t
									ic.NextToken = token
									ic.Shunt(name)
									if ic.ExpressionType != Undefined {
										ic.RaiseError("blank expression!")
									}
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
									ic.ExpressionType = t
									ic.NextToken = token
									ic.Shunt(name)
									if ic.ExpressionType != Undefined {
										ic.RaiseError("blank expression!")
									}
							}
						
						case t.IsSomething():
							var name = token
							ic.Scan('=')
							var value = ic.ScanExpression()
							ic.AssignSomething(name, value)
							
						case User:
							if !ic.GetFlag(InMethod) {
								ic.RaiseError()
							}
							var name = token
							ic.Scan('=')
							var value = ic.ScanExpression()
							ic.SetUserType(ic.LastDefinedType.Name, name, value)
						
						case t.IsUser():
							//Support indexing at any level
							// eg. Monster.Pos.X = 4
							var name = token
							var index string
							for token = ic.Scan(0); token == "."; {
								index = ic.Scan(Name)
								if token = ic.Scan(0); token == "." {
									name = ic.IndexUserType(name, index)
									ic.SetVariable(name, ic.ExpressionType) //This is required for setusertype to recognise.
									ic.SetVariable(name+".", Protected)
								}
							}
							var value string
							if token != "=" {
								value = ic.IndexUserType(name, index)
								
								var b = ic.ExpressionType
								ic.NextToken = token
								ic.Shunt(value)
								if ic.ExpressionType != Undefined {
									ic.RaiseError("blank expression!")
								}
								ic.ExpressionType = b
								
							} else {
								value = ic.ScanExpression()
							}
							
							
								//Set a usertype from within a method.
								if _, ok := ic.LastDefinedType.Detail.Table[name]; ic.GetFlag(InMethod) && ok {
									ic.SetUserType(ic.LastDefinedType.Name, name, value)
										
								} else if index == "" {
									//TODO garbage collection.
									ic.Assembly("PLACE ", value)
									ic.Assembly("RENAME ", name)
								} else {
									ic.SetUserType(name, index, value)
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
				
				} else if ic.GetFlag(InMethod) {
					ic.NextToken = token
					ic.ScanExpression()	
				
				} else {
					ic.RaiseError()
				}
		}
	}
}
