package main

import (
	"text/scanner"
	"fmt"
	"io"
	"strconv"
	"strings"
	"os"
)

/*Turn string literals into numeric strings.
//For example string arguments to a function
//eg. output("A") becomes:

	STRING i+tmp+id
	PUSH 'A' i+tmp+id
	PUSHSTRING i+tmp+id
	
	RUN output.
*/
func ParseString(s *scanner.Scanner, output io.Writer, shunting bool) string {
	unique++
	id := "i+tmp+"+fmt.Sprint(unique)
	
	fmt.Fprintf(output, "STRING %v\n", id)
	
	//This is such a mess :/ 
	var j int
	var arg = s.TokenText()[1:]
	
	//Didn't my master always say 
	//		"NO GOTO STATEMENTS, YOUR MAKING ME BLIND!"
	stringloop:
	arg = strings.Replace(arg, "\\n", "\n", -1)
	for _, v := range arg {
		if v == '"' {
			goto end
		}
		fmt.Fprintf(output, "PUSH %v %v\n", strconv.Itoa(int(v)), id)
	}
	if len(arg) == 0 {
		goto end
	}
	fmt.Fprintf(output, "PUSH %v %v\n", strconv.Itoa(int(' ')), id)
	j++
	arg = string(s.TokenText()[j])
	goto stringloop
	end:
	
	//Return the identifier of the string.
	if shunting {
		return shunt(id, s, output)
	} else {
		return id
	}
}

/*
	Parse array literals.
	For example:
	
		output([1, 2, 3])
	
	becomes:
		
		STRING i+string+id
		PUSH 1 i+string+id
		PUSH 2 i+string+id
		PUSH 3 i+string+id
		
		PUSHSTRING i+string+id
		RUN output
*/
func ParseArray(s *scanner.Scanner, output io.Writer, shunting bool) string {
	unique++
	var id = "i+string+"+fmt.Sprint(unique)

	fmt.Fprintf(output, "STRING %v\n", id)

	//Push all the values.
	for tok := s.Scan(); tok != scanner.EOF; {

		if s.TokenText() == "]" {
			break
		}

		fmt.Fprintf(output, "PUSH %v %v\n", expression(s, output), id)
	
		if s.TokenText() == "]" {
			break
		}
	
		if s.TokenText() != "," {
			fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
			os.Exit(1)
		}
		s.Scan()
	}
	if shunting {
		return shunt(id, s, output)
	} else {
		return id
	}
}

//Parse the return value for a function.
func ParseFunctionReturns(token string, s *scanner.Scanner, output io.Writer, shunting bool) string {	
	if len(functions[token].Returns) > 0 {
		unique++
		id := "i+output+"+fmt.Sprint(unique)
		switch functions[token].Returns[0] {
			case STRING:
				fmt.Fprintf(output, "POPSTRING %v\n", id)
			case NUMBER:
				fmt.Fprintf(output, "POP %v\n", id)
			case FUNCTION:
				fmt.Fprintf(output, "POPFUNC %v\n", id)
			case FILE:
				fmt.Fprintf(output, "POPIT %v\n", id)
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
func ParseFunction(s *scanner.Scanner, output io.Writer, shunting bool) string {
	var token = s.TokenText()
	
	//Currently variadic functions only work with numbers. Why? No reason (Lazyness).
	if functions[token].Variadic {
		unique++
		id := "i+output+"+fmt.Sprint(unique)
		
		fmt.Fprintf(output, "STRING %v\n", id)
		for tok := s.Scan(); tok != scanner.EOF; {
			
			if s.TokenText() == ")" {
				break
			}
			s.Scan()
		
			fmt.Fprintf(output, "PUSH %v %v\n", expression(s, output), id)
			
			if s.TokenText() == ")" {
				break
			}
			
			if s.TokenText() != "," {
				fmt.Println(s.Pos(), "Expecting , found ", s.TokenText())
				os.Exit(1)
			}
		}
	
		fmt.Fprintf(output, "PUSHSTRING %v\n", id)
	} else {

		var i int
		for tok := s.Scan(); tok != scanner.EOF; {
			if s.TokenText() == ")" {
				return token
			}
			s.Scan()
			if s.TokenText() == ")" {
				break
			}
		
			if len(functions[token].Args) > i {
				switch functions[token].Args[i] {
					case STRING:
						fmt.Fprintf(output, "PUSHSTRING %v\n", expression(s, output))
					case NUMBER:
						fmt.Fprintf(output, "PUSH %v\n", expression(s, output))
					case FUNCTION:
						fmt.Fprintf(output, "PUSHFUNC %v\n", expression(s, output))
					case FILE:
						fmt.Fprintf(output, "PUSHIT %v\n", expression(s, output))
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
	}
		
	if functions[token].Local {
		output.Write([]byte("EXE "+token+"\n"))
	} else {
		output.Write([]byte("RUN "+token+"\n"))
	}
	return ParseFunctionReturns(token, s, output, shunting)
}
