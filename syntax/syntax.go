package i

import "github.com/qlova/uct/compiler"

import (
	"github.com/qlova/ilang/syntax/symbols"

	"github.com/qlova/ilang/syntax/software"
	"github.com/qlova/ilang/syntax/print"
	"github.com/qlova/ilang/syntax/convert"
	"github.com/qlova/ilang/syntax/concept"
	"github.com/qlova/ilang/syntax/statement"
	"github.com/qlova/ilang/syntax/expression"
	"github.com/qlova/ilang/syntax/read"
	"github.com/qlova/ilang/syntax/send"
	"github.com/qlova/ilang/syntax/forloop"
	"github.com/qlova/ilang/syntax/ifelse"
	"github.com/qlova/ilang/syntax/error"
	"github.com/qlova/ilang/syntax/errors"
	"github.com/qlova/ilang/syntax/load"
	"github.com/qlova/ilang/syntax/open"
	"github.com/qlova/ilang/syntax/loop"
	"github.com/qlova/ilang/syntax/booleans"
	"github.com/qlova/ilang/syntax/importation"
	"github.com/qlova/ilang/syntax/fixed"
	"github.com/qlova/ilang/syntax/native"
	"github.com/qlova/ilang/syntax/binary"
	"github.com/qlova/ilang/syntax/typer"
	"github.com/qlova/ilang/syntax/languages"
	"github.com/qlova/ilang/syntax/create"
	"github.com/qlova/ilang/syntax/maths"
	"github.com/qlova/ilang/syntax/global"
	
	"github.com/qlova/ilang/syntax/graphics"
	"github.com/qlova/ilang/syntax/update"
	"github.com/qlova/ilang/syntax/input"

	"github.com/qlova/ilang/types/function"
	"github.com/qlova/ilang/types/text"
	"github.com/qlova/ilang/types/letter"
	"github.com/qlova/ilang/types/number"
	"github.com/qlova/ilang/types/decimal"
	"github.com/qlova/ilang/types/list"
	"github.com/qlova/ilang/types/array"
	"github.com/qlova/ilang/types/connection"
	"github.com/qlova/ilang/types/metatype"
	concept_type "github.com/qlova/ilang/types/concept"
)

//Pic'n Mix the 'i' language syntax!
func Syntax() compiler.Syntax {
	var syntax = compiler.NewSyntax("i")
	
	syntax.RegisterStatement(software.Statement)
	syntax.RegisterStatement(software.Exit)
	syntax.RegisterStatement(software.End)
	syntax.RegisterStatement(print.Statement)
	syntax.RegisterStatement(send.Statement)
	syntax.RegisterStatement(read.Statement)
	syntax.RegisterStatement(statement.Statement)
	syntax.RegisterStatement(concept.Return)
	syntax.RegisterStatement(concept.Statement)
	syntax.RegisterStatement(forloop.Statement)
	syntax.RegisterStatement(forloop.End)
	syntax.RegisterStatement(ifelse.If.Statement)
	syntax.RegisterStatement(ifelse.Else.Statement)
	syntax.RegisterStatement(ifelse.End)
	syntax.RegisterStatement(loop.Statement)
	syntax.RegisterStatement(loop.End)
	syntax.RegisterStatement(loop.Break)
	syntax.RegisterStatement(importation.Statement)
	syntax.RegisterStatement(fixed.Statement)
	syntax.RegisterStatement(forloop.Remove)
	syntax.RegisterStatement(convert.Statement)
	syntax.RegisterStatement(error.Statement)
	syntax.RegisterStatement(native.Statement)
	syntax.RegisterStatement(typer.Statement)
	syntax.RegisterStatement(languages.Statement)
	syntax.RegisterStatement(open.Statement)
	syntax.RegisterStatement(errors.Statement)
	syntax.RegisterStatement(errors.End)
	syntax.RegisterStatement(create.Statement)
	syntax.RegisterStatement(global.Statement)
	
	syntax.RegisterStatement(update.Statement)
	
	syntax.RegisterStatement(graphics.Display)
	syntax.RegisterStatement(graphics.Graphics)
	
	syntax.RegisterStatement(function.Statement)
	syntax.RegisterStatement(number.Statement)
	syntax.RegisterStatement(letter.Statement)
	syntax.RegisterStatement(list.Statement)
	syntax.RegisterStatement(text.Statement)
	syntax.RegisterStatement(array.Statement)
	syntax.RegisterStatement(decimal.Statement)
	syntax.RegisterStatement(concept_type.Statement)
	//syntax.RegisterStatement(metatype.Statement)

	syntax.RegisterExpression(expression.Expression)
	syntax.RegisterExpression(expression.NumberOf)
	syntax.RegisterExpression(convert.Expression)
	syntax.RegisterExpression(function.Expression)
	syntax.RegisterExpression(read.Expression)
	syntax.RegisterExpression(load.Expression)
	syntax.RegisterExpression(open.Expression)
	syntax.RegisterExpression(booleans.True)
	syntax.RegisterExpression(booleans.False)
	syntax.RegisterExpression(booleans.Not)
	syntax.RegisterExpression(binary.Expression)
	syntax.RegisterExpression(expression.Negative)
	syntax.RegisterExpression(print.Expression)
	syntax.RegisterExpression(metatype.Expression)
	syntax.RegisterExpression(create.Expression)
	syntax.RegisterExpression(global.Expression)
	syntax.RegisterExpression(maths.Abs)
	syntax.RegisterExpression(input.Keys)
	
	syntax.RegisterExpression(number.Expression)
	syntax.RegisterExpression(error.Expression)
	syntax.RegisterExpression(letter.Expression)
	syntax.RegisterExpression(text.Expression)
	syntax.RegisterExpression(list.Expression)
	syntax.RegisterExpression(array.Expression)
	syntax.RegisterExpression(decimal.Expression)
	syntax.RegisterExpression(concept_type.Expression)
	
	syntax.RegisterType(text.Type)
	syntax.RegisterType(letter.Type)
	syntax.RegisterType(number.Type)
	syntax.RegisterType(decimal.Type)
	syntax.RegisterType(function.Type)
	syntax.RegisterType(list.Type)
	syntax.RegisterType(array.Type)
	syntax.RegisterType(connection.Type)
	syntax.RegisterType(concept_type.Type)
	syntax.RegisterType(metatype.Type)
	
	syntax.RegisterFunction(&number.Method)
	syntax.RegisterFunction(&decimal.Method)
	
	syntax.RegisterFunction(&text.Itoa)
	syntax.RegisterFunction(&text.Join)
	syntax.RegisterFunction(&list.Copy)
	
	syntax.RegisterOperator(symbols.FunctionCallEnd, -1)
	syntax.RegisterOperator(symbols.IndexEnd, -1)
	syntax.RegisterOperator(symbols.Or, 0)
	syntax.RegisterOperator(symbols.And, 1)
	syntax.RegisterOperator(symbols.Equals, 2)
	syntax.RegisterOperator(symbols.More, 2)
	syntax.RegisterOperator(symbols.Less, 2)
	syntax.RegisterOperator(symbols.Plus, 3)
	syntax.RegisterOperator(symbols.Minus, 3)
	syntax.RegisterOperator(symbols.Times, 4)
	syntax.RegisterOperator(symbols.Power, 5)
	syntax.RegisterOperator(symbols.Divide, 4)
	syntax.RegisterOperator(symbols.Modulus, 4)

	syntax.RegisterOperator(symbols.FunctionCallBegin, 6)
	syntax.RegisterOperator(symbols.IndexBegin, 6)
	syntax.RegisterOperator(symbols.Index, 7)
	
	
	syntax.RegisterAlias("×", symbols.Times)
	syntax.RegisterAlias("÷", symbols.Divide)
	
	return syntax
}
