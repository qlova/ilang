package main

import "text/scanner"
import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"io"
	"flag"
	"path"
)

const (
	Iname = scanner.Ident
	
)

//This holds the definition of a function.
type Function struct {
	Exists bool
	Args []TYPE
	Returns []TYPE
	
	//Is this really not a function?
	Ghost bool
	
	//Is this a local?
	Local bool
	
	Inline bool
	Data string
	Loaded bool
	Load string
	
	Variadic bool
}


//Deal with scoping for variables.
type Variables map[string]TYPE

var Scope []Variables

func GainScope() {
	Scope = append(Scope, make(map[string]TYPE))
}

func GetVariable(name string) TYPE {
	for i:=len(Scope)-1; i>=0; i-- {
		if v, ok := Scope[i][name]; ok {
			return v
		}
	}
	return UNDEFINED
}

func NextToken(s *scanner.Scanner, verify rune) string {
	tok := s.Scan()
	if tok != verify {
		RaiseError(s, "Unexpected "+s.TokenText())
	}
	return s.TokenText()
}

func SetVariable(name string, sort TYPE) {
	Scope[len(Scope)-1][name] = sort
}

func Expecting(s *scanner.Scanner, token string) {
	if s.TokenText() != token {
		RaiseError(s, "Expecting "+token+" found "+s.TokenText())
	}
}

func Collect(output io.Writer, name string, typ TYPE) {
	if typ >= USER {
		t := GetType(typ)
		for i, element := range t.Elements {
			if element == STRING || element >= USER {
				fmt.Fprintf(output, "#<GC>\n")
	
				unique++
				var id = unique
				fmt.Fprintf(output, "PLACE %s\n", name)
				fmt.Fprintf(output, "PUSH %v\n", i)
				fmt.Fprintf(output, "GET %s%v\n", "i+gc+", unique)
				fmt.Fprintf(output, "IF %s%v\n", "i+gc+", id)
				if element >= USER {
					fmt.Fprintf(output, "PUSH %s%v\n", "i+gc+", unique)
					fmt.Fprintf(output, "HEAP\n")
					unique++
					fmt.Fprintf(output, "GRAB %s%v\n", "i+gc+", unique)
					Collect(output, fmt.Sprint("i+gc+", unique), element)
				}
				fmt.Fprintf(output, "MUL %s%v %s%v -1\n", "i+gc+", id, "i+gc+", id)
				fmt.Fprintf(output, "PUSH %s%v\n", "i+gc+", id)
				fmt.Fprintf(output, "HEAP\nEND\n")
		
				fmt.Fprintf(output, "#</GC>\n")
			}
		}
	}
}

func LoseScope(output io.Writer) {

	//Erm garbage collection???
	for name, variable := range Scope[len(Scope)-1] {
		if Scope[len(Scope)-1][name+"."] != 1 { //Protected variables
			Collect(output, name, variable)
		}
	}

	Scope = Scope[:len(Scope)-1]
}

func RaiseError(s *scanner.Scanner, message string) {
	fmt.Fprintf(os.Stderr, "[%v] %v\n", s.Pos(), message)
	os.Exit(1)
}

var functions = make( map[string]Function)
var methods = make( map[string]bool)
var unique int

//This means that the result should not be garbage collected at the end of the scope.
var Protected bool

var CurrentFunction Function

var ExpressionType TYPE
var FinalExpressionType TYPE

//TODO put this in a different file. (expression.go?)
func expression(s *scanner.Scanner, output io.Writer, param ...bool) string {
	
	//Do we need to shunt? This is for operator precidence. (defaults to true)
	var shunting bool = len(param) <= 0 || param[0]
	var token = s.TokenText()
	
	defer func() {
		if FinalExpressionType > 0 {
			ExpressionType = FinalExpressionType
			FinalExpressionType = 0
		}
	}()
	
	//Types.
	if string(s.Peek()) != "(" {
		//ExpressionType = ITYPE
		switch token {
			case "number", "text", "letter", "array", "Something", "pipe", "type", "function":
				return fmt.Sprint(int(StringToType[token])) 
		}
	}
	
	//fmt.Println("TOKEN ", token)
	
	if len(token) <= 0 {
		fmt.Println("Empty expression!")
		ExpressionType = UNDEFINED
		return ""
	}
	
	if token == "true" {
		ExpressionType = NUMBER
		if shunting {
			return shunt("1", s, output)
		} else {
			return "1"
		}
	}
	
	if token == "error" {
		ExpressionType = NUMBER
		if shunting {
			return shunt("ERROR", s, output)
		} else {
			return "ERROR"
		}
	}
	
	if token == "false" {
		ExpressionType = NUMBER
		if shunting {
			return shunt("0", s, output)
		} else {
			return "0"
		}
	}

	//If there is a quotation mark, parse the string. 
	if token[0] == '"' {
		ExpressionType = STRING
		return ParseString(s, output, shunting)
	}
	
	if  token[0] == '[' {
		defer func() {
			ExpressionType = ARRAY
		}()
		return ParseArray(s, output, shunting)
	}
	
	//Deal with runes. 
	//	eg. 'a' -> 97
	if len(token) == 3 && token[0] == '\'' && token[2] == '\'' {
		ExpressionType = LETTER
		
		var value = strconv.Itoa(int(s.TokenText()[1]))
		
		if shunting {
			return shunt(value, s, output)
		} else {
			return value
		}
	} else if s.TokenText() == `'\n'` {
		ExpressionType = LETTER
		
		var value = strconv.Itoa(int('\n'))
		
		if shunting {
			return shunt(value, s, output)
		} else {
			return value
		}
	} else if s.TokenText() == `'\r'` {
		ExpressionType = LETTER
		
		var value = strconv.Itoa(int('\r'))

		if shunting {
			return shunt(value, s, output)
		} else {
			return value
		}
	} else if len(token) > 2 && token[0] == '0' && token[1] == 'x' { 
		ExpressionType = NUMBER
		if shunting {
			return shunt(token, s, output)
		} else {
			return token
		}
	}

	
	//Parse method call.
	if methods[token] && s.Peek() == '@' {
		return ParseFunction(s.TokenText(), s, output, shunting)
	}
	
	//Parse function call.
	if functions[token].Exists {
		if s.Peek() != '(' {
			ExpressionType = FUNCTION
			unique++
			id := "i+func+"+fmt.Sprint(unique)
			fmt.Fprintf(output, "SCOPE %v\nTAKE %v\n", token, id)
			return id
		}
		return ParseFunction(s.TokenText(), s, output, shunting)
	}
	
	if token[0] == '-' {
		//Is it a literal number? Then just return it.
		//OR is it a variable?
		if _, err := strconv.Atoi("-"+s.TokenText()); err == nil{
			s.Scan()
			ExpressionType = NUMBER
			if shunting {
				return shunt("-"+s.TokenText(), s, output)
			} else {
				return "-"+s.TokenText()
			}
		} else {
			return shunt("0", s, output, false)
			//fmt.Println(s.Pos(), "Unexpected - sign.")
			//os.Exit(1)
		}
	}

	//Is it a literal number? Then just return it.
	//OR is it a variable?
	if _, err := strconv.Atoi(token); err == nil{
		ExpressionType = NUMBER
		if shunting {
			return shunt(token, s, output)
		} else {
			return token
		}
	}
	
	if sort := GetVariable(token); sort != UNDEFINED {
		//By default, strings are indexed as it's first element when it's a single character.
		if (sort == STRING || sort >= USER) && len(token) == 1 && OperatorFunction {
			if sort >= USER {
				t := GetType(sort)
				ExpressionType = t.Elements[0]
			} else {
				ExpressionType = NUMBER
			}
			
			if ExpressionType == NUMBER {
				unique++
				id := "i+tmp+"+s.TokenText()+fmt.Sprint(unique)
				fmt.Fprintf(output, "PLACE %v\nPUSH 0\nGET %v\n", string(rune(s.TokenText()[0])), id)
				
				if shunting {
					return shunt(id, s, output)
				} else {
					return id
				}
			} else {
				fmt.Println(s.Pos(), "Cannot access type field by default (ONLY NUMBER FIELDS ARE SUPPORTED ATM) Not ", ExpressionType)
				os.Exit(1)
			}
			
			
		} else {
			ExpressionType = sort
		}
		if shunting {
			return shunt(token, s, output)
		} else {
			return token
		}
	}
	
	if GetVariable("__new__") > 0 || methods[CurrentFunctionName] {
		if _, ok := GetType(LastDefinedType).Table[token]; ok {
			ExpressionType = LastDefinedType
				
			token = IndexUserType(s, output, LastDefinedType.String(), token)
			if shunting {
				return shunt(token, s, output)
			} else {
				return token
			}
		}
	}
	
	if sort := GetVariable("gui_"+token); sort != UNDEFINED {
		ExpressionType = sort
		if shunting {
			return shunt("gui_"+token, s, output)
		} else {
			return "gui_"+token
		}
	}
		
	//Do some special maths which people will hate about I.
	// a=2; b=4; ab
	// ab is 8
	//
	//When a is a string, it will index itself "0".
	if OperatorFunction && GetVariable(string(rune(token[0]))) != UNDEFINED {
	
		var lastID, id, id2 string
		for i:=0; i<len(token); i++ {
		
			isvariable := GetVariable(string(rune(token[i])))
			if isvariable >= USER {
				isvariable = STRING
			}
			if isvariable == UNDEFINED {
				//println(token)
				goto notab
			}
			
			if isvariable == STRING && len(token) > i+1 && token[i+1] == 'i' {
				ExpressionType = NUMBER
			
				unique++
				id = "i+tmp+"+s.TokenText()+fmt.Sprint(unique)
				fmt.Fprintf(output, "PLACE %v\nPUSH 1\nGET %v\n", string(rune(token[i])), id)
				i++
			} else if isvariable == STRING {
				ExpressionType = NUMBER
			
				unique++
				id = "i+tmp+"+s.TokenText()+fmt.Sprint(unique)
				fmt.Fprintf(output, "PLACE %v\nPUSH 0\nGET %v\n", string(rune(token[i])), id)
			
			} else if isvariable == NUMBER {
				ExpressionType = NUMBER
				id = string(rune(token[i+1]))
			}
			
			if lastID != "" {
			
				unique++
				id2 = "i+tmp+"+s.TokenText()+fmt.Sprint(unique)
				fmt.Fprintf(output, "VAR %v\n", id2)
				fmt.Fprintf(output, "MUL %v %v %v\n", id2, lastID, id)
				id = id2
			}
		
			lastID = id
			
		}
		
		if shunting {
			return shunt(lastID, s, output)
		} else {
			return lastID
		}
		
		
	}
	notab:
	
	if token[0] == '(' {
		s.Scan()
		if shunting {
			return shunt(expression(s, output), s, output)
		} else {
			return expression(s, output)
		}
	}
	
	ExpressionType = UNDEFINED
	RaiseError(s, "'"+s.TokenText()+"' is undefined!")
	
	if shunting {
		return shunt(s.TokenText(), s, output)
	} else {
		return s.TokenText()
	}
}

var FirstIssue, Issues, Switching bool

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		return
	}
	
	//Open the output file with the file type replaced to .u
	output, err := os.Create(flag.Arg(0)[:len(flag.Arg(0))-2]+".u")
	if err != nil {
		return
	}
	
	ifile, err := os.Create(path.Dir(flag.Arg(0))+"/ilang.u")
	if err != nil {
		return
	}
	
	//Add builtin functions to file.
	builtin(ifile)
			
	
	//Startup the scanner.
	var s scanner.Scanner
	//Keeping track of multiple files when importing.
	var scanners []scanner.Scanner
	
	s.Init(file)
	s.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
	
	//TODO cleanup file from here forward.
	var software, softwareBlock bool
	
	defer func() {
		if !software && GUIEnabled && GUIMain {
			output.Write([]byte("SOFTWARE\n"))
			output.Write([]byte("SHARE gui_main\nRUN gui\n"))
			output.Write([]byte("LOOP\n"))
			output.Write([]byte("SHARE i_newline\n"))
			output.Write([]byte("RUN grab\n"))
			output.Write([]byte("IF ERROR\n"))
			output.Write([]byte("EXIT\n"))
			output.Write([]byte("END\n"))
			output.Write([]byte("REPEAT\n"))
			output.Write([]byte("EXIT\n"))
			LoadFunction("grab")
			LoadFunction("gui")
			LoadFunction("output_m_file")
			LoadFunction("reada_m_file")
		}
	}()
	
	var SwitchCase string
	var FirstSwitch bool
	
	GainScope()
	SetVariable("error", NUMBER)
	
	fmt.Fprintf(ifile, `DATA i_newline "\n"`+"\n")
	fmt.Fprintf(output, `.import ilang`+"\n")
	
	var tok rune
	for {
		tok = s.Scan()
		
		if tok == scanner.EOF {
			if len(scanners) > 0 {
				s = scanners[len(scanners)-1]
				scanners = scanners[:len(Scope)-1]
				continue
			} else {
				return
			}
		}
		
		switch s.TokenText() {
			case "import":
				s.Scan()
			file, err := os.Open(s.TokenText()+".i")
				if err != nil {
					RaiseError(&s, "Cannot import "+s.TokenText()+", does not exist!")
				}
				scanners = append(scanners, s)
				
				s = scanner.Scanner{}
				s.Init(file)
				s.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
				continue
				
			//LOOPS
			
			case "repeat", "break":
				if s.TokenText() == "repeat" {
					LoseScope(output)
				}
				fmt.Fprintf(output, "%v", strings.ToUpper(s.TokenText())+"\n")
			
			case "fork":
				s.Scan()
				function := s.TokenText()
				if !functions[function].Exists {
					fmt.Println(s.Pos(), "Expecting a function but instead, found ", s.TokenText())
					return
				}
				ParseFunction(s.TokenText(), &s, output, false, false)
				fmt.Fprintf(output, "FORK %v\n", function)
				for _, v := range functions[function].Returns {
					fmt.Fprintf(output, "PULL %v\n", v.Push()[4:])
				}
			
			case "do":
				GainScope()
				fmt.Fprintf(output, "LOOP\n")
			
			case "done":
				fmt.Fprintf(output, "REPEAT\n")
			
			case "while":
				s.Scan()
				fmt.Fprintf(output, "IF %v\nERROR 0\nELSE\nBREAK\nEND\n", expression(&s, output))
				fmt.Fprintf(output, "REPEAT\n")
				GainScope()
			
			case "exit":
				fmt.Fprintf(output, "EXIT\n")
			
			case "\n", ";":
			
			case "!":
				output.Write([]byte("ERROR 0\n"))
			
			case "}", "end":
				_, ok := Scope[len(Scope)-1]["loop"]
				if ok {
					LoseScope(output)
					output.Write([]byte("REPEAT\n"))
					continue
				}
			
				nesting, ok := Scope[len(Scope)-1]["elseif"]
				if ok {
					for i:=0; TYPE(i) < nesting; i++ {
						output.Write([]byte("END\n"))
					}
				}
				if len(Scope) > 2 {
					LoseScope(output)
					output.Write([]byte("END\n"))
				} else if softwareBlock {
					softwareBlock = false
					LoseScope(output)
					output.Write([]byte("EXIT\n"))
				} else {
					if OperatorFunction {
						OperatorFunction = false
						output.Write([]byte("SHARE c\n"))
					}
					if GetVariable("__new__") > 0 {
						fmt.Fprintf(output, "SHARE %s\n", LastDefinedType.String())
					}
					LoseScope(output)
					output.Write([]byte("RETURN\n"))
				}
				if Issues && !FirstIssue {
					LoseScope(output)
					output.Write([]byte("END\n"))
					Issues = false
				} else if Issues {
					Issues = false
				}
				
				if Switching && !FirstSwitch {
					LoseScope(output)
					output.Write([]byte("END\n"))
					Switching = false
				} else if Switching {
					Switching = false
				}
				
			
			case "else":
				
				nesting, ok := Scope[len(Scope)-1]["elseif"]
				if !ok {
					nesting = 0
				}
				LoseScope(output)
				fmt.Fprintf(output, "ELSE\n")
				GainScope()
				SetVariable("elseif", nesting)
			
			case "if", "elseif":
				
				if s.TokenText() == "if" {
					GainScope()
					s.Scan()
					fmt.Fprintf(output, "IF %v\n", expression(&s, output))
				} else {
					nesting, ok := Scope[len(Scope)-1]["elseif"]
					if !ok {
						nesting = 0
					}
					LoseScope(output)
					GainScope()
					SetVariable("elseif", nesting+1)
					s.Scan()
					fmt.Fprintf(output, "ELSE\n")
					condition := expression(&s, output)
					fmt.Fprintf(output, "IF %v\n", condition)
				}
				
				
			//Inline universal assembly.
			case ".":
				s.Scan()
				output.Write([]byte(strings.ToUpper(s.TokenText())))
				
				var data bool
				if s.TokenText() == "data" {
					data = true
				}
				
				for tok = s.Scan(); tok != scanner.EOF; {
					if data {
						SetVariable(s.TokenText(), STRING)
					}
				
					if s.TokenText() == "\n" {
						output.Write([]byte("\n"))
						break
					}
					output.Write([]byte(" "+s.TokenText()))
					s.Scan()
				}
			
			case "return":
				s.Scan()
				if CurrentFunction.Exists {
					GainScope()
					if s.TokenText() != "\n" {
						if len(CurrentFunction.Returns) > 0 {
							var ReturnType = CurrentFunction.Returns[0]
							
							if ReturnType != USER {
								fmt.Fprintf(output, "%s %v\n", ReturnType.Push(), expression(&s, output))
							} else {
								var name = expression(&s, output)
								output.Write([]byte("SHARE "+name+"\n"))
								CurrentFunction.Returns[0] = ExpressionType
								Scope[len(Scope)-2][name] = 0
							} 
						}
					}
					Scope = Scope[:len(Scope)-1]
				}
				LoseScope(output)
				GainScope()
				if len(Scope) > 2 {
					output.Write([]byte("RETURN\n"))
				}
			
			case "software":
				software = true
				
				CurrentFunctionName = ""
				
				output.Write([]byte("SOFTWARE\n"))
				if GUIEnabled {
					output.Write([]byte("SHARE gui_main\nRUN gui\n"))
					LoadFunction("gui")
					LoadFunction("output_m_file")
				}
				s.Scan()
				if s.TokenText() != "{" {
					fmt.Println(s.Pos(), "Expecting { found ", s.TokenText())
					return
				}
				GainScope()
				softwareBlock = true
				
				
			
			case "issues":
				output.Write([]byte("IF ERROR\n"))
				output.Write([]byte("VAR ISSUE\n"))
				output.Write([]byte("ADD ISSUE 0 ERROR\n"))
				output.Write([]byte("ADD ERROR 0 0\n"))
				s.Scan()
				if s.TokenText() != "{" {
					fmt.Println(s.Pos(), "Expecting { found ", s.TokenText())
					return
				}
				s.Scan()
				if s.TokenText() != "\n" {
					fmt.Println(s.Pos(), "Expecting newline found ", s.TokenText())
					return
				}
				GainScope()
				FirstIssue = true
				Issues = true
			
			case "switch":	
				s.Scan()
				SwitchCase = expression(&s, output)
				GainScope()
				Expecting(&s, "{")
				s.Scan()
				FirstSwitch = true
				Switching = true
			
			case "default":
				
				if !Switching {
					RaiseError(&s, "Unexpected 'default'")
				}
				LoseScope(output)
				GainScope()
				s.Scan()
				fmt.Fprintf(output, "ELSE\n")
			
			case "case":
				
				if FirstSwitch {
					GainScope()
					s.Scan()
					unique++
					fmt.Fprintf(output, "VAR %v\n", "i+issue+"+fmt.Sprint(unique))
					fmt.Fprintf(output, "SEQ %v %v %v\n", "i+issue+"+fmt.Sprint(unique), SwitchCase, expression(&s, output))
					fmt.Fprintf(output, "IF %v\n", "i+issue+"+fmt.Sprint(unique))
					FirstSwitch = false
				} else {
					nesting, ok := Scope[len(Scope)-1]["elseif"]
					if !ok {
						nesting = 0
					}
					LoseScope(output)
					GainScope()
					SetVariable("elseif", nesting+1)
					s.Scan()
					fmt.Fprintf(output, "ELSE\n")
					fmt.Fprintf(output, "VAR %v\n", "i+issue+"+fmt.Sprint(unique))
					fmt.Fprintf(output, "SEQ %v %v %v\n", "i+issue+"+fmt.Sprint(unique), SwitchCase, expression(&s, output))
					fmt.Fprintf(output, "IF %v\n", "i+issue+"+fmt.Sprint(unique))
				}
			
			case "issue":
				
				if FirstIssue {
					GainScope()
					s.Scan()
					unique++
					fmt.Fprintf(output, "VAR %v\n", "i+issue+"+fmt.Sprint(unique))
					fmt.Fprintf(output, "SEQ %v ISSUE %v\n", "i+issue+"+fmt.Sprint(unique), expression(&s, output))
					fmt.Fprintf(output, "IF %v\n", "i+issue+"+fmt.Sprint(unique))
					FirstIssue = false
				} else {
					nesting, ok := Scope[len(Scope)-1]["elseif"]
					if !ok {
						nesting = 0
					}
					LoseScope(output)
					GainScope()
					SetVariable("elseif", nesting+1)
					s.Scan()
					fmt.Fprintf(output, "ELSE\n")
					fmt.Fprintf(output, "VAR %v\n", "i+issue+"+fmt.Sprint(unique))
					fmt.Fprintf(output, "SEQ %v ISSUE %v\n", "i+issue+"+fmt.Sprint(unique), expression(&s, output))
					fmt.Fprintf(output, "IF %v\n", "i+issue+"+fmt.Sprint(unique))
				}
			
			case "print":
				s.Scan()
				s.Scan()
				var arg = expression(&s, output)
				fmt.Fprintf(output, "%s %s\n", ExpressionType.Push(), arg)
				ParseFunction("text_m_"+ExpressionType.String(), &s, output, false, true, true, true)
				fmt.Fprintf(output, "STDOUT\n")
				
				for s.TokenText() == "," {
					s.Scan()
					var arg = expression(&s, output)
					fmt.Fprintf(output, "%s %s\n", ExpressionType.Push(), arg)
					ParseFunction("text_m_"+ExpressionType.String(), &s, output, false, true, true, true)
					fmt.Fprintf(output, "STDOUT\n")
				}
				
				
				fmt.Fprintf(output, "SHARE i_newline\n")
				fmt.Fprintf(output, "STDOUT\n")
			
			//New type decleration.
			case "type":
				ParseTypeDef(&s, output)
				
			//Compiles function declerations.
			case "function", "method":
				GainScope()
				ParseFunctionDef(&s, output)
			
			case "new":
				GainScope()
				SetVariable("__new__", 1)
				ParseConstructDef(&s, output)
				
			//ParseConstant
			case "const":
				s.Scan()
				var name = s.TokenText()
				s.Scan()
				Expecting(&s, "=")
				s.Scan()
				var value = expression(&s, output, false)
				if ExpressionType != NUMBER && ExpressionType != LETTER {
					RaiseError(&s, "Constant must be a numerical value!")
				} 
				fmt.Fprintf(output, ".const %s %s\n", name, value)
				SetVariable(name, ExpressionType)
				
			case "var", "for":
				//RESET protection
				Protected =false
			
				var forloop = s.TokenText() == "for"
				s.Scan()
				
				var name = s.TokenText()
				s.Scan()
				
				var name2 string
				if forloop && s.TokenText() == "," {
					s.Scan()
					name2 = s.TokenText()
					s.Scan()
				}
				
				if forloop && s.TokenText() == "in" {
					s.Scan()
					array := expression(&s, output)
					
					unique ++
					test := "i+in+"+fmt.Sprint(unique)
					unique ++
					var i, v string
					if name2 != "" {
						i = name
						v = name2
					} else {
						i = "i+in+"+fmt.Sprint(unique)
						v = name
					}
					unique ++
					backup := "i+in+"+fmt.Sprint(unique)
					
					fmt.Fprintf(output, `
VAR %s
VAR %s
LOOP
	ADD %s 0 %s
	PLACE %s
		PUSH %s
		GET %s
	SGE %s %s #%s
	IF %s
		BREAK
	END
	ADD %s %s 1
	
`, i,backup, i, backup, array, i, v, test, i, array, test, backup, i)

					GainScope()
					SetVariable(i, NUMBER)
					SetVariable(v, NUMBER)
					SetVariable("loop", 1)
					continue
				}
				
				//Over in a for loop.
				if forloop && s.TokenText() == "over" {
					s.Scan()
					Expecting(&s, "[")
					s.Scan()
					
					low := expression(&s, output)
					s.Scan()
					high := expression(&s, output)
					
					Expecting(&s, "]")
					s.Scan()
					
					unique ++
					test := "i+over+"+fmt.Sprint(unique)
					unique ++
					backup := "i+back+"+fmt.Sprint(unique)
					fmt.Fprintf(output, `
VAR %s
VAR %s
ADD %s 0 %s
ADD %s 0 %s
LOOP
	SNE %s %s %s
	ADD %s 0 %s
	IF %s
		SLT %s %s %s
		IF %s
			ADD %s %s 1
		ELSE
			SUB %s %s 1
		END
	ELSE
		BREAK
	END
`, name, backup, backup, low, name, low, test, name, high, name, backup, test, test, name, high, test, backup, name, backup, name)
					GainScope()
					SetVariable(name, NUMBER)
					SetVariable("loop", 1)
					continue
				}
				
				if s.TokenText() == "is" {
					s.Scan()
					fmt.Fprintf(output, "ARRAY %v\n", name)
					SetVariable(name, ParseConstructor(&s, output))
					continue
				}
				s.Scan()
				var set = expression(&s, output)
				if ExpressionType == NUMBER || ExpressionType == LETTER {
					fmt.Fprintf(output, "VAR %v\nADD %v 0 %v\n", name, name, set)
				}
				if ExpressionType == STRING || ExpressionType == ARRAY {
					fmt.Fprintf(output, "SHARE %v\nGRAB %v\n", set, name)
				}
				if ExpressionType == FUNCTION || ExpressionType == FILE {
					fmt.Fprintf(output, "RELAY %v\nTAKE %v\n", set, name)
				}
				if ExpressionType >= USER {
					fmt.Fprintf(output, "SHARE %v\nGRAB %v\n", set, name)
				}
				if Protected {
					SetVariable(name+".", 1)
				}
				SetVariable(name, ExpressionType)
				if !forloop {
					continue
				}
				fallthrough
			case "loop":
				fmt.Fprintf(output, "LOOP\n")
				s.Scan()
				if s.TokenText() != "\n" {
					fmt.Fprintf(output, "IF %v\nERROR 0\nELSE\nBREAK\nEND\n", expression(&s, output))
				}
				GainScope()
				SetVariable("loop", 1)
			
			case "gui":
				if !softwareBlock {
					ParseGUIDef(&s, output)
					continue
				}
				fallthrough
			default:
			
				var name = s.TokenText()
				
				if _, ok := StringToType[name]; ok {
					if len(Scope) == 1 {
						ParseOperator(&s, output)
						continue
					}
				}
				
				//Method lines.
				if methods[name] && s.Peek() == '@' {
					expression(&s, output)
					continue
				}
				
				//Function lines.
				if functions[name].Exists {
					var returns = functions[name].Returns
					var f = functions[name]
					f.Returns = nil
					functions[name] = f
						expression(&s, output)
					f.Returns = returns
					functions[name] = f
					continue
				}
				
				s.Scan()
				switch s.TokenText() {
					case "(":
						s.Scan()
						output.Write([]byte("EXE "+name+" \n"))
					case "&":
						s.Scan()
						variable := expression(&s, output)
						if ExpressionType == NUMBER || ExpressionType == LETTER {
							fmt.Fprintf(output, "PLACE %v\n", name)
							fmt.Fprintf(output, "PUT %v\n", variable)
						} else if ExpressionType == STRING {
							output.Write([]byte("JOIN "+name+" "+name+" "+variable+" \n"))
						}
						
					//TODO Allow assigning to non-numeric types?
					case "=", "is":
						
						//Set array in operator function.
						if OperatorFunction && GetVariable(string(rune(name[0]))) >= USER  {
							s.Scan()
							if len(name) == 1 {
								fmt.Fprintf(output, "PLACE %v\nPUSH 0\nSET %v\n", name, expression(&s, output))
							}
							if len(name) == 2 && name[1] == 'i' {
								fmt.Fprintf(output, "PLACE %v\nPUSH 1\nSET %v\n", string(rune(name[0])), expression(&s, output))
							}
							
						} else {
							if GetVariable(name) != NUMBER && GetVariable(name) != LETTER {
								if GetVariable(name) == UNDEFINED {
								
									if GetVariable("__new__") > 0 || methods[CurrentFunctionName]  {
										AssignToUserType(&s, output, LastDefinedType.String(), name)
									} else {
										RaiseError(&s, name+" is undefined!")
									}
									
								//Asigning to "Something", requires that a data field is filled with a pointer or value.
								} else if GetVariable(name) == SOMETHING {
									s.Scan()
									something := expression(&s, output)
									switch ExpressionType {
										case NUMBER, LETTER:
											fmt.Fprintf(output, "PLACE %v\nPUSH 0\nSET %v\nPUSH 1\nSET %v\n", name, something, int(ExpressionType))
										case STRING, ARRAY:
											unique++
											address := fmt.Sprint("i+elem+",unique)
											fmt.Fprintf(output, "SHARE %v\n PUSH 0\nHEAP\nPULL %v\n", something, address)
											fmt.Fprintf(output, "PLACE %v\nPUSH 0\nSET %v\nPUSH 1\nSET %v\n", name, address, int(ExpressionType))
											//methods := GenMethodList(output, NUMBER)
										default:
											if ExpressionType >= USER {
												unique++
												address := fmt.Sprint("i+elem+",unique)
												fmt.Fprintf(output, "SHARE %v\n PUSH 0\nHEAP\nPULL %v\n", something, address)
												fmt.Fprintf(output, "PLACE %v\nPUSH 0\nSET %v\nPUSH 1\nSET %v\n", name, address, int(ExpressionType))
											} else {
												RaiseError(&s, "Cannot assign "+ExpressionType.String()+" to Something (NOT SUPPORTED)")
											}
									}
									
								} else if GetVariable(name) == FILE {
									s.Scan()
									fmt.Fprintf(output, "RELAY %v\nRELOAD %v\n", expression(&s, output), name)
								} else if GetVariable(name) == STRING || GetVariable(name) >= USER || GetVariable(name) == ARRAY {
									s.Scan()
									fmt.Fprintf(output, "PLACE %v\nRENAME %v\n", expression(&s, output), name)
								} else {
									RaiseError(&s, "Cannot assign to "+name+"! Not an assignable type!.")
								}
							} else if GetVariable(name) == NUMBER ||  GetVariable(name) == LETTER {
								s.Scan()
								if name == "error" {
									name = "ERROR"
								}
								fmt.Fprintf(output, "ADD %v 0 %v\n", name, expression(&s, output))
							} else {
								//if GetVariable(name) != SOMETHING && ExpressionType != GetVariable(name) {
									RaiseError(&s, "Cannot assign "+ExpressionType.String()+" to "+GetVariable(name).String())
								//}
							}
						}
					
					default:
						if s.TokenText() == "[" {
							t := GetVariable(name)
							if t == STRING || t == ARRAY {
								s.Scan()
								var index = expression(&s, output, true)
								
								s.Scan()
								if s.TokenText() != "=" {
									RaiseError(&s, "Expecting =, found "+s.TokenText())
								} 
								s.Scan()
								fmt.Fprintf(output, "PLACE %v\nPUSH %v\nSET %v\n", name, index, expression(&s, output))
							} else {
								RaiseError(&s, "Cannot index "+name+", not an array!!! ("+t.String()+")")
							}
							continue
						}
					
						if (len(s.TokenText()) > 0 && s.TokenText()[0] == '.') || s.TokenText() == "." {
						
							
							//Index structures.
							// Like type Complex { real, imag }
							// Complex.real = 2
							
							if t := GetVariable(name); t >= USER {
								s.Scan()
								stringdex := s.TokenText()
								s.Scan()
							
								AssignToUserType(&s, output, name, stringdex)
							} 
							
						} else {
						
							//i++ 
							if GetVariable(name) == STRING || GetVariable(name) == ARRAY {
								if string(s.Peek()) == "=" {
									var op = s.TokenText()
									s.Scan()
									s.Scan()
									switch op {
										case "+":
											fmt.Fprintf(output, "JOIN %v %v %v\n", name, name, expression(&s, output))	
										default:
											RaiseError(&s, "Unexpected "+s.TokenText())
									}
								} else {
									RaiseError(&s, "Unexpected "+s.TokenText())
								}
								continue
							}
							
							if GetVariable(name) == NUMBER {
								if s.TokenText() == "+" && string(s.Peek()) != "=" {
									s.Scan()
									if s.TokenText() == "+" {
										fmt.Fprintf(output, "ADD %v %v 1\n", name, name)
									} else {
										RaiseError(&s, "Unexpected "+s.TokenText()+", expecting '+'")
									}
								} else if s.TokenText() == "-" && string(s.Peek()) != "=" {
									s.Scan()
									if s.TokenText() == "-" {
										fmt.Fprintf(output, "ADD %v %v 1\n", name, name)
									} else {
										RaiseError(&s, "Unexpected "+s.TokenText()+", expecting '-'")
									}
								} else if string(s.Peek()) == "=" {
									var op = s.TokenText()
									s.Scan()
									s.Scan()
									switch op {
										case "+":
											fmt.Fprintf(output, "ADD %v %v %v\n", name, name, expression(&s, output))
										case "-":
											fmt.Fprintf(output, "SUB %v %v %v\n", name, name, expression(&s, output))
										case "*":
											fmt.Fprintf(output, "MUL %v %v %v\n", name, name, expression(&s, output))
										case "^":
											fmt.Fprintf(output, "POW %v %v %v\n", name, name, expression(&s, output))	
										case "/":
											fmt.Fprintf(output, "DIV %v %v %v\n", name, name, expression(&s, output))	
										default:
											RaiseError(&s, "Unexpected "+s.TokenText())
									}
								} else {
									RaiseError(&s, "Unexpected "+s.TokenText())
								}
							} else {
					
								if name == "" {
									return	
								}
								RaiseError(&s, "Unexpected "+s.TokenText())
							}
						}
				}
				
		}
	}
}
