package main

func (ic *Compiler) AssembleVar(name string, value string) {

	ic.Assembly("%v %v", ic.ExpressionType.Push, value)
	ic.Assembly("%v %v", ic.ExpressionType.Pop, name)
	ic.SetVariable(name, ic.ExpressionType)
	
	list := ic.ExpressionType
	list.List = true
	list.User = false
	if ic.GetFlag(InMethod) && ic.ExpressionType != List &&
		ic.LastDefinedType.Detail.Elements[ic.LastDefinedType.Detail.Table[name]] == List {

		ic.LastDefinedType.Detail.Elements[ic.LastDefinedType.Detail.Table[name]] = list
	}
}
