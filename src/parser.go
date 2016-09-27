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
	
	fmt.Fprintf(output, "ARRAY %v\n", id)
	
	//This is such a mess :/ 
	var j int
	var arg = s.TokenText()[1:]
	arg = strings.Replace(arg, "\\n", "\n", -1)
	arg = strings.Replace(arg, "\\r", "\r", -1)
	arg = strings.Replace(arg, "\\\n", "\\n", -1)
	
	var backflipping = false
	var backflip string
	
	//Didn't my master always say 
	//		"NO GOTO STATEMENTS, YOUR MAKING ME BLIND!"
	stringloop:
	for _, v := range arg {
		if v == '"' {
			goto end
		}
		
		if !backflipping && v == '`' {
			backflipping = true
			backflip  = ""
			continue
		}
		
		if backflipping {
			
			if v != '`' {
				backflip += string(v)
			} else {
				backflipping = false
				
				anotherlabel:
				
				if GetVariable(backflip) == STRING {
					fmt.Fprintf(output, "PLACE %v\n", id)
					fmt.Fprintf(output, "JOIN %v %v %v\n", id, id, backflip)
					fmt.Fprintf(output, "PLACE %v\n", id)
					continue
				}
				if GetVariable(backflip) == UNDEFINED {
				
					if GetVariable("__new__") > 0 || methods[CurrentFunctionName] {
						ExpressionType = LastDefinedType		
						var add = IndexUserType(s, output, LastDefinedType.String(), backflip)
						SetVariable(add, ExpressionType)
						backflip = add
						
						goto anotherlabel
						
					} else {
						RaiseError(s, "`"+backflip+"` is undefined!")
					}
				} else {
					function := functions["text_m_"+GetVariable(backflip).String()]
					if !function.Exists {
						RaiseError(s, "`"+backflip+"` cannot be turned into text! ("+GetVariable(backflip).String()+")")
					}
					
					fmt.Fprintf(output, "%s %s\n", GetVariable(backflip).Push(), backflip)
					if function.Inline {
						fmt.Fprintf(output, "%v\n", function.Data)
					} else {
						fmt.Fprintf(output, "RUN text_m_%v\n", GetVariable(backflip))
					}
					unique++
					fmt.Fprintf(output, "GRAB i+backflip+%v\n", unique)
					
					fmt.Fprintf(output, "PLACE %v\n", id)
					fmt.Fprintf(output, "JOIN %v %v i+backflip+%v\n", id, id, unique)
					fmt.Fprintf(output, "PLACE %v\n", id)
					continue
				}
			}	
			
		} else {
			fmt.Fprintf(output, "PUT %v\n", strconv.Itoa(int(v)))
		}
	}
	
	if len(arg) == 0 {
		goto end
	}
	fmt.Fprintf(output, "PUT %v\n", strconv.Itoa(int(' ')))
	j++
	arg = string(s.TokenText()[j])
	goto stringloop
	end:
	
	ExpressionType = STRING
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

	fmt.Fprintf(output, "ARRAY %v\n", id)

	//Push all the values.
	for tok := s.Scan(); tok != scanner.EOF; {

		if s.TokenText() == "]" {
			break
		}

		fmt.Fprintf(output, "PLACE %v\nPUT %v\n", id, expression(s, output))
	
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
