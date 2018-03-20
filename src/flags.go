package ilang

import "fmt"

func NewFlag() Type {
	var t Type
	t.Int = TypeIota
	TypeIota++
	t.Name = "flag_"+fmt.Sprint(t.Int)
	return t
}

func NewNamedFlag(name string) Type {
	var t Type
	t.Int = TypeIota
	TypeIota++
	t.Name = "flag_"+name
	return t
}

var Protected = NewNamedFlag("Protected")
var InFunction = NewNamedFlag("InFunction")
var InMethod = NewNamedFlag("InMethod")
var FirstCase = NewNamedFlag("FirstCase")

var Unused = NewNamedFlag("Unused")
var Used = NewNamedFlag("Used")

//This will return the value of a scopped flag.
func (c *Compiler) GetFlag(sort Type) bool {
	for i:=len(c.Scope)-1; i>=0; i-- {
		if _, ok := c.Scope[i][sort.Name]; ok {
			return true
		}
	}
	return false
}

//This will return the value of a scopped flag.
func (c *Compiler) GetScopedFlag(sort Type) bool {
	if _, ok := c.Scope[len(c.Scope)-1][sort.Name]; ok {
		return true
	}
	return false
}

//Set the type of a variable, this is akin to creating or assigning a variable.
func (c *Compiler) SetFlag(flag Type) {
	c.Scope[len(c.Scope)-1][flag.Name] = flag
}

//Set the type of a variable, this is akin to creating or assigning a variable.
func (c *Compiler) UnsetFlag(flag Type) {
	delete(c.Scope[len(c.Scope)-1], flag.Name) 
}
