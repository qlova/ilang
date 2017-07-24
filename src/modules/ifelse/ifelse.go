package ifelse
import "github.com/qlova/ilang/src"

func init() {
	ilang.RegisterToken([]string{
		"if",
	}, ScanIf)
	
	ilang.RegisterToken([]string{
		"else",
	}, ScanElse)
	
	ilang.RegisterToken([]string{
		"elseif",
	}, ScanElseIf)
	
	ilang.RegisterListener(If, IfEnd)
	ilang.RegisterListener(Else, ElseIfEnd)
	ilang.RegisterListener(ElseIf, ElseIfEnd)
}

var If = ilang.NewFlag()
var Else = ilang.NewFlag()
var ElseIf = ilang.NewFlag()

func ElseIfEnd(ic *ilang.Compiler) {
	nesting, ok := ic.Scope[len(ic.Scope)-1]["flag_nesting"]
	if ok {
		for i:=0; i < nesting.Int; i++ {
			ic.Assembly("END")
		}
	}
	ic.LoseScope()
	ic.Assembly("END")
}

func IfEnd(ic *ilang.Compiler) {
	ic.Assembly("END")
}

func ScanIf(ic *ilang.Compiler) {
	var expression = ic.ScanExpression()
	if ic.ExpressionType != ilang.Number {
		ic.RaiseError("if statements must have numeric conditions!")
	}
	ic.Assembly("IF ", expression)
	ic.GainScope()
	ic.SetFlag(If)
}

func ScanElse(ic *ilang.Compiler) {
	nesting, ok := ic.Scope[len(ic.Scope)-1]["flag_nesting"]
	if !ok {
		nesting.Int = 0
	}
	
	if ic.GetScopedFlag(If) {
		ic.UnsetFlag(If)
	} else if ic.GetScopedFlag(ElseIf) {
		ic.UnsetFlag(ElseIf)
		ic.LoseScope()
		
		nesting, ok = ic.Scope[len(ic.Scope)-1]["flag_nesting"]
		if !ok {
			nesting.Int = 0
		}
		
	} else {
		ic.RaiseError("You cannot have an else without an if!")
	}
	
	ic.LoseScope()
	ic.Assembly("ELSE")
	ic.GainScope()
	ic.SetVariable("flag_nesting", ilang.Type{Int:nesting.Int})
	ic.GainScope()
	ic.SetFlag(Else)
}

func ScanElseIf(ic *ilang.Compiler) {
	nesting, ok := ic.Scope[len(ic.Scope)-1]["flag_nesting"]
	if !ok {
		nesting.Int = 0
	}
	
	if ic.GetScopedFlag(If) {
		ic.UnsetFlag(If)
	} else if ic.GetScopedFlag(ElseIf) {
		ic.UnsetFlag(ElseIf)
		ic.LoseScope()
		
		nesting, ok = ic.Scope[len(ic.Scope)-1]["flag_nesting"]
		if !ok {
			nesting.Int = 0
		}
		
	} else if ic.GetScopedFlag(Else) {
		ic.RaiseError("Cannot use ifelse after an else...")
	} else {
		ic.RaiseError("Cannot use ifelse without an if!")
	}
	
	ic.LoseScope()
	ic.Assembly("ELSE")
	var expression = ic.ScanExpression()
	ic.Assembly("IF ", expression)
	ic.GainScope()
	ic.SetVariable("flag_nesting", ilang.Type{Int:nesting.Int+1})
	ic.GainScope()
	ic.SetFlag(ElseIf)
}
