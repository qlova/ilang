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
	
	//Didn't my master always say 
	//		"NO GOTO STATEMENTS, YOUR MAKING ME BLIND!"
	stringloop:
	arg = strings.Replace(arg, "\\n", "\n", -1)
	arg = strings.Replace(arg, "\\r", "\r", -1)
	arg = strings.Replace(arg, "\\\n", "\\n", -1)
	for _, v := range arg {
		if v == '"' {
			goto end
		}
		fmt.Fprintf(output, "PUT %v\n", strconv.Itoa(int(v)))
	}
	if len(arg) == 0 {
		goto end
	}
	fmt.Fprintf(output, "PUT %v\n", strconv.Itoa(int(' ')))
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
