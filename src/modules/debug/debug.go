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
			if ic.ExpressionType.SubType != nil {
				fmt.Println("\tSubtype: ", ic.ExpressionType.SubType.Name)
			}
			fmt.Println("[RAW]" , ic.ExpressionType)
		
		case "variable":
			name := ic.Scan(ilang.Name)
			variable := ic.GetVariable(name)
			
			fmt.Println("[DEBUG] Variable ", name)
			fmt.Println("\tType: ", variable.Name)
			if variable.SubType != nil {
				fmt.Println("\tSubtype: ", variable.SubType.Name)
			}
		case "function":
			name := ic.Scan(ilang.Name)
			function, ok := ic.DefinedFunctions[name]
				
			fmt.Println("[DEBUG] Function ", name)
			if !ok {
				fmt.Println("Does not exist!")
			} else {
				fmt.Println("\tArguments: ", len(function.Args))
				fmt.Println("\tReturns: ", len(function.Returns))
			}
		case "type":
			name := ic.Scan(ilang.Name)
			t, ok := ic.DefinedTypes[name]
			
			if !ok {
				fmt.Println("Does not exist!")
			}
			
			fmt.Println("[DEBUG] Type ", name)
			for subt, index := range t.Detail.Table {
				fmt.Println("\t"+subt, ic.LastDefinedType.Detail.Elements[index])
			}
			
		default:
			ic.RaiseError("Unknown debug mode '", mode, "'")
	}
}
