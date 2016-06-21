package main

//This is an operator such as + - *
//It contains all the information required to compile them.
type Operator struct {
	code string
	shunt bool
	mode int
}

//Opp is a standard arithmetic operator.
func opp(a string, shunt ...bool) Operator {
	return Operator{code:a, shunt:len(shunt)<=0}
}

//Pow is an operator such as x²
func pow(a string, shunt ...bool) Operator {
	return Operator{code:a, shunt:len(shunt)<=0, mode:1}
}

//App is an operator which calls a function.
func app(a string, shunt ...bool) Operator {
	return Operator{code:a, shunt:len(shunt)<=0, mode:2}
}

//Operator data.
//Contains format for compliation.
var Operators = map[string]Operator{
	"/": 	opp( "VAR %v\nDIV %v %v %v\n", true),
	"÷": 	opp( "VAR %v\nDIV %v %v %v\n", true),
	"+": 	opp( "VAR %v\nADD %v %v %v\n"),
	"-": 	opp( "VAR %v\nSUB %v %v %v\n"),
	
	//Should these be kept??
	"and":	opp( "VAR %v\nMUL %v %v %v\n"),
	"or":	opp( "VAR %v\nADD %v %v %v\n"),
	
	"*":	opp( "VAR %v\nMUL %v %v %v\n", true),
	"×":	opp( "VAR %v\nMUL %v %v %v\n", true),
	"mod": 	opp( "VAR %v\nMOD %v %v %v\n", true),
	"^": 	opp( "VAR %v\nPOW %v %v %v\n", true),
	"&":	opp( "STRING %v\nJOIN %v %v %v\n"),
	"=":	opp( "VAR %v\nSEQ %v %v %v\n"),
	"!=":	opp( "VAR %v\nSNE %v %v %v\n"),
	"<=":	opp( "VAR %v\nSLE %v %v %v\n"),
	"<":	opp( "VAR %v\nSLT %v %v %v\n"),
	">":	opp( "VAR %v\nSGT %v %v %v\n"),
	">=":	opp( "VAR %v\nSGE %v %v %v\n"),
	
	"²":	pow( "VAR %v\nMUL %v %v %v\n", true),
	
	"==":	app( "PUSHSTRING %v\n PUSHSTRING %v\nRUN strings.equal\nPOP %v\n"),
	"@":	app( "PUSHSTRING %v\nPUSH %v\nRUN hash\nPOP %v\n", true),
	"?":	app( "PUSH %v\nPUSH %v\nRUN unhash\nPOPSTRING %v\n", true),
	
	"!": Operator{},
}







