package connection

import "github.com/qlova/uct/compiler"

var Name = compiler.Translatable{
	compiler.English: "connection",
}

var Type = compiler.Type {
	Name: Name,
	
	Base: compiler.PIPE,
}
