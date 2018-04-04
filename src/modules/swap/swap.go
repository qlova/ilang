package swap

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/modules/method"

func init() {
	ilang.RegisterToken([]string{
		"swap", 		//English
	}, ScanSwap)
}

func ScanSwap(ic *ilang.Compiler) {
	ic.Scan('(')
	var a = ic.Scan(0)
	var t = ic.GetVariable(a)
	if t == ilang.Undefined {
		ic.RaiseError(a, " is not a defined variable!")
	}
	
	ic.Scan(',')
	var b = ic.Scan(0)
	
	if ic.GetVariable(b) == ilang.Undefined {
		ic.RaiseError(b, " is not a defined variable!")
	}
	if ic.GetVariable(b) != t {
		ic.RaiseError(a, " and ", b," are not the same type!")
	}
	
	ic.Scan(')')
	
	switch (t.Push) {
		case "PUSH":
		
			var tmp = ic.Tmp("swap")
			ic.Assembly("VAR %s", tmp)
		
			ic.Assembly("ADD %s %s 0", tmp, a)
			ic.Assembly("ADD %s %s 0", a, b)
			ic.Assembly("ADD %s %s 0", b, tmp)
		
		case "SHARE":

			var tmp = ic.Tmp("swap")
			ic.Assembly("ARRAY ", tmp)
			ic.Assembly("SHARE ", a)
			ic.Assembly("RENAME ", tmp)
			
			ic.Assembly("SHARE ", b)
			ic.Assembly("RENAME ", a)
			
			ic.Assembly("SHARE ", tmp)
			ic.Assembly("RENAME ", b)
		
		case "RELAY":
			var tmp = ic.Tmp("swap")
			ic.Assembly("ARRAY ", tmp)
			ic.Assembly("RELAY ", a)
			ic.Assembly("RELOAD ", tmp)
			
			ic.Assembly("RELAY ", b)
			ic.Assembly("RELOAD ", a)
			
			ic.Assembly("RELAY ", tmp)
			ic.Assembly("RELOAD ", b)
		default:
			ic.RaiseError("Cannot swap variables ",a, " and ", b, "... Unswappable types!")
	}
	
	if ic.GetFlag(method.Flag) {
		method.Sync(ic, a, b)
	}
}
