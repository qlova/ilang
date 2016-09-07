package main

import (
	"text/scanner"
	"fmt"
	"io"
)


//This is an operator such as + - *
//It contains all the information required to compile them.
type Operator struct {

	Assembly string
	Precidence bool
	
	A, B TYPE
	
	ExpressionType TYPE
}

var Operations = make(map[string]map[TYPE]map[TYPE]Operator)

//Opp is a standard arithmetic operator.
func NewOperator(a TYPE, o string, b TYPE, asm string, p bool, args ...TYPE) {
	var typ TYPE = 0
	if len(args) > 0 {
		typ = args[0]
	}
	opp := Operator{
		A:a,
		B:b,
		Assembly:asm,
		Precidence:p,
		ExpressionType:typ,
	}
	if _, ok := Operations[o]; !ok {
		Operations[o] = make(map[TYPE]map[TYPE]Operator)
	}
	if _, ok := Operations[o][a]; !ok {
		Operations[o][a] = make(map[TYPE]Operator)
	}
	Operations[o][a][b] = opp
}

func GetOperator(sym string, a TYPE, b TYPE) (Operator, bool) {
	if _, ok := Operations[sym]; ok {
		if _, ok := Operations[sym][a]; ok {
			if o, ok := Operations[sym][a][b]; ok {
				return o, true
			}
		}
	}
	return Operator{}, false
}


func IsOperator(sym string) bool {
	if _, ok := Operations[sym]; ok {
		return true
	}
	return false
}

func OperatorPrecident(sym string) bool {
	return map[string]bool{"=":true, "!=":true,">=":true,"<=":true,"<":true,">":true,
		"+": true,
		"-": true,
	}[sym]
}

func init() {
	NewOperator(NUMBER, "/", NUMBER, "VAR %c\nDIV %c %a %b", true)
	NewOperator(NUMBER, "÷", NUMBER, "VAR %c\nDIV %c %a %b", true)
	
	NewOperator(NUMBER, "+", NUMBER, "VAR %c\nADD %c %a %b", false)
	NewOperator(NUMBER, "-", NUMBER, "VAR %c\nSUB %c %a %b", false)
	
	NewOperator(NUMBER, "*", NUMBER, "VAR %c\nMUL %c %a %b", true)
	NewOperator(NUMBER, "×", NUMBER, "VAR %c\nMUL %c %a %b", true)
	
	NewOperator(NUMBER, "mod", NUMBER, "VAR %c\nMOD %c %a %b", true)
	NewOperator(NUMBER, "^",   NUMBER, "VAR %c\nPOW %c %a %b", true)
	
	NewOperator(NUMBER, "=", NUMBER, "VAR %c\nSEQ %c %a %b", true)
	NewOperator(NUMBER, "!=",NUMBER, "VAR %c\nSNE %c %a %b", true)
	NewOperator(NUMBER, "<", NUMBER, "VAR %c\nSLT %c %a %b", true)
	NewOperator(NUMBER, ">", NUMBER, "VAR %c\nSGT %c %a %b", true)
	NewOperator(NUMBER, "<=",NUMBER, "VAR %c\nSLE %c %a %b", true)
	NewOperator(NUMBER, ">=",NUMBER, "VAR %c\nSGE %c %a %b", true)
	
	NewOperator(ITYPE, "=", ITYPE, "VAR %c\nSEQ %c %a %b", true)
	
	NewOperator(STRING, "+", STRING, "ARRAY %c\nJOIN %c %a %b", false)
	NewOperator(STRING, "=", STRING, "SHARE %a\n SHARE %b\nRUN strings.equal\nPULL %c\n", false)
	NewOperator(STRING, "!=", STRING, "SHARE %a\n SHARE %b\nRUN strings.equal\nPULL %c\nDIV %c %c 0\n", false)
	
	NewOperator(NUMBER, "²", UNDEFINED, "VAR %c\nPOW %c %a %a", true)
	
	NewOperator(STRING, "#", NUMBER, "SHARE %a\nPUSH %b\nRUN hash\nPULL %c\n", true)
	NewOperator(NUMBER, "?", NUMBER, "PUSH %a\nPUSH %b\nRUN unhash\nGRAB %c\n", true,  STRING)
}

var OperatorFunction bool

func ParseOperator(s *scanner.Scanner, output io.Writer) {
	var A, B TYPE
	var symbol string
	
	A = StringToType[s.TokenText()]
	
	s.Scan()
		symbol = s.TokenText()
		
	s.Scan()
		B = StringToType[s.TokenText()]
		
	s.Scan()
	if s.TokenText() == "{" {
		s.Scan()
	}
	
	NewOperator(A, symbol, B, "SHARE %a\n SHARE %b\nRUN "+fmt.Sprint(A, "_", symbol, "_", B)+"\nGRAB %c\n", OperatorPrecident(symbol))
	
	GainScope()
	fmt.Fprintf(output, "FUNCTION %s_%s_%s\n", A, symbol, B)
	fmt.Fprintf(output, "GRAB b\nGRAB a\nARRAY c\n")
	for range DefinedTypes[A-USER].Elements {
		fmt.Fprintf(output, "PUT 0\n")
	}
	OperatorFunction = true
	
	SetVariable("c", A)
	SetVariable("a", A)
	SetVariable("b", B)
}
