package main

func (ic *Compiler) ScanSwitch() {
	var expression = ic.ScanExpression()
	if ic.ExpressionType != Number {
		ic.RaiseError("switch statements must have numeric conditions!")
	}
	ic.Scan('{')
	ic.GainScope()
	ic.SetFlag(Type{Name: "flag_switch", Push: expression})
	
	for {
		token := ic.Scan(0)
		if token != "\n" {
			if token != "case" {
				ic.RaiseError("Expecting case")
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
}

func (ic *Compiler) ScanDefault() {
	if ic.GetVariable("flag_switch") == Undefined {
		ic.RaiseError("'default' must be within a 'switch' block!")
	}
	ic.LoseScope()
	ic.Assembly("ELSE")
	ic.GainScope()
}

func (ic *Compiler) ScanCase() {
	if ic.GetVariable("flag_switch") == Undefined {
		ic.RaiseError("'case' must be within a 'switch' block!")
	}

	var expression = ic.ScanExpression()
	var condition = ic.Tmp("case")
	nesting, ok := ic.Scope[len(ic.Scope)-2]["flag_nesting"]
	if !ok {
		nesting.Int = 0
	}
	
	ic.LoseScope()
	
	ic.Assembly("ELSE")
	ic.SetVariable("flag_nesting", Type{Int:nesting.Int+1})
	
	ic.Assembly("VAR ", condition)
	ic.Assembly("SEQ %v %v %v", condition, expression, ic.GetVariable("flag_switch").Push)
	ic.Assembly("IF ",condition)
	ic.GainScope()
}
