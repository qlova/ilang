package errors

import "github.com/qlova/uct/compiler"
import "strings"

func AssignmentMismatch(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Cannot assign a "+a.Name[compiler.English]+" value to a variable of type "+b.Name[compiler.English],
	}
}

func ExpectingType(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Expecting a value of type "+a.Name[compiler.English]+" instead got a value of type "+b.Name[compiler.English],
	}
}

func UnknownType(a string) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Unknown Type: "+a,
	}
}



func Single(a compiler.Type, symbol string, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "The relationship "+a.Name[compiler.English]+symbol+b.Name[compiler.English]+" is not defined!",
	}
}

func IsInvalidName(name string) bool {
	return strings.Contains(name, "_")
}


func InvalidName(name string) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "Invalid name: "+name,
	}
}

func Inconsistent(a, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "The usage here of the '"+a.Name[compiler.English]+"' type is inconsistent with the\n use of the '"+b.Name[compiler.English]+"' type before this!", 
	}
}

func NoSuchElement(a string, b compiler.Type) compiler.Translatable {
	return compiler.Translatable {
		compiler.English: "No such element '"+a+"' in type '"+b.Name[compiler.English]+"'!", 
	}
}
