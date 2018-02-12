package p

import "github.com/qlova/ilang/src"
import "github.com/qlova/ilang/src/modules/method"
import _ "github.com/qlova/ilang/src/modules/for"
import "github.com/qlova/ilang/src/types/list"

import (
	"strings"
	"text/scanner"
)

func init() {
	ilang.RegisterToken([]string{
		"print", 		//English 
		"afdrukken", 	//Dutch
		"印刷", 			//Chinese
		"Распечатать", 	//Russian
		"打印",			//Japanese
	}, ScanPrint)
}

//TODO change to use plugins instead of hacky scanner method.
func List(ic *ilang.Compiler, name string) {
	ic.Scanners = append(ic.Scanners, ic.Scanner)
	
	ic.Scanner = &scanner.Scanner{}
	ic.Scanner.Init(strings.NewReader(`
		for id, value in `+name+`
			write(value)
			if id < len(`+name+`)-1
				write(',')
			end
		end
		write('\n')
	`))
	ic.Scanner.Position.Filename = "print.go"
	ic.Scanner.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
}

func ScanPrint(ic *ilang.Compiler) {
	ic.Scan('(')
	ic.DisableOwnership = true
	arg := ic.ScanExpression()
	
	if (ic.ExpressionType == ilang.Undefined) {
		ic.RaiseError(ic.LastToken, " is undefined!")
	}	
	
	if !ic.ExpressionType.Empty() {
		ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
	}
	
	if ic.ExpressionType.Name == list.Type.Name || ic.ExpressionType == ilang.Array {
		ic.SetVariable(arg, ic.ExpressionType)
		List(ic, arg)
	} else {
		method.Call(ic, "text", ic.ExpressionType)
		ic.Assembly("STDOUT")
	}
	
	for {
		token := ic.Scan(0)
		if token != "," {
			if token != ")" {
				ic.RaiseError("Unexpected ", token)
			}
			break
		}
		arg := ic.ScanExpression()
		if !ic.ExpressionType.Empty() {
			ic.Assembly("%v %v", ic.ExpressionType.Push, arg)
		}
		if ic.ExpressionType.Name == list.Type.Name {
			List(ic, arg)
		} else {
			method.Call(ic, "text", ic.ExpressionType)
			ic.Assembly("STDOUT")
		}
	}
	
	ic.Assembly("SHARE i_newline")
	ic.Assembly("STDOUT")
	ic.DisableOwnership = false
}
