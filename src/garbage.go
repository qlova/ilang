package ilang

import "strings"

func (ic *Compiler) CollectGarbage() {
	//Erm garbage collection???
	for name, variable := range ic.Scope[len(ic.Scope)-1] {
		if strings.Contains(name, "_") {
			var ok = false
			if ic.LastDefinedType.Detail != nil {
				_, ok = ic.LastDefinedType.Detail.Table[strings.Split(name, "_")[0]]
			}
			if variable == Unused && !(ic.GetFlag(InMethod) && ok ) {
				ic.RaiseError("unused variable! ", strings.Split(name, "_")[0])
			} 
		}
	
		if ic.Scope[len(ic.Scope)-1][name+"."] != Protected { //Protected variables
			if ic.GetFlag(InMethod) && name == ic.LastDefinedType.Name {
				continue
			}
			
			//Possible memory leak, TODO check up on this.
			if _, ok := ic.DefinedTypes[name]; ic.GetFlag(InMethod) && ok {
				continue
			}
			
			if variable.IsUser() != Undefined && !variable.Empty() {
				ic.Assembly("SHARE ", name)
				ic.Assembly(ic.RunFunction("collect_m_"+variable.Name))
			}
		}
	}
}

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
				ic.Library(ic.RunFunction("collect_m_"+element.Name))
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
