package ilang

/*
	Scan a usertype statement for example,
		usertype = newusertype
*/
func (ic *Compiler) ScanUserStatement() {
	var name = ic.Scan(0)
	var usertype = ic.GetVariable(name)
	var token string
	
	if usertype == User {
	//This means we are in a method and we are defining a usertype's type.
	
		if !ic.GetFlag(InMethod) {
			ic.RaiseError()
		}
		ic.Scan('=')
		var value = ic.ScanExpression()
		ic.SetUserType(ic.LastDefinedType.Name, name, value)
		return
	}
	
	//TODO CLEAN UP THIS MESS!
	
	//Support indexing at any level
	// eg. Monster.Pos.X = 4
	var index string
	for token = ic.Scan(0); token == "."; {
		index = ic.Scan(Name)
		if token = ic.Scan(0); token == "." {
			name = ic.IndexUserType(name, index)
			ic.SetVariable(name, ic.ExpressionType) //This is required for setusertype to recognise.
			ic.SetVariable(name+".", Protected)
		}
	}
	
	if _, ok := usertype.Detail.Table[index]; index != "" && ok {
	if maybelist := usertype.Detail.Elements[usertype.Detail.Table[index]]; 
			(maybelist == List || maybelist.List) && token == "+" {
		ic.Scan('=')
		
		list := ic.IndexUserType(name, index)
		var listtype = ic.ExpressionType
		
		value := ic.ScanExpression()
		
		if listtype == List && ic.ExpressionType.Push == "PUSH" {
			usertype.Detail.Elements[usertype.Detail.Table[index]] = Array
			ic.Assembly("PLACE ", list)
			ic.Assembly("PUT ", value)
			return
		}
		
		if listtype == List {
			typedlist := ic.ExpressionType
			typedlist.List = true
			typedlist.User = false
			listtype = typedlist
			usertype.Detail.Elements[usertype.Detail.Table[index]] = typedlist
		}
		
		ic.PutList(listtype, list, value)
		return
	}
	}
	
	var value string
	if token != "=" {
		value = ic.IndexUserType(name, index)
	
		var b = ic.ExpressionType
		ic.NextToken = token
		ic.Shunt(value)
		if ic.ExpressionType != Undefined {
			ic.RaiseError("blank expression!")
		}
		ic.ExpressionType = b
	
	} else {
		value = ic.ScanExpression()
	}


	//Set a usertype from within a method.
	if _, ok := ic.LastDefinedType.Detail.Table[name]; ic.GetFlag(InMethod) && ok {
		ic.SetUserType(ic.LastDefinedType.Name, name, value)
		
	} else if index == "" {
		if !usertype.Empty() {
			//TODO garbage collection.
			ic.Assembly("PLACE ", value)
			ic.Assembly("RENAME ", name)
		}
	} else {
		ic.SetUserType(name, index, value)
	}
}
