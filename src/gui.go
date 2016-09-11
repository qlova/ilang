package main

import (
	"text/scanner"
	"fmt"
	"io"
	"strings"
)

var GUIEnabled bool

// gui main { <button> </button> }
func ParseGUIDef(s *scanner.Scanner, output io.Writer) {
	s.Scan()
	
	fmt.Fprintf(output, "DATA gui_%s \"", s.TokenText())
	SetVariable("gui_"+s.TokenText(), STRING)
	s.Scan()
	
	Expecting(s, "{")
	
	var design string
	var jsbraces = 0
	s.Scan()
	for {
		s.Scan()
		if s.TokenText() == "}"  {
	 		if jsbraces == 0 {
				break
			} else {
				jsbraces--
			}
		}
		if s.TokenText() == "{" {
			jsbraces++
		}
		if s.TokenText() != "\n" {
			if strings.Contains(s.TokenText(), "\"") {
				design += strings.Replace(s.TokenText(), "\"", "\\\"", -1)
			} else {
				design += s.TokenText()
			}
			if string(s.Peek()) == " " {
				design += " "		
			}
		}
	}
	
	fmt.Fprintf(output, "%s\"\n", design)
	GUIEnabled = true
}
