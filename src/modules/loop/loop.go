package loop

import "github.com/qlova/ilang/src"

var Loop = ilang.NewFlag()

func init() {
	ilang.RegisterToken([]string{"loop"}, ScanLoop)
	ilang.RegisterListener(Loop, EndLoop)
	
	ilang.RegisterToken([]string{"break"}, func(ic *ilang.Compiler) {
		ic.CollectGarbage()
		ic.Assembly("BREAK")
	})
}

func ScanLoop(ic *ilang.Compiler) {
	ic.Assembly("LOOP")
	ic.GainScope()
	ic.NextToken = ic.Scan(0)
	if ic.NextToken != "{" {
		condition := ic.ScanExpression()
		ic.Assembly("SEQ ", condition, " 0 ", condition)
		ic.Assembly("IF ", condition)
		ic.Assembly("BREAK")
		ic.Assembly("END")
	}
	ic.Scan('{')
	ic.SetFlag(Loop)
}

func EndLoop(ic *ilang.Compiler) {
	if ic.LastToken != "}" {
		ic.RaiseError("Loop must be ended with a '}' token, found '", ic.LastToken, "'")
	}
	ic.Assembly("REPEAT")
}
