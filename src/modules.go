package ilang

//Pretty much the compiler API.

var EnglishTokens = make(map[string]func(*Compiler))
var Tokens = make(map[string]func(*Compiler))
var Listeners = make(map[Type]func(*Compiler))
var Functions = make(map[string]Function)
var Expressions = make([]func(*Compiler) string, 0, 16)
var Statements = make(map[Type]func(*Compiler))
var Symbols = make(map[string]func(*Compiler) Type)
var Defaults = make([]func(*Compiler) bool, 0, 4)
var Variables = make([]func(*Compiler, string) Type, 0, 4)
var Shunts = make(map[string][]func(*Compiler, string) string)
var Constructors = make([]func(*Compiler, Type), 0, 4)

var Casts = make([]func(*Compiler, string, Type, Type) string, 0, 4)

var FunctionBuilders =  make([]func(string) *Function, 0, 4)

var Collections = make([]func(*Compiler, Type), 0, 4)


var SpecialOperators = make([]func(string, Type, Type) *Operator, 0, 4)

//Register a statement token.
func RegisterToken(tokens []string, f func(*Compiler)) {
	for _, name := range tokens {
		Tokens[name] = f
	}
	EnglishTokens[tokens[0]] = f
}

func RegisterListener(listener Type, f func(*Compiler)) {
	Listeners[listener] = f
}

func RegisterStatement(listener Type, f func(*Compiler)) {
	Statements[listener] = f
}

func RegisterDefault(f func(*Compiler) bool) {
	Defaults = append(Defaults, f)
}

func RegisterFunction(name string, f Function) {
	Functions[name] = f
}

func RegisterFunctionBuilder(f func(string) *Function) {
	FunctionBuilders = append(FunctionBuilders, f)
}

func RegisterCast(f func(*Compiler, string, Type, Type) string) {
	Casts = append(Casts, f)
}

func RegisterSymbol(name string, f func(*Compiler) Type) {
	Symbols[name] = f
}

func RegisterVariable(f func(*Compiler, string) Type) {
	Variables = append(Variables, f)
}

func RegisterConstructor(f func(*Compiler, Type)) {
	Constructors = append(Constructors, f)
}

func RegisterSpecialOperator(f func(string, Type, Type) *Operator) {
	SpecialOperators = append(SpecialOperators, f)
}

//Register a new type of expression, return an empty string if you don't recognise anything.
func RegisterExpression(f func(*Compiler) string) {
	Expressions = append(Expressions, f)
}

func RegisterShunt(token string, f func(*Compiler, string) string) {
	Shunts[token] = append(Shunts[token], f)
}

func RegisterCollection(f func(*Compiler, Type)) {
	Collections = append(Collections, f)
}
