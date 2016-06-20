package main

import "text/scanner"
import (
	"os"
	"fmt"
	"strings"
	"strconv"
	"io"
	"flag"
)

//These are the 4 types in I.
const (
	FUNCTION = iota
	STRING
	NUMBER
	FILE
	
	UNDEFINED
)

//This holds the definition of a function.
type Function struct {
	Exists bool
	Args []int
	Returns []int
	
	//Is this a local?
	Local bool
	
	Variadic bool
}


//Deal with scoping for variables.
type Variables map[string]int

var Scope []Variables

func GainScope() {
	Scope = append(Scope, make(map[string]int))
}

func GetVariable(name string) int {
	for i:=len(Scope)-1; i>=0; i-- {
		if v, ok := Scope[i][name]; ok {
			return v
		}
	}
	return UNDEFINED
}

func SetVariable(name string, sort int) {
	Scope[len(Scope)-1][name] = sort
}

func LoseScope() {
	Scope = Scope[:len(Scope)-1]
}

var functions = make( map[string]Function)
var methods = make( map[string]bool)
var unique int

var CurrentFunction Function

var ExpressionType int
func expression(s *scanner.Scanner, output io.Writer, param ...bool) string {
	
	//Do we need to shunt? This is for operator precidence. (defaults to true)
	var shunting bool = len(param) <= 0 || param[0]
	var token = s.TokenText()
	
	//fmt.Println("TOKEN ", token)
	
	if len(token) <= 0 {
		fmt.Println("Empty expression!")
		return ""
	}
	
	if token == "true" {
		if shunting {
			return shunt("1", s, output)
		} else {
			return "1"
		}
	}
	
	if token == "false" {
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
		ExpressionType = STRING
		return ParseArray(s, output, shunting)
	}
	
	//Deal with runes. 
	//	eg. 'a' -> 97
	if len(token) == 3 && token[0] == '\'' && token[2] == '\'' {
		ExpressionType = NUMBER
		defer s.Scan()
		return strconv.Itoa(int(s.TokenText()[1]))
	} else if s.TokenText() == `'\n'` {
		ExpressionType = NUMBER
		defer s.Scan()
		return strconv.Itoa(int('\n'))
	}

	
	//Parse method call.
	if methods[token] && s.Peek() == '@' {
		return ParseFunction(s, output, shunting)
	}
	
	//Parse function call.
	if functions[token].Exists {
		if s.Peek() != '(' {
			ExpressionType = FUNCTION
			unique++
			id := "i+func+"+fmt.Sprint(unique)
			fmt.Fprintf(output, "FUNC %v %v\n", id, token)
			return id
		}
		return ParseFunction(s, output, shunting)
	}
	
	if token[0] == '-' {
		s.Scan()
		//Is it a literal number? Then just return it.
		//OR is it a variable?
		if _, err := strconv.Atoi("-"+s.TokenText()); err == nil{
			ExpressionType = NUMBER
			if shunting {
				return shunt("-"+s.TokenText(), s, output)
			} else {
				return "-"+s.TokenText()
			}
		} else {
			fmt.Println("Unexpected - sign.")
			os.Exit(1)
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
		ExpressionType = sort
		if shunting {
			return shunt(token, s, output)
		} else {
			return token
		}
	}
		
	//Do some special maths which people will hate about I.
	// a=2; b=4; ab
	// ab is 8
	if GetVariable(string(rune(token[0]))) != UNDEFINED {
		ExpressionType = NUMBER
		if len(token) == 2 {
			if GetVariable(string(rune(token[1]))) != UNDEFINED {
				unique++
				id := "i+tmp+"+s.TokenText()+fmt.Sprint(unique)
				fmt.Fprintf(output, "VAR %v\n", id)
				fmt.Fprintf(output, "MUL %v %v %v\n", id, string(rune(s.TokenText()[0])), string(rune(s.TokenText()[1])))
				
				if shunting {
					return shunt(id, s, output)
				} else {
					return id
				}
			}
		}
	}
	
	if token[0] == '(' {
		s.Scan()
		if shunting {
			return shunt(expression(s, output), s, output)
		} else {
			return expression(s, output)
		}
	}
	
	if shunting {
		return shunt(s.TokenText(), s, output)
	} else {
		return s.TokenText()
	}
}

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
	
	//Add builtin functions to file.
	builtin(output)
	
	//Startup the scanner.
	var s scanner.Scanner
	s.Init(file)
	s.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
	
	//TODO cleanup file from here forward.
	var tok rune
	for tok != scanner.EOF {
		tok = s.Scan()
		
		switch s.TokenText() {
			case "\n", ";":
			
			case "!":
				output.Write([]byte("ERROR 0\n"))
			
			case "}", "end":
				output.Write([]byte("END\n"))
				LoseScope()
			
			case "else":
				fmt.Fprintf(output, "ELSE\n")
				LoseScope()
				GainScope()
			
			case "if", "elseif":
				
				if s.TokenText() == "if" {
					GainScope()
					s.Scan()
					fmt.Fprintf(output, "IF %v\n", expression(&s, output))
				} else {
					LoseScope()
					GainScope()
					s.Scan()
					fmt.Fprintf(output, "ELSEIF %v\n", expression(&s, output))
				}
				
				
			//Inline universal assembly.
			case ".":
				s.Scan()
				output.Write([]byte(strings.ToUpper(s.TokenText())))
				for tok = s.Scan(); tok != scanner.EOF; {
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
					if s.TokenText() != "\n" {
						if len(CurrentFunction.Returns) > 0 {
							switch CurrentFunction.Returns[0] {
								case NUMBER:
									output.Write([]byte("PUSH "+expression(&s, output)+"\n"))
								case STRING:
									output.Write([]byte("PUSHSTRING "+expression(&s, output)+"\n"))
								case FUNCTION:
									output.Write([]byte("PUSHFUNC "+expression(&s, output)+"\n"))
								case FILE:
									output.Write([]byte("PUSHIT "+expression(&s, output)+"\n"))
							}
						}
					}
				}
				output.Write([]byte("RETURN\n"))
			
			case "software":
				output.Write([]byte("ROUTINE\n"))
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
			
			case "issues":
				output.Write([]byte("IF ERROR\nADD ERROR 0 0\n"))
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
				
			//Compiles function declerations.
			case "function", "method":
				GainScope()
				var name string
				var function Function
				
				var method bool = s.TokenText() == "method"
				var methodType = NUMBER
				
				
				// function name(param1, param2) returns {
				output.Write([]byte("SUBROUTINE "))
				s.Scan()
				if method {
					switch s.TokenText() {
						case "~":
							methodType = FILE
							s.Scan()
					}
				}
				if method {
					fmt.Fprintf(output, "%v_m_%v\n", s.TokenText(), methodType)
				} else {
					output.Write([]byte(s.TokenText()+"\n"))
				}
				name = s.TokenText()
				s.Scan()
				if s.TokenText() != "(" {
					fmt.Println(s.Pos(), "Expecting ( found ", s.TokenText())
					return
				}
				
				//We need to reverse the POP's because of stack pain.
				var toReverse []string
				for tok = s.Scan(); tok != scanner.EOF; {
					var popstring string
					if s.TokenText() == ")" {
						break
					}
					//String arguments.
					if s.TokenText() == "[" {
						//Update our function definition with a string argument.
						function.Args = append(function.Args, STRING)
						
						popstring += "POPSTRING "
						s.Scan()
						if s.TokenText() != "]" {
							fmt.Println(s.Pos(), "Expecting ] found ", s.TokenText())
							return
						}
						s.Scan()
					//Other type of string argument. (Variadic)
					} else if s.TokenText() == "." {
						
						//Update our function definition with a string argument.
						function.Args = append(function.Args, STRING)
						function.Variadic = true
						
						popstring += "POPSTRING "
						s.Scan()
						if s.TokenText() != "." {
							fmt.Println(s.Pos(), "Expecting . found ", s.TokenText())
							return
						}
						s.Scan()
					//Function arguments.
					} else if s.TokenText() == "(" {
						
						//Update our function definition with a string argument.
						function.Args = append(function.Args, FUNCTION)
						
						popstring += "POPFUNC "
						s.Scan()
						if s.TokenText() != ")" {
							fmt.Println(s.Pos(), "Expecting ) found ", s.TokenText())
							return
						}
						s.Scan()
					} else {
						//Update our function definition with a numeric argument.
						function.Args = append(function.Args, NUMBER)
						
						popstring += "POP "
					}
					popstring += s.TokenText()+"\n"
					toReverse = append(toReverse, popstring)
					s.Scan()
					if s.TokenText() == ")" {
						break
					}
					if s.TokenText() != "," {
						fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
						return
					}
					s.Scan()
				}
				for i := len(toReverse)-1; i>=0; i-- {
					output.Write([]byte(toReverse[i]))
				}
				if method {
					switch methodType {
						case FILE:
							fmt.Fprintf(output, "POPIT self\n")
					}
				}
				s.Scan()
				if s.TokenText() != "{" {
					if s.TokenText() != "[" {
						function.Returns = append(function.Returns, NUMBER)
					} else {
						function.Returns = append(function.Returns, STRING)
						s.Scan()
						if s.TokenText() != "]" {
							fmt.Println(s.Pos(), "Expecting ] found ", s.TokenText())
							return
						}
					}
					s.Scan()
					if s.TokenText() != "{" {	
						fmt.Println(s.Pos(), "Expecting { found ", s.TokenText())
						return
					}
				}
				s.Scan()
				if s.TokenText() != "\n" {
					fmt.Println(s.Pos(), "Expecting newline found ", s.TokenText())
					return
				}
				
				function.Exists = true
				
				if method {
					methods[name] = true
					functions[name+"_m_"+fmt.Sprint(methodType)] = function	
				} else {
					functions[name] = function
				}
				
				CurrentFunction = function
			case "var":
				s.Scan()
				if s.TokenText() == "[" {
					s.Scan()
					if s.TokenText() != "]" {
						fmt.Println(s.Pos(), "Expecting ] found ", s.TokenText())
						return
					}
					s.Scan()
					var name = s.TokenText()
					s.Scan()
					s.Scan()
					n := expression(&s, output)
					fmt.Fprintf(output, "PUSHSTRING %v\n", n)
					fmt.Fprintf(output, "POPSTRING %v\n", name)
					SetVariable(name, STRING)
				} else if s.TokenText() == "(" {
					s.Scan()
					if s.TokenText() != ")" {
						fmt.Println(s.Pos(), "Expecting ) found ", s.TokenText())
						return
					}
					s.Scan()
					var name = s.TokenText()
					s.Scan()
					s.Scan()
					n := expression(&s, output)
					fmt.Fprintf(output, "PUSHFUNC %v\n", n)
					fmt.Fprintf(output, "POPFUNC %v\n", name)
					SetVariable(name, FUNCTION)
				} else if s.TokenText() == "~" {
					s.Scan()
					var name = s.TokenText()
					s.Scan()
					s.Scan()
					n := expression(&s, output)
					fmt.Fprintf(output, "PUSHIT %v\n", n)
					fmt.Fprintf(output, "POPIT %v\n", name)
					SetVariable(name, FILE)
				} else {
					var name = s.TokenText()
					s.Scan()
					s.Scan()
					fmt.Fprintf(output, "VAR %v %v\n", name, expression(&s, output))
					SetVariable(name, NUMBER)
				}
			
			default:
			
				var name = s.TokenText()
				
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
						output.Write([]byte("PUSH "+expression(&s, output)+" "+name+" \n"))
					case "=":
						/*if functions[s.TokenText()].Exists && s.Peek() != '(' {
							
							functions[name] = functions[s.TokenText()]
							f := functions[name] 
							f.Local = true
							functions[name] = f
							output.Write([]byte("FUNC "+name+" "+s.TokenText()+"\n"))
							
						}*/
						s.Scan()
						fmt.Fprintf(output, "ADD %v 0 %v\n", name, expression(&s, output))
					default:
						if len(s.TokenText()) > 0 && s.TokenText()[0] == '.' {
							var index = s.TokenText()[1:]
							s.Scan()
							if s.TokenText() != "=" {
								fmt.Println(s.Pos(), "Expecting = found ", s.TokenText())
								return
							}
							s.Scan()
							output.Write([]byte("SET "+name+" "+index+" "+expression(&s, output)+"\n"))
							
						} else {
					
					
							if name == "" {
								return	
							}
							fmt.Println(s.Pos(), "Unexpected ", name)
							return
						}
				}
				
		}
	}
}
