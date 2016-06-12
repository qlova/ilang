package main

import (
	"fmt"
	"text/scanner"
	"io"
	"os"
	"strconv"
)

//This is an expression shunter. It takes the current identifyer and shunts it into the next operator.
//I don't think this has anything to do with the shunting yard algorithim, I just like the term.
func shunt(name string, s *scanner.Scanner, output io.Writer) string {

		//Scan the next token.
		s.Scan()
		
		//If it is one of these characters, then we have finished our shunt.
		switch s.TokenText() {
			case ")", ",", "\n", "]", ";":
				return name
		}
		
		//I love doing the shunting.
		if operator, ok := Operators[s.TokenText()]; ok {
		
			//Multi-opps.
			if _, ok := Operators[string(s.Peek())]; ok {
				operator = Operators[s.TokenText()+string(s.Peek())]
				s.Scan()
				if operator.code == "" {
					fmt.Println("[SHUNTING YARD] Invalid operator matchup: "+s.TokenText()+string(s.Peek()))
					os.Exit(1)
				}
			}
		
			//Here we create the unique name for the shunting result.
			unique++
			id := "i+shunt+"+fmt.Sprint(unique)
			
			//Operators have some defined formats (can be found in operators.go)
			s.Scan()
			switch operator.mode {
				case 0:
					fmt.Fprintf(output, operator.code, id, id, name, expression(s, output, operator.shunt))
				case 1:
					fmt.Fprintf(output, operator.code, id, id, name, name)
				case 2:
					fmt.Fprintf(output, operator.code, name, expression(s, output, operator.shunt), id)
			}
			
			//This is for operator precidence.
			if !operator.shunt {
				return shunt(id, s, output)
			}
			return id
		}
		
		//Special case for indexing arrays.
		if _, err := strconv.Atoi(s.TokenText()[1:]); err == nil && s.TokenText()[0] == '.' {
			unique++
			output.Write([]byte("INDEX "+name+" "+s.TokenText()[1:]+" i+shunt+"+fmt.Sprint(unique)+"\n"))
			s.Scan()
			return "i+shunt+"+fmt.Sprint(unique)
		} else if s.TokenText() == "." {
			s.Scan()
			unique++
			output.Write([]byte("INDEX "+name+" "+s.TokenText()+" i+shunt+"+fmt.Sprint(unique)+"\n"))
			s.Scan()
			return "i+shunt+"+fmt.Sprint(unique)
		}
		
		fmt.Println("[SHUNTING YARD] Unexpected ", s.TokenText(), "("+name+")")
		os.Exit(1)
		return ""
}
