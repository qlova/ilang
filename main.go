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

type Function struct {
	Exists bool
	Args []bool
	Returns []bool
}

var variables = make( map[string]bool)
var functions = make( map[string]Function)
var unique int

func shunt(name string, s *scanner.Scanner, output io.Writer) string {
		s.Scan()
		switch s.TokenText() {
			case ")", ",", "\n", "]":
				return name
			case "/":
				unique++
				output.Write([]byte("VAR i+shunt+"+fmt.Sprint(unique)+"\n"))
				s.Scan()
				output.Write([]byte("DIV i+shunt+"+fmt.Sprint(unique)+" "+name+" "+expression(s, output)+"\n"))
				return "i+shunt+"+fmt.Sprint(unique)
			case "+":
				unique++
				output.Write([]byte("VAR i+shunt+"+fmt.Sprint(unique)+"\n"))
				s.Scan()
				output.Write([]byte("ADD i+shunt+"+fmt.Sprint(unique)+" "+name+" "+expression(s, output)+"\n"))
				return "i+shunt+"+fmt.Sprint(unique)
			case "²":
				unique++
				output.Write([]byte("VAR i+shunt+"+fmt.Sprint(unique)+"\n"))
				s.Scan()
				output.Write([]byte("MUL i+shunt+"+fmt.Sprint(unique)+" "+name+" "+name+"\n"))
				return "i+shunt+"+fmt.Sprint(unique)
			default:
				println(name, s.TokenText())
			
		}
		return ""
}

func expression(s *scanner.Scanner, output io.Writer) string {

	//Turn string literals into numeric strings.
	//For example string arguments to a function
	//eg. output("A")
	// ->
	// STRING i+tmp+id
	// PUSH 'A' i+tmp+id
	// PUSHSTRING i+tmp+id
	// RUN output
	if s.TokenText()[0] == '"' {
				
		unique++
		var newarg string = "STRING i+tmp+"+fmt.Sprint(unique)+"\n"
		var j int
		var arg = s.TokenText()[1:]
		
		stringloop:
		arg = strings.Replace(arg, "\\n", "\n", -1)
		for _, v := range arg {
			if v == '"' {
				goto end
			}
			newarg += "PUSH "+strconv.Itoa(int(v))+" i+tmp+"+fmt.Sprint(unique)+"\n"
		}
		if len(arg) == 0 {
			goto end
		}
		newarg += "PUSH "+strconv.Itoa(int(' '))+" i+tmp+"+fmt.Sprint(unique)+"\n"
		j++
		//println(arg)
		arg = string(s.TokenText()[j])
		goto stringloop
		end:
		//println(newarg)
		output.Write([]byte(newarg))
		s.Scan()
		return "i+tmp+"+fmt.Sprint(unique)
	}


	//Is it a literal number?
	if _, err := strconv.Atoi(s.TokenText()); err == nil {
		return shunt(s.TokenText(), s, output)
	} else {
	
		var name = s.TokenText()
	
		if functions[name].Exists  {

			var i int
			for tok := s.Scan(); tok != scanner.EOF; {
				
				s.Scan()
				if s.TokenText() == ")" {
					break
				}
				
				if len(functions[name].Args) > i {
					if functions[name].Args[i] {
						output.Write([]byte("PUSHSTRING "+expression(s, output)+"\n"))
					} else {
						output.Write([]byte("PUSH "+expression(s, output)+"\n"))
					}
				}
				
				if s.TokenText() == ")" {
					break
				}
				if s.TokenText() != "," {
					fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
					os.Exit(1)
				}
				
				
			}
			s.Scan()		
			unique++
			output.Write([]byte("RUN "+name+"\n"))
			if len(functions[name].Returns) > 0 {
				output.Write([]byte("POP i+output+"+fmt.Sprint(unique)+"\n"))
			}			
			return "i+output+"+fmt.Sprint(unique)
		}
	
		//Is it a variable?
		if variables[s.TokenText()] {
			return shunt(s.TokenText(), s, output)
			
		} else {
			
			// a=2; b=4; ab
			if variables[string(rune(s.TokenText()[0]))] {
				if len(s.TokenText()) == 2 {
					if variables[string(rune(s.TokenText()[1]))] {
						unique++
						output.Write([]byte("VAR i+tmp+"+s.TokenText()+fmt.Sprint(unique)+"\n"))
						output.Write([]byte("MUL i+tmp+"+s.TokenText()+fmt.Sprint(unique)+" "+
							string(rune(s.TokenText()[0]))+" "+
							string(rune(s.TokenText()[1]))+"\n"))
						
						return shunt("i+tmp+"+s.TokenText()+fmt.Sprint(unique), s, output)
					}
				}
			}
			
		}
	
	}
	return shunt(s.TokenText(), s, output)
}

func main() {
	flag.Parse()

	file, err := os.Open(flag.Arg(0))
	if err != nil {
		return
	}
	
	output, err := os.Create(flag.Arg(0)+".u")
	if err != nil {
		return
	}
	
	output.Write([]byte(
`
SUBROUTINE output
	POPSTRING data
	PUSHSTRING data
	STDOUT 
END
`	))
	functions["output"] = Function{Exists:true, Args:[]bool{true}}
	
	var s scanner.Scanner
	s.Init(file)
	s.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
	
	var tok rune
	for tok != scanner.EOF {
		tok = s.Scan()
		
		switch s.TokenText() {
			case "\n":
				
			
			case "}":
				output.Write([]byte("END\n"))
				
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
				output.Write([]byte("PUSH "+expression(&s, output)+"\n"))
			
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
					if s.TokenText() == "[" {
						
						//Update our function definition with a string argument.
						function.Args = append(function.Args, true)
						
						popstring += "POPSTRING "
						s.Scan()
						if s.TokenText() != "]" {
							fmt.Println(s.Pos(), "Expecting ] found ", s.TokenText())
							return
						}
						s.Scan()
					} else {
						//Update our function definition with a numeric argument.
						function.Args = append(function.Args, false)
						
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
						function.Returns = append(function.Returns, false)
					} else {
						function.Returns = append(function.Returns, true)
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
					expression(&s, output)
					continue
				}
				
				s.Scan()
				switch s.TokenText() {
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
							
						} else {
							variables[name] = true
							output.Write([]byte("VAR "+name+" "+s.TokenText()+"\n"))
						}
					default:
						fmt.Println(s.Pos(), "Unexpected ", name)
						return
				}
				
		}
	}
}
