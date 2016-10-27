package main

import "fmt"

func NewFlag() Type {
	var t Type
	t.Int = TypeIota
	TypeIota++
	t.Name = "flag_"+fmt.Sprint(t.Int)
	return t
}

var Protected = NewFlag()
var Software = NewFlag()
var InFunction = NewFlag()
var InMethod = NewFlag()
var Issues = NewFlag()
var Issue = NewFlag()
var FirstCase = NewFlag()
var Loop = NewFlag()
var ForLoop = NewFlag()
var New = NewFlag()

var Unused = NewFlag()
var Used = NewFlag()
