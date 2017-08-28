package ilang

import "strings"
import "github.com/gedex/inflector"

var TypeIota int

type Type struct {
	Name, Push, Pop string
	Int int
	User bool
	List bool
	
	Super string
	
	Decimal bool
	
	Detail *UserType
	Interface *Interface
	SubType *Type
}

func (t Type) Equals(b Type) bool {

	if t.Name == "list" {
		if b.Name == "list" {
			if t.SubType == nil || b.SubType == nil {
				return true
			}
		}
	}

	if t.Name != b.Name {
		return false
	}
	if t.SubType != nil && b.SubType != nil {
		return t.SubType.Equals(*b.SubType)
	}
	if t.SubType != nil || b.SubType != nil {
		return false
	}
	return true
}

func (t Type) DefaultValue() string {
	switch t.Push {
		case "PUSH":
			return "0"
		case "SHARE":
			return "backup"
		case "RELAY":
			return "open"
	}
	return ""
}

func (t Type) IsUser() Type {
	if t.User {
		return t
	} else {
		return Undefined
	}
}


func (t Type) IsMatrix() Type {
	if t.Name == "matrix" {
		return t
	} else {
		return Undefined
	}
}

func (t Type) IsArray() Type {
	if t.Name == "array" {
		return t
	} else {
		return Undefined
	}
}

func (t Type) IsList() Type {
	if t.List {
		return t
	} else {
		return Undefined
	}
}

//TODO this will break something types.
func (t Type) Empty() bool { 
	return t.Detail != nil && len(t.Detail.Elements) == 0
}

type UserType struct {	
	Elements []Type
	Table map[string]int
	SubElements map[int]Type
}

func NewUserType(name string) Type {
	t := NewType(name, "SHARE", "GRAB")
	t.User = true
	t.Detail = new(UserType)
	t.Detail.Table = make(map[string]int)
	return t
}

var string2type = map[string]Type{}

func NewType(name string, options ...string) Type {
	var t Type
	t.Name = name
	
	if len(options) == 2 {
		t.Pop = options[1]
		t.Push = options[0]
	}
	
	t.Int = TypeIota
	TypeIota++
	
	string2type[name] = t
	
	return t
}

func GetType(name string) Type {
	return string2type[name]
}

var Undefined = NewType("undefined")
var Number = NewType("number", "PUSH", "PULL")

var Text = NewType("text", "SHARE", "GRAB")
var Array = NewType("array", "SHARE", "GRAB")
var Matrix = NewType("matrix", "SHARE", "GRAB")

var Variadic = NewFlag()

func (ic *Compiler) ScanSymbolicType() Type {
	var result Type = Undefined
	var symbol = ic.Scan(0)
	
	if symbol == "." && ic.Peek() == "." {
		ic.Scan('.')
		symbol = ".."
	}
	
	if f, ok := Symbols[symbol]; ok {
		return f(ic)
	}
	
	switch symbol {
		case `""`:
			result = Text
			
		default:
			result = Number
			ic.NextToken = symbol
			return result
	}
	return result
}

//Check if the given type exists or not.
func (ic *Compiler) TypeExists(name string) bool {
	_, ok := ic.DefinedTypes[name]
	return ok
}

func (ic *Compiler) CallType(name string) string {
	if name == "text" {
		var array = ic.Tmp("user")
		ic.Assembly("ARRAY ", array)
		return array
	}
	
	if ic.DefinedTypes[name].Empty() {
		return ""
	} else {
		var array = ic.Tmp("user")
		ic.Assembly("ARRAY ", array)
		for range ic.DefinedTypes[name].Detail.Elements {
			ic.Assembly("PUT 0")
		}
		return array
	}
}

func (list Type) GetComplexName() string {
	if list.SubType == nil {
		return list.Name
	}
	
	if *list.SubType == list {
		panic("SELF REFERENCING LIST!")
	}

	var strings []string
	var t = list
	for {
		strings = append(strings, t.Name)
		if t.SubType == nil || (t == Type{}) {
			break
		}
		t = *t.SubType
	}
	var serialised string
	for i:=len(strings)-1; i >= 0; i-- {
		serialised += strings[i]
	}
	return serialised
}	

//Get a numeric value which represents the type.
func (ic *Compiler) GetPointerTo(name string) string {
	if ic.ExpressionType.Push == "SHARE" {
		var pointer = ic.Tmp("pointer")
		ic.Assembly("PUSH 0")
		ic.Assembly("SHARE ", name)
		ic.Assembly("HEAP")
		ic.Assembly("PULL ", pointer)
		return pointer
	}
	return name
}

func (ic *Compiler) Dereference(pointer string) string {
	if ic.ExpressionType.Push == "SHARE" {
		var value = ic.Tmp("deref")
		ic.Assembly("IF ", pointer)
		ic.Assembly("PUSH ", pointer)
		ic.Assembly("HEAP")
		ic.Assembly("ELSE ")
		
		var tmp = ic.CallType(ic.ExpressionType.Name)
		ic.Assembly("SHARE ", tmp)
		
		ic.Assembly("END")
		ic.Assembly("GRAB ", value)
		return value
	}
	return pointer
}

//This scans a new type definition and creates the type.
//eg. type Point { x, y }
func (ic *Compiler) ScanType() {
	var name = ic.Scan(Name)
	
	//This is for the grate engine.
	//Are we declaring a game?
	if name == "Game" {
		ic.Game = true
	}
	
	t := NewUserType(name)
	ic.DefinedTypes[name] = t
	
	switch ic.Scan(0) {
		case "{":
		case "is": //Inheritance eg. type WeightedPoint is Point { weight }
			super := ic.Scan(Name)
			t = ic.DefinedTypes[super]
			t.Super = t.Name
			t.Name = name
			switch ic.Scan(0) {
				case "\n":
					ic.DefinedTypes[name] = t
					ic.LastDefinedType = t
					return
				case "{":
				default:
					ic.RaiseError()
			}
		default:
			ic.RaiseError()
	}
		
	ic.InsertPlugins(name)
	//What are the elements?
	for {
		var token = ic.Scan(0)
		if token == "}" {
			break
		}
		if token != "," && token != "\n" {
			ic.NextToken = token
		}
		
		MemberType := ic.ScanSymbolicType()
		
		ident := ic.Scan(Name)
		if ident == "}" {
			break
		}
		
		//Embedded structs which are inferred.
		//eg.
		/*
			type Member {}
			type Base {
				Member() //This will be accessed as 'member'.
			}
		*/
		if ic.Peek() == "(" {
			if MemberType != Number {
				ic.RaiseError("Unexpected (")
			}
			ic.Scan('(')
			ic.Scan(')')
			var ok bool
			MemberType, ok = ic.DefinedTypes[ident]
			if !ok {
				ident := inflector.Singularize(ident)
				MemberType, ok = ic.DefinedTypes[ident]
				MemberType.List = true
				MemberType.User = false
				if !ok {
					ic.RaiseError("No such type! ", ident)
				}
			}
			ident = strings.ToLower(ident)
		}
		
		
		if ident != "\n" { 
			t.Detail.Elements = append(t.Detail.Elements, MemberType)
			t.Detail.Table[ident] = len(t.Detail.Elements)-1
		}
		
	}
	
	ic.DefinedTypes[name] = t

	ic.LastDefinedType = t
}

//This scans a type literal.
// eg. 
/*
	type Object {value}
	
	software { var o = Object{22} }
*/
func (ic *Compiler) ScanTypeLiteral() string {
	return ic.ScanConstructor()
}

func (ic *Compiler) ScanConstructor() string {
	var name = ic.Scan(Name)
					
	if _, ok := ic.DefinedTypes[name]; !ok {
		ic.RaiseError(name+" is an unrecognised type!")
	}
	
	var token = ic.Scan(0)
	
	/*if ic.Peek() == ")" && token == "(" {
		ic.ExpressionType = InFunction
		ic.NextToken = "("
		return name
	}*/
	
	var array = ic.Tmp("constructor")
	
	ic.Assembly("ARRAY ", array)
	//This is effectively a constructor.
	if token == "{" {
		var i int
		for {
			
			var expr = ic.ScanExpression()
			ic.Assembly("PLACE ", array)
			if ic.ExpressionType.Push == "PUSH" {
				ic.Assembly("PUT %v", expr)
			} else {
				var tmp = ic.Tmp("heap")
				ic.Assembly(ic.ExpressionType.Push," ", expr)
				ic.Assembly("PUSH 0")
				if ic.ExpressionType.Push == "RELAY" {
					ic.Assembly("HEAPIT")
				} else {
					ic.Assembly("HEAP")
				}
				ic.Assembly("PULL ", tmp)
				ic.Assembly("PUT ", tmp)
			}
			if i >= len(ic.DefinedTypes[name].Detail.Elements) {
				ic.RaiseError("Too many arguments passed to constructor!")
			}
			if ic.ExpressionType != ic.DefinedTypes[name].Detail.Elements[i] {
				ic.RaiseError("Mismatched types! Argument (%v) of constructor should be '%v'", i+1, 
					ic.DefinedTypes[name].Detail.Elements[i])
			}
			token = ic.Scan(0)
			for token == "\n" {
				token = ic.Scan(0)
			}
			if token == "}" {
				break
			} else if token != "," {
				ic.Expecting(",")
			}
			i++
		}
		for j := range ic.DefinedTypes[name].Detail.Elements {
			if j > i {
				ic.Assembly("PUT 0")
			}
		}
	} else if token == "\n" || token == ")" {
		for range ic.DefinedTypes[name].Detail.Elements {
			ic.Assembly("PUT 0")
		}
		if token == ")" {
			ic.NextToken = ")"
		}	
	} else {
		ic.RaiseError()
	}
	ic.ExpressionType = ic.DefinedTypes[name]
	return array
}
