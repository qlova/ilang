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
type Scope map[string]Type

type Compiler struct {
	Scope []Scope //This holds the current scope, it contains flags and variables.
	SwapScope []Scope
	
	Output io.Writer //This is where the Compiler will output to.
	Lib io.Writer	//This is a Lib file which the compiler will write builtin functions to.
	
	Header bool 	//Are we still in the header of the file? or has a function been declared.
	Game bool		//Is this a game package?
	NewGame bool	//Does a NewGame constructor exist?
	UpdateGame bool //Does an update Game method exist?
	DrawGame bool	//Does a DrawGame method exist?
	
	Scanner *scanner.Scanner	//This is the current scanner.
	Scanners []*scanner.Scanner //Multiple files, multiple scanners.
	
	CurrentLine string //All the tokens up to this point.
	CurrentLineReset bool
	
	LastToken string
	NextToken string			//You can overide the next token, this will be returned by the next call to scan.
	NextNextToken string
	NextNextNextToken string
	
	DefinedTypes map[string]Type			//A map of all the user created types.
	DefinedFunctions map[string]Function	//A map of all the user created functions.
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
	DisableOwnership bool //Disable taking ownership.
	TakingExpression bool
	ProtectExpression bool //Protect expression values from being garbage collected
	
	//A counter for tmp variables so names do not clash.
	Unique int
	
	//These variables keep track of importing directories and packages.
	FileDepth int
	Dirs []string
	
	Stop bool //If this variable is set, the compiler will stop.
	
	Translation bool
	Translated bool
	Language string
	
	
	//Set variables. (Sets.i)
	SetItemCount int
	SetItems map[string]int
	
	//Optimisation
	LastLine string
	ProgramDir string
	
	Aliases map[string]string
	
	DisableOutput bool
}

func (ic *Compiler) SwapOutput() {
	ic.Assembly("\n")
	ic.Output, ic.Lib = ic.Lib, ic.Output
	ic.Scope, ic.SwapScope = ic.SwapScope, ic.Scope
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
	if !ic.DisableOutput {
		var raw = ic.asm(asm...)
		ic.optimise(raw)
	}
}

func (ic *Compiler) optimise(asm string) {
	var lines = strings.Split(asm, "\n")
	for _, line := range lines {
		//line = strings.TrimSpace(line)
		
		var a, b = strings.TrimSpace(line), strings.TrimSpace(ic.LastLine)
		
		if strings.Contains(a, "_") &&  strings.Contains(b, "_") && len(a) > 4 && len(b) > 4 && a[:4] == "PUSH" && b[:4] == "PULL" && a[4:] == b[4:] {
			ic.LastLine = "#opt"
		} else {
			if ic.LastLine != "#opt" {
				fmt.Fprintln(ic.Output, ic.LastLine)
			}
			ic.LastLine = line	
		}
	}
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
	var r = c.scan(verify)
	if a, ok := c.Aliases[r]; ok {
		return a
	}
	return r
}

func (c *Compiler) scan(verify rune) string {
	if c.NextToken != "" {
		var text = c.NextToken
		if verify > 0  && rune(text[0]) != verify {
			text = strconv.Quote(text)
			c.RaiseError("Unexpected "+text+", expecting "+string(verify))
		}
		c.NextToken = c.NextNextToken
		c.NextNextToken = c.NextNextNextToken
		c.NextNextNextToken = ""
		
		c.LastToken = text
		return text
	}

	//Plugin injection.
	if len(c.Insertion) > 0 {
		if c.P == len(c.Insertion) {
			c.P = 0
			c.I = 0
			c.Insertion = nil
			c.Lines = 0
			return c.scan(verify)
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
		
			c.LastToken = text
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
			if c.Scanner.Filename == "" {
				os.Chdir(c.ProgramDir)
			}
			
			return c.Scan(verify)
		} else if !c.Stop {
			
			//Final cleanup and tasks.
			for _, t := range c.DefinedTypes {
				c.Collect(t)
			}
			
			fmt.Fprintf(c.Lib, `DATA i_newline "\n"`+"\n")
			c.LoadFunction("strings.equal")
			c.LoadFunction("strings.compare")
			c.LoadFunction("i_base_number")
		
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
					c.Assembly("FUNCTION new_m_Graphics")
					c.GainScope()
					c.Assembly("ARRAY game")
					for range c.DefinedTypes["Graphics"].Detail.Elements {
						c.Assembly("PUT 0")
					}
					c.Assembly("SHARE game")
					c.LoseScope()
					c.Assembly("RETURN")
				}
					c.Assembly("FUNCTION Graphics")
					c.GainScope()
					c.Assembly("ARRAY game")
					for range c.DefinedTypes["Graphics"].Detail.Elements {
						c.Assembly("PUT 0")
					}
					c.Assembly("SHARE game")
					c.LoseScope()
					c.Assembly("RETURN")
				if !c.UpdateGame {
					c.Assembly("FUNCTION update_m_Graphics")
					c.Assembly("RETURN")
				}
				if !c.DrawGame {
					c.Assembly("FUNCTION draw_m_Graphics")
					c.Assembly("RETURN")
				}
			
				c.Assembly("SOFTWARE")
				c.GainScope()
				c.Assembly("RUN grate")
				c.LoseScope()
				c.Assembly("EXIT")
			}
			
			c.Stop = true
			c.LastToken = ""
			return ""
		}
	}
	c.LastToken = c.Scanner.TokenText()
	
	
	if c.CurrentLineReset {
		c.CurrentLine = ""
		c.CurrentLineReset = false
	}
	if len(c.LastToken) > 1 {
		c.CurrentLine += " "+c.LastToken
	} else {
		c.CurrentLine += c.LastToken
	}
	if c.LastToken == "\n" {
		c.CurrentLineReset = true
	}
	
	return c.Scanner.TokenText()
}

func (c *Compiler) Expecting(token string) {
	c.RaiseError("Expecting "+token)
}

func (ic *Compiler) LoseScope() {

	ic.CollectGarbage()

	if len(ic.Scope) == 0 {
		ic.RaiseError()
	}
	
	//Prep our listeners.
	//This allows modules to listen when a flag has fallen out of scope.
	var cache = make(map[Type]bool)
	for listener := range Listeners {
		cache[listener] = ic.GetScopedFlag(listener)
	}			
	
	ic.Scope = ic.Scope[:len(ic.Scope)-1]
	
	for listener, f := range Listeners {
		if cache[listener] {
			f(ic)
		}
	}
	
}

func (c *Compiler) TokenText() string {
	if len(c.Insertion) > 0 && c.P < len(c.Insertion) {
		return c.Insertion[c.P].Tokens[c.I-1]
	}
	return c.Scanner.TokenText()
}

func (c *Compiler) RaiseError(message ...interface{}) {
	
	//Let the user know what line it is!
	fmt.Print(c.Scanner.Pos().Line, ": ", c.CurrentLine)
	c.NextToken = ""
	char := len(c.CurrentLine)
	if char == 0 {
		char = 1
	}
	if !c.CurrentLineReset {
		for {
			tok := c.Scan(0)
			fmt.Print(tok)
			if len(tok) > 0 {
				fmt.Print(" ")
			}
			if tok == "\n" {
				break
			}
		}
	} else {
		char--
	}
	fmt.Print(strings.Repeat(" ", len(fmt.Sprint(c.Scanner.Pos().Line))+1), " ", strings.Repeat(" ", char-1))
	fmt.Println("^")

	if len(message) == 0 {
		fmt.Fprintf(os.Stderr, "    %v\n", "Unexpected "+strconv.Quote(c.TokenText()))
	} else {
		fmt.Fprintf(os.Stderr, "    %v\n", fmt.Sprint(message...))
	}
	if os.Getenv("PANIC") == "1" { 
		panic("PANIC=1")
	}
	os.Exit(1)
}

func NewCompiler(input io.Reader) Compiler {
	var s scanner.Scanner
	s.Init(input)
	s.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
	
	c := Compiler{
		Scanner: &s,
		DefinedFunctions: make(map[string]Function),
		DefinedTypes: make(map[string]Type),
		LastDefinedType: Undefined,
		Plugins: make(map[string][]Plugin),
		SetItems: make(map[string]int),
		Aliases: make(map[string]string),
	}
	
	c.Builtin()
	
	for name, f := range Functions {
		c.DefinedFunctions[name] = f
	}
	
	return c
}

//This will run through the compliation process.
func (ic *Compiler) Compile() {
	ic.GainScope()
	ic.SetVariable("error", Number)
	
	ic.Assembly(".import ilang")
	
	ic.Header = true
	
	ic.ProgramDir, _ = os.Getwd()
	
	ic.SwapScope = append(ic.SwapScope, ic.Scope[0])
	
	for ic.ScanAndCompile() {
		
	}
}

func (ic *Compiler) ScanAndCompile() bool {
	token := ic.Scan(0)
	if ic.Stop {
		//Output the rest of the buffered assembly to the file.
		ic.Assembly("\n")
		return false
	}
	
	//It might be nice to have a Compiler.Register(token, ScanFunc)
	if f, ok := Tokens[token]; ok {
		f(ic)
		return true
	}
	
	//These are all the tokens in ilang.
	switch token {
		case "\n", ";":
		
		case "plugin":
			ic.ScanPlugin()
		
		case "@":
			ic.Language = ic.Scan(Name)
			if ic.Language == "ch" {
				ic.Language = "zh-CN"
			}
			ic.Translation = true
		
		case "type":
			ic.Header = false
			ic.ScanType()
		
		case "gui":
			ic.Header = false
			ic.ScanGui()
		
		case "fork":
			name := ic.Scan(Name)
			ic.Scan('(')
			ic.Fork = true
			ic.ScanFunctionCall(name)
			ic.Scan(')')

		case "end":
			ic.LoseScope()
			
		case "}":					
			ic.LoseScope()
			
			
		default:

			var skip bool
			for _, f := range Defaults {
				if f(ic) {
					skip = true
					break
				}
			}
			if skip {
				return true
			}
			
			ic.NextToken = token
			ic.ScanStatement()
	}
	return true
}
