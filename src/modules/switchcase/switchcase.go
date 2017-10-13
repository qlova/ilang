package switchcase

import "github.com/qlova/ilang/src"

func init() {
	ilang.RegisterToken([]string{"switch"}, ScanSwitch)
	ilang.RegisterToken([]string{"case"}, ScanCase)
	ilang.RegisterToken([]string{"default"}, ScanDefault)
	
	ilang.RegisterListener(Switch, SwitchEnd)
}

var Switch = ilang.NewFlag()
var Default = ilang.NewFlag()

func SwitchEnd(ic *ilang.Compiler) {

	//We are in a typeloop.
	if ic.GetVariable("flag_type")  != ilang.Undefined{
		return
	}

	for i:=0; i < (ic.GetVariable("flag_nesting").Int); i++ {
		ic.Assembly("END")
	}
	ic.Assembly("END")
	ic.LoseScope()
	ic.Assembly("#ACTIVATE")	
}

func ScanSwitch(ic *ilang.Compiler) {

	//We are in a typeloop.
	if ic.GetVariable("flag_type") != ilang.Undefined {
		var check = ic.Scan(0)
		if check == "type" {	
			SwitchType(ic)
		} else {
			ic.NextToken = check
		}
		return
	}

	var expression = ic.ScanExpression()
	if ic.ExpressionType != ilang.Number {
		ic.RaiseError("switch statements must have numeric conditions!")
	}
	ic.Scan('{')
	ic.Assembly("#SWITCH")
	ic.GainScope()
	
	ic.SetFlag(ilang.Type{Name: "flag_switch", Push: expression})
	ic.SetFlag(ilang.Type{Name: "flag_nesting", Int: 0})
	
	//Find first case.
	for {
		token := ic.Scan(0)
		if token != "\n" {
			if token != "case" {
				ic.RaiseError("Expecting 'case', found ", token)
			}
			break
		}
	}
	expression = ic.ScanExpression()
	var condition = ic.Tmp("case")
	ic.Assembly("VAR ", condition)
	ic.Assembly("SEQ %v %v %v", condition, expression, ic.GetVariable("flag_switch").Push)
	ic.Assembly("IF ",condition)
	ic.GainScope()
	ic.SetFlag(Switch)
}

func ScanDefault(ic *ilang.Compiler) {
	if !ic.GetFlag(Switch) {
		ic.RaiseError("'default' must be within a 'switch' block!")
	}
	ic.UnsetFlag(Switch)
	ic.LoseScope()
	ic.Assembly("ELSE")
	ic.GainScope()
	ic.SetFlag(Switch)
	ic.SetFlag(Default)
}

func ScanCase(ic *ilang.Compiler) {
	if !ic.GetScopedFlag(Switch) {
		ic.RaiseError("a 'case' must be within a 'switch' block!")
	}
	
	if ic.GetScopedFlag(Default) {
		ic.RaiseError("default must be at the end of the switch statement!")
	}
	
	//We are in a typeloop.
	if ic.GetVariable("flag_type")  != ilang.Undefined {
		CaseType(ic)
		return
	}

	var expression = ic.ScanExpression()
	var condition = ic.Tmp("case")
	
	ic.UnsetFlag(Switch)
	
	ic.LoseScope()
	
	ic.UpdateVariable("flag_nesting", ilang.Type{Int: ic.GetVariable("flag_nesting").Int+1 })
	
	ic.Assembly("ELSE")
	
	ic.Assembly("VAR ", condition)
	ic.Assembly("SEQ %v %v %v", condition, expression, ic.GetVariable("flag_switch").Push)
	ic.Assembly("IF ",condition)
	ic.GainScope()
	ic.SetFlag(Switch)
}
