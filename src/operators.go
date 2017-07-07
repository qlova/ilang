package ilang


//This is an operator such as + - *
//It contains all the information required to compile them.
type Operator struct {

	Assembly string
	Precidence bool
	
	A, B Type
	
	ExpressionType Type
}

var Operations = make(map[string]map[Type]map[Type]Operator)

//Opp is a standard arithmetic operator.
func NewOperator(a Type, o string, b Type, asm string, p bool, args ...Type) {
	var typ Type = a
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
		Operations[o] = make(map[Type]map[Type]Operator)
	}
	if _, ok := Operations[o][a]; !ok {
		Operations[o][a] = make(map[Type]Operator)
	}
	Operations[o][a][b] = opp
}

func GetOperator(sym string, a Type, b Type) (Operator, bool) {
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
		"+=": true,
		"-=": true,
		"*=": true,
		"/=": true,
	}[sym]
}

func init() {
	NewOperator(Number, "/", Number, "VAR %c\nDIV %c %a %b", true)
	NewOperator(Number, "÷", Number, "VAR %c\nDIV %c %a %b", true)
	
	NewOperator(Number, "+", Number, "VAR %c\nADD %c %a %b", false)
	NewOperator(Number, "-", Number, "VAR %c\nSUB %c %a %b", false)
	
	NewOperator(Decimal, "+", Decimal, "VAR %c\nADD %c %a %b", false)
	NewOperator(Decimal, "-", Decimal, "VAR %c\nSUB %c %a %b", false)
	NewOperator(Decimal, "/", Decimal, "VAR %t\nVAR %c\nMUL %t %a 1000000\nDIV %c %t %b", true)
	NewOperator(Decimal, "*", Decimal, "VAR %c\nMUL %c %a %b\nDIV %c %c 1000000", true)
	
	NewOperator(Number, "++", Undefined, "ADD %a %a 1", false, Undefined)
	NewOperator(Number, "--", Undefined, "SUB %a %a 1", false, Undefined)
	
	NewOperator(Number, "+=", Number, "ADD %a %a %b", false, Undefined)
	NewOperator(Number, "-=", Number, "SUB %a %a %b", false, Undefined)
	NewOperator(Number, "*=", Number, "MUL %a %a %b", false, Undefined)
	NewOperator(Number, "/=", Number, "DIV %a %a %b", false, Undefined)
	
	NewOperator(Decimal, "+=", Decimal, "ADD %a %a %b", false, Undefined)
	NewOperator(Decimal, "-=", Decimal, "SUB %a %a %b", false, Undefined)
	NewOperator(Decimal, "*=", Decimal, "MUL %c %a %b\nDIV %c %c 1000000", false, Undefined)
	NewOperator(Decimal, "/=", Decimal, "VAR %t\nMUL %t %a 1000000\nDIV %c %t %b", false, Undefined)
	
	//Set Operations.
	NewOperator(Set, "=", Set, "VAR %c\nSEQ %c %a %b", true, Number)
	NewOperator(Set, "+", Set, "VAR %c\nMUL %c %a %b", false)
	NewOperator(Set, "-", Set, "VAR %c\nDIV %c %a %b", false)
	NewOperator(Set, "+=", Set, "MUL %a %a %b", false, Undefined)
	NewOperator(Set, "-=", Set, "DIV %a %a %b", false, Undefined)
	NewOperator(Set, "<=", Set, "VAR %c\nMOD %c %b %a\nDIV %c %c 0", true, Number)
	NewOperator(Set, ">", Set, "VAR %c\nMOD %c %b %a\n", true, Number)
	NewOperator(Set, ">=", Set, "VAR %c\nMOD %c %b %a\nVAR %t\nSEQ %t %a %b\nADD%c %c %t", true, Number)
	NewOperator(Set, "<", Set, "VAR %c\nMOD %c %b %a\nDIV %c %c 0\nVAR %t\nSNE %t %a %b\nMUL%c %c %t", true, Number)
	
	NewOperator(Number, "or", Number, "VAR %c\nADD %c %a %b", false)
	
	NewOperator(Number, "and", Number, "VAR %c\nMUL %c %a %b", true)
	
	NewOperator(Number, "*", Number, "VAR %c\nMUL %c %a %b", true)
	NewOperator(Number, "×", Number, "VAR %c\nMUL %c %a %b", true)
	
	NewOperator(Number, "mod", Number, "VAR %c\nMOD %c %a %b", true)
	NewOperator(Decimal, "mod", Decimal, "VAR %c\nMOD %c %a %b", true)
	NewOperator(Number, "^",   Number, "VAR %c\nPOW %c %a %b", true)
	
	NewOperator(Decimal, "^", Number, "VAR %c\nPOW %c %a %b\nVAR %t\nPOW %t 100 %b\nMUL %t %t 1000000\nDIV %c %c %t", true, Decimal)
	
	NewOperator(Letter, "=", Letter, "VAR %c\nSEQ %c %a %b", true, Number)
	NewOperator(Letter, "!=",Letter, "VAR %c\nSNE %c %a %b", true, Number)
	
	NewOperator(Number, "=", Number, "VAR %c\nSEQ %c %a %b", true)
	NewOperator(Number, "!=",Number, "VAR %c\nSNE %c %a %b", true)
	NewOperator(Number, "<", Number, "VAR %c\nSLT %c %a %b", true)
	NewOperator(Number, ">", Number, "VAR %c\nSGT %c %a %b", true)
	NewOperator(Number, "<=",Number, "VAR %c\nSLE %c %a %b", true)
	NewOperator(Number, ">=",Number, "VAR %c\nSGE %c %a %b", true)
	
	NewOperator(Decimal, "=", Decimal, "VAR %c\nSEQ %c %a %b", true, Number)
	NewOperator(Decimal, "!=",Decimal, "VAR %c\nSNE %c %a %b", true, Number)
	NewOperator(Decimal, "<", Decimal, "VAR %c\nSLT %c %a %b", true, Number)
	NewOperator(Decimal, ">", Decimal, "VAR %c\nSGT %c %a %b", true, Number)
	NewOperator(Decimal, "<=",Decimal, "VAR %c\nSLE %c %a %b", true, Number)
	NewOperator(Decimal, ">=",Decimal, "VAR %c\nSGE %c %a %b", true, Number)
	
	NewOperator(Itype, "=", Itype, "VAR %c\nSEQ %c %a %b", true)
	
	NewOperator(Text, "+", Text, "ARRAY %c\nJOIN %c %a %b", false)
	NewOperator(Array, "+", Array, "ARRAY %c\nJOIN %c %a %b", false)
	
	NewOperator(Text, "--", Undefined, "PLACE %a\nPOP %t", false, Undefined)
	
	NewOperator(Text, "+=", Letter, "PLACE %a\nPUT %b", false, Undefined)
	NewOperator(Text, "+=", Text, "JOIN %a %a %b", false, Undefined)
	NewOperator(Array, "+=", Array, "JOIN %a %a %b", false, Undefined)
	NewOperator(Array, "+=", Number, "PLACE %a\nPUT %b", false, Undefined)
	
	NewOperator(Text, "=", Text, "SHARE %a\n SHARE %b\nRUN strings.equal\nPULL %c\n", false, Number)
	NewOperator(Text, "!=", Text, "SHARE %a\n SHARE %b\nRUN strings.equal\nPULL %c\nDIV %c %c 0\n", false, Number)
	
	NewOperator(Number, "²", Undefined, "VAR %c\nPOW %c %a %a", true)
	
	NewOperator(Text, "#", Number, "SHARE %a\nPUSH %b\nRUN hash\nPULL %c\n", true, Number)
	NewOperator(Number, "?", Number, "PUSH %a\nPUSH %b\nRUN unhash\nGRAB %c\n", true,  Text)
	
	//For grate engine.
	NewOperator(Number, "px", Undefined, "VAR %c\nMUL %c %a 10", true)
}
