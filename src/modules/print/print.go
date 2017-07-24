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

func List(ic *ilang.Compiler, name string) {
	ic.Scanners = append(ic.Scanners, ic.Scanner)
	
	ic.Scanner = &scanner.Scanner{}
	ic.Scanner.Init(strings.NewReader(`
		for value in `+name+`
			print(value)
		end
	`))
	ic.Scanner.Position.Filename = "print.go"
	ic.Scanner.Whitespace= 1<<'\t' | 1<<'\r' | 1<<' '
}

func ScanPrint(ic *ilang.Compiler) {
	ic.Scan('(')
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
	
	for {
		token := ic.Scan(0)
		if token != "," {
			if token != ")" {
				ic.RaiseError()
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
}
