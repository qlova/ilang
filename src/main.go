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

var variables = make( map[string]bool)
var functions = make( map[string]Function)
var unique int

func expression(s *scanner.Scanner, output io.Writer, param ...bool) string {
	
	//Do we need to shunt? This is for operator precidence. (defaults to true)
	var shunting bool = len(param) <= 0 || param[0]
	var token = s.TokenText()
	
	if len(token) <= 0 {
		fmt.Println("Empty expression!")
		return ""
	}

	//If there is a quotation mark, parse the string. 
	if token[0] == '"' {
		return ParseString(s, output, shunting)
	}
	
	if  token[0] == '[' {
		return ParseArray(s, output, shunting)
	}
	
	//Deal with runes. 
	//	eg. 'a' -> 97
	if len(token) == 3 && token[0] == '\'' && token[2] == '\'' {
		defer s.Scan()
		return strconv.Itoa(int(s.TokenText()[1]))
	} else if s.TokenText() == `'\n'` {
		defer s.Scan()
		return strconv.Itoa(int('\n'))
	}

	
	//Parse function call.
	if functions[token].Exists {
		return ParseFunction(s, output, shunting)
	}

	//Is it a literal number? Then just return it.
	//OR is it a variable?
	if _, err := strconv.Atoi(token); err == nil || variables[token] {
		if shunting {
			return shunt(token, s, output)
		} else {
			return token
		}
	}
		
	//Do some special maths which people will hate about I.
	// a=2; b=4; ab
	// ab is 8
	if variables[string(rune(token[0]))] {
		if len(token) == 2 {
			if variables[string(rune(token[1]))] {
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
				
			
			case "}", "end":
				output.Write([]byte("END\n"))
			
			case "else":
				fmt.Fprintf(output, "ELSE\n")
			
			case "if", "elseif":
				
				if s.TokenText() == "if" {
					s.Scan()
					fmt.Fprintf(output, "IF %v\n", expression(&s, output))
				} else {
					s.Scan()
					fmt.Fprintf(output, "ELSEIF %v\n", expression(&s, output))
				}
				
			//Inline universal assembly.
			case ".":
				s.Scan()
				output.Write([]byte(strings.ToUpper(s.TokenText()+" ")))
				for tok = s.Scan(); tok != scanner.EOF; {
					if s.TokenText() == "\n" {
						output.Write([]byte("\n"))
						break
					}
					output.Write([]byte(s.TokenText()))
					s.Scan()
				}
			
			case "return":
				s.Scan()
				if s.TokenText() != "\n" {
					output.Write([]byte("PUSH "+expression(&s, output)+"\n"))
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
				
			//Compiles function declerations.
			case "function":
				var name string
				var function Function
				
				// function name(param1, param2) returns {
				output.Write([]byte("SUBROUTINE "))
				s.Scan()
				output.Write([]byte(s.TokenText()+"\n"))
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
				functions[name] = function
			default:
			
				var name = s.TokenText()
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
						variables[name] = true
						output.Write([]byte("PUSH "+expression(&s, output)+" "+name+" \n"))
					case "=":
						// a = 
						s.Scan()
						if s.TokenText() == "[" {
							//a = [12,32,92]
							output.Write([]byte("STRING "+name+"\n"))
							
							for tok = s.Scan(); tok != scanner.EOF; {
							
								if s.TokenText() == "]" {
									break
								}
							
								output.Write([]byte("PUSH "+expression(&s, output)+" "+name+"\n"))
								
								if s.TokenText() == "]" {
									break
								}
								
								if s.TokenText() != "," {
									fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
									return
								}
								s.Scan()
							}
						
						} else if s.TokenText()[0] == '"' {
							//Turn string literals into numeric strings.
							//For example string arguments to a function
							//eg. output("A")
							// ->
							// STRING i+tmp+id
							// PUSH 'A' i+tmp+id
							// PUSHSTRING i+tmp+id
							// RUN output
								var newarg string = "STRING "+name+"\n"
								var j int
								var arg = s.TokenText()[1:]
		
								stringloop:
								arg = strings.Replace(arg, "\\n", "\n", -1)
								for _, v := range arg {
									if v == '"' {
										goto end
									}
									newarg += "PUSH "+strconv.Itoa(int(v))+" "+name+"\n"
								}
								if len(arg) == 0 {
									goto end
								}
								newarg += "PUSH "+strconv.Itoa(int(' '))+" "+name+"\n"
								j++
								//println(arg)
								arg = string(s.TokenText()[j])
								goto stringloop
								end:
								//println(newarg)
								output.Write([]byte(newarg))
								s.Scan()
						
						} else {
							if functions[s.TokenText()].Exists && s.Peek() != '(' {
								
								functions[name] = functions[s.TokenText()]
								f := functions[name] 
								f.Local = true
								functions[name] = f
								output.Write([]byte("FUNC "+name+" "+s.TokenText()+"\n"))
								
							} else if functions[s.TokenText()].Exists {
								
								if len(functions[s.TokenText()].Returns) > 0 {
									if functions[s.TokenText()].Returns[0] == FILE {
										variables[name] = true
										fmt.Fprintf(output, "PUSHIT %v\n", expression(&s, output))
										fmt.Fprintf(output, "POPIT %v\n", name)
									} else if functions[s.TokenText()].Returns[0] == STRING {
										variables[name] = true
										fmt.Fprintf(output, "PUSHSTRING %v\n", expression(&s, output))
										fmt.Fprintf(output, "POPSTRING %v\n", name)
									} else if functions[s.TokenText()].Returns[0] == FUNCTION {
										variables[name] = true
										fmt.Fprintf(output, "PUSHFUNC %v\n", expression(&s, output))
										fmt.Fprintf(output, "POPFUNCTION %v\n", name)
									} else {
										variables[name] = true
										output.Write([]byte("VAR "+name+" "+expression(&s, output)+"\n"))
									}
								} else {
									fmt.Println(s.Pos(), "Function ", s.TokenText(), " output cannot be assigned to a value!")
									return
								}
								
							} else {
						
								variables[name] = true
								output.Write([]byte("VAR "+name+" "+expression(&s, output)+"\n"))
							}
						}
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
