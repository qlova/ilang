package ilang

//Pretty much the compiler API.

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

func RegisterToken(tokens []string, f func(*Compiler)) {
	for _, name := range tokens {
		Tokens[name] = f
	}
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


func RegisterSymbol(name string, f func(*Compiler) Type) {
	Symbols[name] = f
}

func RegisterVariable(f func(*Compiler, string) Type) {
	Variables = append(Variables, f)
}

func RegisterConstructor(f func(*Compiler, Type)) {
	Constructors = append(Constructors, f)
}

func RegisterExpression(f func(*Compiler) string) {
	Expressions = append(Expressions, f)
}

func RegisterShunt(token string, f func(*Compiler, string) string) {
	Shunts[token] = append(Shunts[token], f)
}
