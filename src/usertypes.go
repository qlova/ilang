package ilang

/*
	Scan a usertype statement for example,
		usertype = newusertype
		usertype.element = value
*/
/*
func (ic *Compiler) ScanUserStatement() {
	//The name is the usertype we want to index.
	var name = ic.Scan(0)
	var usertype = ic.GetVariable(name) //This is the specific user TYPE.
	
	var token string
	var value string
	
	//The specific usertype has not been defined yet.
	if usertype == User {
	
		//This means we are in a method and we are defining a usertype's type.
		if ic.GetFlag(InMethod) {
		
			ic.Scan('=')
			var value = ic.ScanExpression()
			ic.SetUserType(ic.LastDefinedType.Name, name, value)
			return
		
		
		//Here, we set the type for a unset user variable.
		//eg. var u = {}; u = Type()
		} else {
			
			ic.Scan('=')
			var value = ic.ScanExpression()
			if !ic.ExpressionType.User {
				ic.RaiseError("Cannot set ", name, " to type ", ic.ExpressionType.Name, "! Must be a user defined type!")
			}
			ic.Assembly("PLACE ", value)
			ic.Assembly("RENAME ", name)
			
			ic.UpdateVariable(name, ic.ExpressionType)
			
			return
		}
	}
	
	//TODO CLEAN UP THIS MESS!
	
	//Support indexing at any level.
	// eg. Monster.Pos.X = 4
	var index string
	for token = ic.Scan(0); token == "."; {
		index = ic.Scan(Name)
		if token = ic.Scan(0); token == "." {
			name = ic.IndexUserType(name, index)
			
			if ic.ExpressionType.User {
				usertype = ic.ExpressionType
			}
			ic.SetVariable(name, ic.ExpressionType) //This is required for setusertype to recognise.
			ic.SetVariable(name+".", Protected)
		}
	}
	
	//Figure out the value which needs to be assigned.
	if token == "=" {
		value = ic.ScanExpression()
	} else {
		value = ic.IndexUserType(name, index)
		ic.SetVariable(value, ic.ExpressionType)
		ic.NextToken = value
		ic.NextNextToken = token
		ic.ScanStatement()
	}
	
	//Function and pipe calls do not need to be reassigned. 
	if token == "(" || (token == "-" && ic.Peek() == "-") {
		return
	}

	//Assign the value to the userdata.
	if index == "" {
		if !usertype.Empty() {
			//TODO garbage collection.
			ic.Assembly("PLACE ", value)
			ic.Assembly("RENAME ", name)
		}
	} else {
		if ic.GetVariable(value) != Undefined {
			ic.ExpressionType = ic.GetVariable(value)
		}
		ic.SetUserType(name, index, value)
	}
	
	if _, ok := usertype.Detail.Table[index]; index != "" && ok {
	if maybelist := usertype.Detail.Elements[usertype.Detail.Table[index]]; 
			(maybelist == List || maybelist.List) && (token == "+" || token == "-") {
			
		list := ic.IndexUserType(name, index)
		var listtype = ic.ExpressionType	
		
		if token == "-" {
			ic.Scan('-')
			//TODO garbage collect.
			ic.Assembly("PLACE ", list)
			ic.Assembly("IF #",list,"\nPOP ", ic.Tmp("emptypop"), "\nEND")
			
			return
		}
		ic.Scan('=')
		
		value = ic.ScanExpression()
		
		if listtype == List && ic.ExpressionType.Push == "PUSH" {
			usertype.Detail.Elements[usertype.Detail.Table[index]] = Array
			ic.Assembly("PLACE ", list)
			ic.Assembly("PUT ", value)
			value = list
			ic.ExpressionType = Array
			ic.SetUserType(name, index, value)
			return
			
		} else {
		
			if listtype == List {
				typedlist := ic.ExpressionType
				typedlist.List = true
				typedlist.User = false
				listtype = typedlist
				usertype.Detail.Elements[usertype.Detail.Table[index]] = typedlist
			}
		
			ic.PutList(listtype, list, value)
		}
	}
	} else {
	
		value = ""
		//Modify the usertype value.
		if token != "=" {
			value = ic.IndexUserType(name, index)
			ic.SetVariable(value, ic.ExpressionType)
			ic.NextToken = value
			ic.NextNextToken = token
			ic.ScanStatement()
		} else {
			value = ic.ScanExpression()
		}
	
		println(token)
	}

	if token == "=" || token == "+" || token == "-" || token == "/" || token == "*" {
		println("Gonna set ", name, "[", index, "] =", value, " ", ic.GetVariable(value).Name)
		
	
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
	}*/
//}

func (ic *Compiler) IndexUserType(name, element string) string {
	var t UserType
	if ic.GetVariable(name) != Undefined {
		t = *ic.GetVariable(name).Detail
		ic.SetVariable(name+"_use", Used)
	} else {
		t = *ic.ExpressionType.Detail
	}
	
	//Deal with indexing Something types.
	/*if GetVariable(name) == SOMETHING {
		switch element {
			case "number":
				ExpressionType = NUMBER
				fmt.Fprintf(output, "PLACE %s\n", name)
				fmt.Fprintf(output, "PUSH 0\n")
				fmt.Fprintf(output, "GET %s%v\n", "i+user+", unique)
				return "i+user+"+fmt.Sprint(unique)
		}
	}*/
	
	if index, ok := t.Table[element]; !ok {
		ic.RaiseError(name+" does not have an element named "+element)
	} else {
	
		var tmp = ic.Tmp("index")
		ic.ExpressionType = t.Elements[index]
	
		switch t.Elements[index].Push {
			case "PUSH":
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUSH ", index)
				ic.Assembly("GET ", tmp)
				return tmp
			
			case "SHARE", "RELAY":
				ic.Assembly("PLACE ", name) //The array we are indexing, the place.
				ic.Assembly("PUSH ", index) //Push the index onto the stack.
				ic.Assembly("GET ", tmp)	//Get the value of the array at the index on the stack.
				ic.Assembly("IF ",tmp)		//If there is a valid address, (greater than zero)
				ic.GainScope()
				
				//Retrieve the array.
				ic.Assembly("PUSH ", tmp)
				if t.Elements[index].Push == "RELAY" {
					ic.Assembly("HEAPIT")
				} else {
					ic.Assembly("HEAP")
				}
				tmp = ic.Tmp("index")
				ic.Assembly(t.Elements[index].Pop, " ", tmp)
				ic.Assembly(t.Elements[index].Push, " ", tmp)
				ic.LoseScope()
				
				ic.Assembly("ELSE") //We will return a new array.
				ic.GainScope()
				ic.Assembly("ARRAY ", tmp)
				if t.Elements[index].User {
				for range t.Elements[index].Detail.Elements {
					ic.Assembly("PUT 0")
				}
				}
				
				//First we will assign it to the usertype.
				if t.Elements[index].Push == "SHARE" {
					var tmp2 = ic.Tmp("index")
					
					ic.Assembly("SHARE ", tmp)
					ic.Assembly("PUSH 0")
					ic.Assembly("HEAP")
					ic.Assembly("PULL ", tmp2)
				
					ic.Assembly("PLACE ", name)
					ic.Assembly("PUSH ", index)
					ic.Assembly("SET ", tmp2) 
				}
				
				ic.Assembly("SHARE ", tmp)
				if t.Elements[index].Push == "RELAY" {
					ic.Assembly("OPEN")
				}
				ic.LoseScope()
				ic.Assembly("END")
				ic.Assembly(t.Elements[index].Pop, " ", tmp)
				
				return tmp
				
			default:
				ic.RaiseError(name+" cannot index "+element+", type is unindexable!!!")
		}
	}
	return ""
}

func (ic *Compiler) SetUserType(name, element, value string) {
	var t UserType
	if ic.GetVariable(name) != Undefined {
		t = *ic.GetVariable(name).Detail
		ic.SetVariable(name+"_use", Used)
	} else {
		ic.RaiseError("Cannot set type without type identity!")
	}
	
	if index, ok := t.Table[element]; !ok {
		ic.RaiseError(name+" does not have an element named "+element)
	} else {
	
		if (t.Elements[index].Name == "thing") || ic.ExpressionType.Name == "matrix" {
			t.Elements[index] = ic.ExpressionType
			
		}
	
		if !ic.ExpressionType.Equals(t.Elements[index]) {
			ic.RaiseError("Type mismatch, cannot assign '",ic.ExpressionType.Name,"', to a element of type '",t.Elements[index].Name,"'")		
		}

		switch t.Elements[index].Push {
			case "PUSH":
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUSH ", index)
				ic.Assembly("SET ", value)
			
			case "SHARE", "RELAY":
				
				//TODO garbage collect
				var tmp = ic.Tmp("index")
				ic.Assembly(t.Elements[index].Push, " ", value)
				ic.Assembly("PUSH 0")
				if t.Elements[index].Push == "RELAY" {
					ic.Assembly("HEAPIT")
				} else {
					ic.Assembly("HEAP")
				}
				ic.Assembly("PULL ", tmp)
				
				ic.Assembly("PLACE ", name)
				ic.Assembly("PUSH ", index)
				ic.Assembly("SET ", tmp)
				
			default:
				ic.RaiseError(name+" cannot index "+element+", type is unindexable!!!")
		}
	}
}
