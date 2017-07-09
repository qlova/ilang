package ilang

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

func (ic *Compiler) ScanListStatement() {
	var name = ic.Scan(0)
	var t  = ic.GetVariable(name)
	var token = ic.Scan(0)
	
	//TODO CLEAN THIS UP!
	switch token {
		case "-":
			ic.Scan('-')
			ic.Assembly("PLACE ", name)
			ic.Assembly("POP ", ic.Tmp("cut"))
			
		case "&", "+":
		
			if token == "+" {
				ic.Scan('=')
			}
			value := ic.ScanExpression()
			
			if t == List && ic.ExpressionType.Push == "PUSH" {
				ic.UpdateVariable(name, Array)
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUT ", value)
				return
			}
			
			if t == List {
				list := ic.ExpressionType
				list.List = true
				list.User = false
				t = list
				ic.UpdateVariable(name, list)
				//println(name)
				if ic.GetFlag(InMethod) {
					ic.LastDefinedType.Detail.Elements[ic.LastDefinedType.Detail.Table[name]] = t
				}
			}
			
			//This appends elements to a list {..}
			if t.List {
				ic.PutList(t, name, value)
				
			} else {
			
				if ic.ExpressionType.Push != "PUSH" {
					ic.RaiseError("Only numeric values can be added to arrays.")
				}
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUT ", value)
			}
		case "[":
			var index = ic.ScanExpression()
			ic.Scan(']')
			ic.Scan('=')
			var value = ic.ScanExpression()
			
			ic.Set(name, index, value)
		case "=":
			value := ic.ScanExpression()
			if ic.ExpressionType != t {
				ic.RaiseError("Only ",t.Name," values can be assigned to ",name,".")
			}
			
			if _, ok := ic.LastDefinedType.Detail.Table[name]; ic.GetFlag(InMethod) && ok {
				ic.SetUserType(ic.LastDefinedType.Name, name, value)	
			} else {									
				ic.Assembly("PLACE ", value)
				ic.Assembly("RENAME ", name)
			}
		case "has":
			if ic.GetFlag(InMethod) {
				ic.SetUserType(ic.LastDefinedType.Name, name, ic.ScanList())
			} else {
				ic.AssembleVar(name, ic.ScanList())
			}
			
		default:
			ic.ExpressionType = t
			ic.NextToken = token
			ic.Shunt(name)
			if ic.ExpressionType != Undefined {
				ic.RaiseError("blank expression!")
			}
	}
}

func (ic *Compiler) NewListOf(t Type) string {
	if t.Empty() {
		ic.RaiseError("Cannot create an array of ",t.Name,"! (The type has no size)")
	}

	t.List = true
	t.User = false
	ic.ExpressionType = t
	var list = ic.Tmp("list")
	ic.Assembly("ARRAY ", list)
	return list
}

func (ic *Compiler) ScanList() string {
	var name = ic.Scan(Name)
	
	t, ok := ic.DefinedTypes[name]
	if !ok {
		if i, ok := ic.DefinedInterfaces[name]; !ok {
			ic.RaiseError(name+" is an unrecognised type!")
		} else {
			t = i.GetType()
		}
	}
	t.List = true
	t.User = false
	
	var list = ic.Tmp("list")
	
	ic.Scan('(')
	if tok := ic.Scan(0); tok != "s" {
		ic.NextToken = tok
		size := ic.ScanExpression()
		if ic.ExpressionType != Number {
			ic.RaiseError("Expecting list size!")
		}
		ic.Assembly("PUSH ", size)
		ic.Assembly("MAKE")
		ic.Assembly("GRAB ", list)
	} else {
		ic.Assembly("ARRAY ", list)
	}
	ic.Scan(')')
	
	ic.ExpressionType = t
	
	
	
	return list
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
