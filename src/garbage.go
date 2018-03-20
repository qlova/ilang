package ilang

import "strings"

func (ic *Compiler) CollectGarbage() {
	//Erm garbage collection???
	for name, variable := range ic.Scope[len(ic.Scope)-1] {
		if strings.Contains(name, "_used") {
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
			
			if (variable.IsUser() != Undefined)  || (variable.SubType != nil && (variable.SubType.Push != "PUSH" || variable.SubType.SubType != nil)) {
				
				if !variable.Empty() {
				ic.Collect(variable)
				ic.Assembly(variable.Push, " ", name)
				ic.Assembly(ic.RunFunction("collect_m_"+variable.GetComplexName()))
				}
			}
		}
	}
}

var AlreadyGeneratedACollectionMethodFor = make(map[Type]bool)

func (t Type) Free(pointer string) string {
	if (t.IsUser() == Undefined && t.SubType == nil) || (t.SubType.Push == "PUSH" && t.SubType.SubType == nil) {
		return ""
	}
	if t.Empty() {
		return ""
	}
	
	return "IF "+pointer+"\nPUSH "+pointer+
	"\nHEAP\nRUN collect_m_"+t.GetComplexName()+
	"\nMUL "+pointer+" -1 "+pointer+"\nPUSH "+pointer+"\nHEAP\nEND"
}

func (ic *Compiler) Collect(t Type) {
	if AlreadyGeneratedACollectionMethodFor[t] {
		return
	}
	
	if t.Empty() {
		return
	}
	
	if t.IsUser() == Undefined  {
		//TODO collect lists, tables and other complex types!
		
		for _, collection := range Collections {
			collection(ic, t)
		}
		AlreadyGeneratedACollectionMethodFor[t] = true
		
		return
	}

	
	var scope = ic.Scope
	ic.GainScope()
	ic.Library("FUNCTION collect_m_", t.GetComplexName())
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
			ic.Library("PUSH ", tmp)
			ic.Library("HEAP")
			ic.Library("MUL %v %v -1", tmp, tmp)
			if element.IsUser() != Undefined {
				ic.Library("PUSH ", tmp)
				ic.Library("HEAP")
				
				var member = ic.Tmp("member")
				ic.Library("GRAB ", member)
				ic.Library("IF #", member)
					ic.Library("SHARE ", member)
					ic.Library(ic.RunFunction("collect_m_"+element.Name))
				ic.Library("END")
			}
			ic.LoseScope()
			ic.Library("END")
		}
	}
	ic.LoseScope()
	ic.Library("RETURN")
	ic.Scope = scope
	
	AlreadyGeneratedACollectionMethodFor[t] = true
}
