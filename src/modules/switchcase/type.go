package switchcase

import "github.com/qlova/ilang/src"

func SwitchType(ic *ilang.Compiler) {
	ic.Scan('{')
	ic.Assembly("#SWITCH_TYPE")
	ic.GainScope()
	ic.SetFlag(Switch)
	
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
	CaseType(ic)
}

func CaseType(ic *ilang.Compiler) {
	var token = ic.Scan(0)

	if t := ilang.GetType(token); (t == ilang.Type{}) {
		ic.RaiseError(token, " is not a defined type!")
	} else if ic.GetVariable("element_type") == t {
		return
	} else {
		for {
			token := ic.Scan(0)
			if token == "case" || token == "default" || token == "}" {
				ic.NextToken = token
				return
			}
		}	
	}
}
