package ilang

func (ic *Compiler) Collect(t Type) {
	if t.IsUser() == Undefined || t.Empty() {
		return
	}
	
	var scope = ic.Scope
	ic.GainScope()
	ic.Library("FUNCTION collect_m_", t.Name)
	ic.GainScope()
	ic.Library("GRAB variable")
	ic.Library("PLACE variable")
	
	for i, element := range t.Detail.Elements {
		if element == Text || element.IsUser() != Undefined {
			var tmp = ic.Tmp("gc")
			ic.Library("PUSH ", i)
			ic.Library("GET ", tmp)
			ic.Library("IF ", tmp)
			ic.GainScope()
			if element.IsUser() != Undefined {
				ic.Library("PUSH ", tmp)
				ic.Library("HEAP")
				ic.Library("RUN collect_m_", element.Name)
			}
			ic.Library("MUL %v %v -1", tmp, tmp)
			ic.Library("PUSH ", tmp)
			ic.Library("HEAP")
			ic.LoseScope()
			ic.Library("END")
		}
	}
	ic.LoseScope()
	ic.Library("RETURN")
	ic.Scope = scope
}
