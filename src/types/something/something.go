package something

import "github.com/qlova/ilang/src"

var Type = ilang.NewType("something", "SHARE", "GRAB")

func init() {
	ilang.RegisterStatement(Type, ScanStatement)
	ilang.RegisterSymbol("?", ScanSymbol)
	ilang.RegisterExpression(ScanExpression)
	ilang.RegisterShunt(".", Shunt)
	
	ilang.RegisterFunction("something", ilang.Method(Type, true, "PUSH 1\nMAKE"))
}

func ScanSymbol(ic *ilang.Compiler) ilang.Type {
	
	if ic.Peek() == "{" {
		return ScanInterface(ic)
	}
	
	return Type
}

func Shunt(ic *ilang.Compiler, name string) string {
	if ic.ExpressionType == Type {
		var cast = ic.Scan(ilang.Name)
		return ic.Shunt(Index(ic, name, cast))
	}
	return ""
}

func ScanExpression(ic *ilang.Compiler) string {
	var token = ic.LastToken
	if token == "?" {
		ic.ExpressionType = Type
		var tmp = ic.Tmp("something")
		ic.Assembly("ARRAY ", tmp)
		ic.Assembly("PUT 0")
		return tmp
	}
	return ""
}

func ScanStatement(ic *ilang.Compiler) {

	var name = ic.Scan(ilang.Name)
	ic.Scan('=') //TODO Maybe support other types of statements?
	var value = ic.ScanExpression()
	
	Assign(ic, name, value)
	
}

func Index(ic *ilang.Compiler, name string, cast string) string {
	switch cast {
		case "number", "letter", "decimal":
			var test = ic.Tmp("test")
			ic.Assembly("PUSH 2")
			ic.Assembly("PLACE ", name)
			ic.Assembly("GET ", test)
			ic.Assembly("SEQ %v %v %v", test, test, ilang.GetType(cast).Int)
			ic.Assembly("IF ", test)
			
			var test2 = ic.Tmp("test")
			var num = ic.Tmp("number")
			ic.Assembly("PUSH 0")
			ic.Assembly("GET ", test2)
			ic.Assembly("PUSH ", test2)
			
			ic.Assembly("ELSE")
			ic.Assembly("ERROR 404")
			ic.Assembly("PUSH 0")
			ic.Assembly("END")
			ic.Assembly("PULL ", num)
			
			ic.ExpressionType = ilang.GetType(cast)
			if ic.ExpressionType == ilang.Undefined {
				ic.RaiseError("Cannot cast something to ", cast)
			}
			
			return num
		default:
			var test = ic.Tmp("test")
			ic.Assembly("PUSH 2")
			ic.Assembly("PLACE ", name)
			ic.Assembly("GET ", test)
			ic.Assembly("SEQ %v %v %v", test, test, ilang.GetType(cast).Int)
			ic.Assembly("IF ", test)
			
			var address = ic.Tmp("address")
			ic.Assembly("PUSH 0")
			ic.Assembly("PLACE ", name)
			ic.Assembly("GET ", address)
			ic.Assembly("PUSH ", address)
			ic.Assembly("HEAP")
			var txt = ic.Tmp("txt")
			
			ic.Assembly("ELSE")
			ic.Assembly("ERROR 404")
			ic.Assembly("ARRAY ", txt)
			ic.Assembly("SHARE ", txt)
			ic.Assembly("END")
			
			ic.Assembly("GRAB ", txt)
			
			ic.ExpressionType = ilang.GetType(cast)
			if ic.ExpressionType == ilang.Undefined {
				ic.RaiseError("Cannot cast something to ", cast)
			}
			
			return txt
	}
	return ""
}

func Assign(ic *ilang.Compiler, name string, value string) {
	switch ic.ExpressionType.Push {
	
		case "PUSH":
			var tmp = ic.Tmp("number")
			ic.Assembly("ARRAY ", tmp)
			ic.Assembly("PUT ", value)
			ic.Assembly("PUT 0")
			ic.Assembly("PUT ", ic.ExpressionType.Int)
			for _, v := range ic.ExpressionType.Name {
				ic.Assembly("PUT ", byte(v))
			}
			ic.Assembly("SHARE ", name)
			ic.Assembly("RUN collect_m_something")
			ic.Assembly("PLACE ", tmp)
			ic.Assembly("RENAME ", name)
			
		case "SHARE":
			var tmp = ic.Tmp("text")
			ic.Assembly("ARRAY ", tmp)
			ic.Assembly("SHARE ", value)
			ic.Assembly("PUSH 0")
			ic.Assembly("HEAP")
			var address = ic.Tmp("address")
			ic.Assembly("PULL ", address)
			ic.Assembly("PUT ", address)
			ic.Assembly("PUT 1")
			ic.Assembly("PUT ", ic.ExpressionType.Int)
			for _, v := range ic.ExpressionType.Name {
				ic.Assembly("PUT ", byte(v))
			}
			ic.Assembly("SHARE ", name)
			ic.Assembly("RUN collect_m_something")
			ic.Assembly("PLACE ", tmp)
			ic.Assembly("RENAME ", name)
			
		default:
			ic.RaiseError(ic.ExpressionType.Name, " is not a something")
	}
}
