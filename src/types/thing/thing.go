package thing

import "github.com/qlova/ilang/src"

var Type = ilang.NewType("thing", "SHARE", "GRAB")

func init() {
	ilang.RegisterStatement(Type, ScanStatement)
	ilang.RegisterSymbol("{", ScanSymbol)
	ilang.RegisterExpression(ScanExpression)
	
	ilang.RegisterFunction("thing", ilang.Method(Type, true, "PUSH 0\nMAKE"))
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	ic.Scan('}')
	return Type
}

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	if token == "{" {
		ic.Scan('}')
		ic.ExpressionType = Type
		
		var tmp = ic.Tmp("user")
		ic.Assembly("ARRAY ", tmp)
		return tmp
	}
	return ""
}

/*
	Scan a usertype statement for example,
		usertype = newusertype
		usertype.element = value
*/
func ScanStatement(ic *ilang.Compiler) {
	//The name is the usertype we want to index.
	var name = ic.Scan(ilang.Name)
	var usertype = ic.GetVariable(name) //This is the specific user TYPE.
	
	var token string
	var value string
	
	//The specific usertype has not been defined yet.
	if usertype == Type {
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
	
	//TODO CLEAN UP THIS MESS!
	
	//Support indexing at any level.
	// eg. Monster.Pos.X = 4
	var index string
	for token = ic.Scan(0); token == "."; {
		index = ic.Scan(ilang.Name)
		if token = ic.Scan(0); token == "." {
			name = ic.IndexUserType(name, index)
			
			if ic.ExpressionType.User {
				usertype = ic.ExpressionType
			}
			ic.SetVariable(name, ic.ExpressionType) //This is required for setusertype to recognise.
			ic.SetVariable(name+".", ilang.Protected)
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
		if ic.GetVariable(value) != ilang.Undefined {
			ic.ExpressionType = ic.GetVariable(value)
		}
		ic.SetUserType(name, index, value)
	}
}
