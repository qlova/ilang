package debug

import "github.com/qlova/ilang/src"
import "fmt"

func init() {
	ilang.RegisterToken([]string{"debug"}, ScanDebug)
}

func ScanDebug(ic *ilang.Compiler) {
	var mode = ic.Scan(0)
	
	
	switch mode {
		case "expression":	
			ic.ScanExpression()

			fmt.Println("[DEBUG] Expression")
			fmt.Println("\tType: ", ic.ExpressionType.Name)
		
		case "variable":
			name := ic.Scan(ilang.Name)
			variable := ic.GetVariable(name)
			
			fmt.Println("[DEBUG] Variable ", name)
			fmt.Println("\tType: ", variable.Name)
			if variable.SubType != nil {
				fmt.Println("\tSubtype: ", variable.SubType.Name)
			}
		
		default:
			ic.RaiseError("Unknown debug mode '", mode, "'")
	}
}
