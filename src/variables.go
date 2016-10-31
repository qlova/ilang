package main

import "strings"

//Create the assembly for a new variable and keep track of it.
//Type is inferred by the Compiler's ExpressionType value.
func (ic *Compiler) CreateVariable(name, value string) {
	ic.AssembleVar(name, value)
}

//Assign a value to a variable, this will generate assembly.
func (ic *Compiler) AssembleVar(name string, value string) {

	ic.Assembly("%v %v", ic.ExpressionType.Push, value)
	ic.Assembly("%v %v", ic.ExpressionType.Pop, name)
	ic.SetVariable(name, ic.ExpressionType)
	
	list := ic.ExpressionType
	list.List = true
	list.User = false
	if ic.GetFlag(InMethod) && ic.ExpressionType != List &&
		ic.LastDefinedType.Detail.Elements[ic.LastDefinedType.Detail.Table[name]] == List {

		ic.LastDefinedType.Detail.Elements[ic.LastDefinedType.Detail.Table[name]] = list
	}
}

//Set the type of a variable, this is akin to creating or assigning a variable.
func (c *Compiler) SetVariable(name string, sort Type) {
	if !strings.Contains(name, "_") && sort != Protected {
		c.SetVariable(name+"_use", Unused)
		for i:=len(c.Scope)-1; i>=0; i-- {
			if v, ok := c.Scope[i][name]; ok && v != List && v != User {
				c.RaiseError("Duplicate variable name!", name, "(", v.Name, ")")
			}
		}
	}
	c.Scope[len(c.Scope)-1][name] = sort
}

//This will return the type of the variable. UNDEFINED for undefined variables.
func (ic *Compiler) GetVariable(name string) Type {
	for i:=len(ic.Scope)-1; i>=0; i-- {
		if v, ok := ic.Scope[i][name]; ok {
			ic.Scope[i][name+"_use"] = Used
			return v
		}
	}
	
	//Allow table values to be indexed in a method.
	if ic.GetFlag(InMethod) {
		if _, ok := ic.LastDefinedType.Detail.Table[name]; ok {
			var value = ic.IndexUserType(ic.LastDefinedType.Name, name)
			ic.AssembleVar(name, value)
			ic.SetVariable(name+"_use", Used)
			return ic.ExpressionType
		}
	}
	
	return Undefined
}
