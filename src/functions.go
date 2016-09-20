package main

import (
	"text/scanner"
	"fmt"
	"io"
	"os"
	"strings"
)

var loaded = make(map[string]bool)

func LoadFunction(name string) {
	if functions[name].Data != "" && !loaded[name] && !functions[name].Inline {
		IFILE.Write([]byte(functions[name].Data))
	
		loaded[name] = true
	}
}

func ParseFunctionDef(s *scanner.Scanner, output io.Writer) {
	var name string
	var function Function
	
	var method bool = s.TokenText() == "method"
	var methodType = LastDefinedType
	
	
	// function name(param1, param2) returns {
	output.Write([]byte("FUNCTION "))
	s.Scan()
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
	for tok := s.Scan(); tok != scanner.EOF; {
		var popstring string
		if s.TokenText() == ")" {
			break
		}
		var T TYPE
		//String arguments.
		if s.TokenText() == "[" {
			//Update our function definition with a string argument.
			function.Args = append(function.Args, STRING)
			
			popstring += "GRAB "
			
			T = STRING
			
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
			
			popstring += "GRAB "
			
			T = STRING
			
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
			
			T = FUNCTION
			
			popstring += "TAKE "
			s.Scan()
			if s.TokenText() != ")" {
				fmt.Println(s.Pos(), "Expecting ) found ", s.TokenText())
				return
			}
			s.Scan()
		//File arguments.
		} else if s.TokenText() == "|" {
			
			//Update our function definition with a string argument.
			function.Args = append(function.Args, FILE)
			
			T = FILE
			
			popstring += "TAKE "
			s.Scan()
		} else {
			//Update our function definition with a numeric argument.
			function.Args = append(function.Args, NUMBER)
			
			T = NUMBER
			
			popstring += "PULL "
		}
		SetVariable(s.TokenText(), T)
		
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
	
	//The method variable needs to be popped and put into scope.
	if method {
		name := LastDefinedType.String()
		fmt.Fprintf(output, "GRAB %s\n", name)
		SetVariable(name, LastDefinedType)
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
}

//Parse the return value for a function.
func ParseFunctionReturns(token string, s *scanner.Scanner, output io.Writer, shunting bool) string {	
	if len(functions[token].Returns) > 0 {
		unique++
		id := "i+output+"+fmt.Sprint(unique)
		switch functions[token].Returns[0] {
			case STRING:
				fmt.Fprintf(output, "GRAB %v\n", id)
				ExpressionType = STRING
			case NUMBER:
				fmt.Fprintf(output, "PULL %v\n", id)
				ExpressionType = NUMBER
			case FUNCTION:
				fmt.Fprintf(output, "TAKE %v\n", id)
				ExpressionType = FUNCTION
			case FILE:
				fmt.Fprintf(output, "TAKE %v\n", id)
				ExpressionType = FILE
		}
		
		if shunting {
			return shunt(id, s, output)
		}	
		return id
	}	
	
	return ""	
}

/*
	Parse a function.
	eg.	output(text(20)) becomes:
	
		PUSH 20
		RUN text
		POPSTRING i+output+id
		
		PUSHSTRING i+output+id
		RUN output
*/
func ParseFunction(name string, s *scanner.Scanner, output io.Writer, shunting bool, calling ...bool) string {
	var token = name

	var methodType TYPE
	
	var call = len(calling) == 0 || calling[0]
	var noreturns = len(calling) > 1 && calling[1]

	//Currently variadic functions only work with numbers. Why? No reason (Lazyness).
	if functions[token].Variadic {
		unique++
		id := "i+output+"+fmt.Sprint(unique)
		
		fmt.Fprintf(output, "ARRAY %v\n", id)
		for tok := s.Scan(); tok != scanner.EOF; {
			
			if s.TokenText() == ")" {
				break
			}
			s.Scan()
		
			fmt.Fprintf(output, "PLACE %v\nPUT %v\n", id, expression(s, output))
			
			if s.TokenText() == ")" {
				break
			}
			
			if s.TokenText() != "," {
				fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
				os.Exit(1)
			}
		}
	
		fmt.Fprintf(output, "SHARE %v\n", id)
	} else {

		var i int
		if s.TokenText() != "," {
			s.Scan()
		}
		for {
			
			if s.TokenText() == "@" {
				s.Scan()
				var name string
				var sort TYPE
				if s.TokenText() == "(" {
					name = expression(s, output, false)
					sort = ExpressionType
				} else {
					name = s.TokenText()
					sort = GetVariable(s.TokenText())
				}
				methodType = sort
				if sort == UNDEFINED && (
					name == "inbox" || name == "outbox" ){
						methodType = FILE
				}
				switch sort {
					case STRING:
						fmt.Fprintf(output, "SHARE %v\n", name)
					case NUMBER:
						fmt.Fprintf(output, "PUSH %v\n", name)
					case FUNCTION:
						fmt.Fprintf(output, "RELAY %v\n", name)
					case FILE:
						fmt.Fprintf(output, "RELAY %v\n", name)
						
					default:
						fmt.Fprintf(output, "SHARE %v\n", name)
				}
				s.Scan()
				token = token+"_m_"+methodType.String()
			}
		
			if s.TokenText() == ")" {
				return token
			}
			
			if s.TokenText() == "," && !(len(functions[token].Args) > i) {
				break
			}
			
			s.Scan()
			if s.TokenText() == ")" {
				break
			}
		
			if len(functions[token].Args) > i {
				var argument = expression(s, output)
				
				if ExpressionType != functions[token].Args[i] {
					if methods[token] {
					
						//Special something calls.
						if ExpressionType == SOMETHING {
							
							//s.Scan()
							fmt.Fprintf(output, "PLACE %v\n", argument)
							fmt.Fprintf(output, "PUSH 1\nGET itype\nIF 1\nVAR itypetest\nIF itypetest\nERROR 1\n")
							var ends = 0
							for key, f := range functions {
								split := strings.Split(key, "_m_")
								if len(f.Args) == 1 && len(split) > 0 && split[0] == token {
									fmt.Fprintf(output, "ELSE\nSEQ itypetest itype %v\nIF itypetest\n", int(f.Args[0]))
									switch f.Args[0] {
										case NUMBER:
											fmt.Fprintf(output,"PUSH 0\nGET data\nPUSH data\n")
									}
									fmt.Fprintf(output,"RUN %s\n", key)
									ends++
								}
							}
							
							for i := 0; i < ends; i ++ {
								fmt.Fprintf(output, "END\n")
							}
							fmt.Fprintf(output, "END\nEND\n")
							if call && !noreturns {
								return ParseFunctionReturns(token, s, output, shunting)
							}
							return ""
						}
					
						token = token+"_m_"+ExpressionType.String()
						
						if !functions[token].Exists || len(functions[token].Args) != 0 {
							RaiseError(s, fmt.Sprintf("%v is not defined for arguments of type %v!", 
								token, ExpressionType.String()))
						}
						
						fmt.Fprintf(output, "SHARE %v\n", argument)
						goto endTypeCheck
						
					} else {	
						RaiseError(s, fmt.Sprintf("Type mismatch! Argument %v of '%v()' expects %v, got %v (%v)", 
							fmt.Sprint(i+1), token, functions[token].Args[i].String(), ExpressionType.String(), argument))
					}
				}
				fmt.Fprintf(output, "%v %v\n", functions[token].Args[i].Push(), argument)
			} else {
				break
			}
			endTypeCheck:
		
			if s.TokenText() == ")" {
				break
			}
			if s.TokenText() != "," {
				fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
				os.Exit(1)
			}
			i++
		}	
	}
	
	if call && !functions[token].Ghost {
		if functions[token].Local {
			output.Write([]byte("EXE "+token+"\n"))
		} else if functions[token].Inline {
			 output.Write([]byte(functions[token].Data+"\n"))
		} else {
			output.Write([]byte("RUN "+token+"\n"))
		}
		
		//Write the function to the ifile if it is builtin.
		//println(token, functions[token].Data != "", !functions[token].Loaded, !functions[token].Inline)
		LoadFunction(token)
		
		if !noreturns {
			return ParseFunctionReturns(token, s, output, shunting)
		}
	}
	return ""
}
