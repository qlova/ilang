package ilang

import "strings"

var OnVariableMarked = make([]func(*Compiler, string, string), 0, 4)

//OnVariableMarked callback, run when a variable is marked.
func RegisterOnVariableMarked(callback func(*Compiler, string, string)) {
	OnVariableMarked = append(OnVariableMarked, callback)
}

//Create the assembly for a new variable and keep track of it.
//Type is inferred by the Compiler's ExpressionType value.
func (ic *Compiler) CreateVariable(name, value string) {
	ic.AssembleVar(name, value)
}

//Assign a value to a variable, this will generate assembly.
func (ic *Compiler) AssembleVar(name string, value string) {

	if !ic.ExpressionType.Empty() {
		ic.Assembly("%v %v", ic.ExpressionType.Push, value)
	}
	
	if !ic.ExpressionType.Empty() {
		ic.Assembly("%v %v", ic.ExpressionType.Pop, name)
	}
	ic.SetVariable(name, ic.ExpressionType)
}

//Set the type of a variable, this is akin to creating or assigning a variable.
func (c *Compiler) SetVariable(name string, sort Type) {
	if !strings.Contains(name, "_") && sort != Protected {
		c.SetVariable(name+"_use", Unused)
		for i:=len(c.Scope)-1; i>=0; i-- {
			if v, ok := c.Scope[i][name]; ok /*&& v != List && v != User*/ {
				c.RaiseError("Duplicate variable name!", name, "(", v.Name, ")")
			}
		}
	}
	c.Scope[len(c.Scope)-1][name] = sort
}

//Mark the variable with a specific feature.
func (ic *Compiler) MarkVariable(name string, mark string) {
	for i:=len(ic.Scope)-1; i>=0; i-- {
		if _, ok := ic.Scope[i][name]; ok {
			ic.Scope[i][name+"_"+mark] = Used
			
			//Run the OnVariableMarked callbacks
			for _, callback := range OnVariableMarked {
				callback(ic, name, mark)
			}
			
		}
	}
}

//This will update a variable to be a new type.
func (ic *Compiler) UpdateVariable(name string, sort Type) {
	for i:=len(ic.Scope)-1; i>=0; i-- {
		if _, ok := ic.Scope[i][name]; ok {
			ic.Scope[i][name] = sort
		}
	}
}

//This will return the type of the variable. UNDEFINED for undefined variables.
func (ic *Compiler) GetVariable(name string) Type {
	for i:=len(ic.Scope)-1; i>=0; i-- {
		if v, ok := ic.Scope[i][name]; ok {
			ic.Scope[i][name+"_use"] = Used
			return v
		}
	}
	
	for _, f := range Variables {
		t := f(ic, name)
		if t != Undefined {
			return t
		}
	}
	
	return Undefined
}
