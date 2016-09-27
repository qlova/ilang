package main

import (
	"text/scanner"
	"fmt"
	"io"
	"os"
	"strings"
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


var loaded = make(map[string]bool)

var CurrentFunctionName string

func LoadFunction(name string) {
	if functions[name].Data != "" && !loaded[name] && !functions[name].Inline {
		IFILE.Write([]byte(functions[name].Data))
	
		loaded[name] = true
	}
}

func ParseConstructDef(s *scanner.Scanner, output io.Writer) {
	name := NextToken(s, Iname)
	fmt.Fprintf(output, "FUNCTION %s\n", name)
	fmt.Fprintf(output, "ARRAY %v\n", LastDefinedType.String())
	
	for range DefinedTypes[LastDefinedType-USER].Elements {
		fmt.Fprintf(output, "PUT 0\n")
	}
	
	SetVariable(LastDefinedType.String(), LastDefinedType)
	SetVariable(LastDefinedType.String()+".", 1)
	
	var function Function
	function.Exists = true
	function.Returns = []TYPE{LastDefinedType}
	functions[name] = function
	CurrentFunction = function

	s.Scan()	
	Expecting(s, "{")
	s.Scan()
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
		
		//Identfy the type and add it to the function.
		var ArgumentType = ParseSymbolicType(s)
		
		if ArgumentType == MULTIPLE {
			function.Variadic = true
			ArgumentType = ARRAY
		}
		function.Args = append(function.Args, ArgumentType)
		popstring += ArgumentType.Pop()+" "
		
		SetVariable(s.TokenText(), ArgumentType)
		
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
		SetVariable(name+".", 1)
	}
	
	
	s.Scan()
	
	//Find out the return value.
	if s.TokenText() != "{" || (s.TokenText() == "{" && s.Peek() == '}') {
		var ReturnType = ParseSymbolicType(s)
		function.Returns = append(function.Returns, ReturnType)
		if ReturnType == NUMBER {
			s.Scan()
		}
	}
	Expecting(s, "{")
	
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
	CurrentFunctionName = name
}

//Parse the return value for a function.
func ParseFunctionReturns(token string, s *scanner.Scanner, output io.Writer, shunting bool) string {	
	if len(functions[token].Returns) > 0 {
		unique++
		id := "i+output+"+fmt.Sprint(unique)
		
		var ReturnType = functions[token].Returns[0]
		
		fmt.Fprintf(output, "%s %v\n", ReturnType.Pop(), id)
		ExpressionType = ReturnType
		
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
	var noargs = len(calling) > 2 && calling[2]

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
		
		for !noargs {
			
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
				
				//PUSH type
				fmt.Fprintf(output, "%s %v\n", sort.Push(), name)

				s.Scan()
				token = token+"_m_"+methodType.String()
			}
		
			if s.TokenText() == ")" {
				if len(functions[token].Args) > i {
					RaiseError(s, token+" requires "+fmt.Sprint(len(functions[token].Args))+" arguments!")
				}
				return token
			}
			
			if s.TokenText() == "," && !(len(functions[token].Args) > i) {
				break
			}
			
			s.Scan()
			if s.TokenText() == ")" {
				if len(functions[token].Args) > i {
					RaiseError(s, token+" requires "+fmt.Sprint(len(functions[token].Args))+" arguments!")
				}
				break
			}
						
			if len(functions[token].Args) < i {
				RaiseError(s, token+" requires "+fmt.Sprint(len(functions[token].Args))+" arguments!")
			}
		
			if methods[token] || len(functions[token].Args) > i {
				var argument = expression(s, output)
				
				if methods[token] && ExpressionType != UNDEFINED && len(functions[token].Args) == 0 {
					if methods[token] {
					
						//Special something calls.
						if ExpressionType == SOMETHING {
							
							//s.Scan()
							fmt.Fprintf(output, "PLACE %v\n", argument)
							fmt.Fprintf(output, "PUSH 1\nGET itype\nIF 1\nVAR itypetest\nIF itypetest\nERROR 1\n")
							var ends = 0
							for key, _ := range functions {
								split := strings.Split(key, "_m_")
								if len(split) > 1 && split[0] == token {
									fmt.Fprintf(output, "ELSE\nSEQ itypetest itype %v\nIF itypetest\n", int(StringToType[split[1]]))
									switch StringToType[split[1]] {
										case NUMBER, LETTER, ITYPE:
											fmt.Fprintf(output,"PUSH 0\nGET data\nPUSH data\n")
										case STRING, ARRAY:
											unique++
											fmt.Fprintf(output, "PUSH %v\n", 0)
											fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				
											fmt.Fprintf(output, "PUSH %s%v\n", "i+user+", unique)
											fmt.Fprintf(output, "HEAP\n")
										default:
											if ExpressionType >= USER {
												unique++
												fmt.Fprintf(output, "PUSH %v\n", 0)
												fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				
												fmt.Fprintf(output, "PUSH %s%v\n", "i+user+", unique)
												fmt.Fprintf(output, "HEAP\n")
											}
									}
									if functions[key].Inline {
										output.Write([]byte(functions[key].Data+"\n"))
									} else {
										//println(key, functions[key].Exists)
										output.Write([]byte("RUN "+key+"\n"))
									}
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
						
						//PUSH type.
						fmt.Fprintf(output, "%s %v\n", ExpressionType.Push(), argument)
						
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
			i++
		
			if s.TokenText() == ")" {
				if len(functions[token].Args) > i {
					RaiseError(s, token+" requires "+fmt.Sprint(len(functions[token].Args))+" arguments!")
				}
				break
			}
			if s.TokenText() != "," {
				fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
				os.Exit(1)
			}
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
		LoadFunction(functions[token].Load)
		
		if !noreturns {
			return ParseFunctionReturns(token, s, output, shunting)
		}
	}
	return ""
}
