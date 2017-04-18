package main

func (t Type) MakeList() Type {
	t.List = true
	t.User = false
	return t
}

func (t Type) ListType() Type {
	t.List = false
	if t.Int > Something.Int {
		t.User = true
	}
	return t
}

//Add a value to a list of type t.
func (ic *Compiler) PutList(t Type, list string, value string) {
		if ic.ExpressionType.Name != t.Name {
		if t.Name == "Something" {
			var tmp = ic.Tmp("something")
			ic.Assembly("ARRAY ", tmp)
			ic.SetVariable(tmp, t)
			ic.AssignSomething(tmp, value)
			ic.SetVariable(tmp, Undefined)
			value = tmp	
		} else {
		ic.RaiseError("Type mismatch! Cannot add a ", ic.ExpressionType.Name,
			 " to a List of ", t.Name)
		}
	}

	var tmp = ic.Tmp("index")
	ic.Assembly("SHARE ", value)
	ic.Assembly("PUSH 0")
	ic.Assembly("HEAP")
	ic.Assembly("PULL ", tmp)

	ic.Assembly("PLACE ", list)
	ic.Assembly("PUT ", tmp)
}
