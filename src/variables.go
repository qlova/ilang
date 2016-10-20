package main

func (ic *Compiler) AssembleVar(name string, value string) {
	ic.Assembly("%v %v", ic.ExpressionType.Push, value)
	ic.Assembly("%v %v", ic.ExpressionType.Pop, name)
	ic.SetVariable(name, ic.ExpressionType)
}
