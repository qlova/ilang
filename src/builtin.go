package ilang

//YES, this file is a mess. Maybe this needs to be bindata.


func BlankMethod(r Type) Function {
	return Function {
		Exists: true,
		Inline: true,
		Returns: []Type{r},
	}
}

func SimpleMethod(r Type, data string) Function {
	return Function{
		Exists: true,
		Returns: []Type{r},
		Data: data,
	}
}

func Method(r Type, Inline bool, Data string, load ...string) Function {
	if len(load) > 0 {
		return Function{
			Exists: true,
			Returns: []Type{r},
			Inline:Inline,
			Data:Data,
			Method:true,
			Import:load[0],
		}
	} else {
		return Function{
			Exists: true,
			Returns: []Type{r},
			Inline:Inline,
			Data:Data,
			Method:true,
		}
	}
}

func InlineFunction(a []Type, data string, r []Type) Function {
	return Function{
		Exists: true,
		Args: a,
		Returns: r,
		Data: data,
		Inline: true,
	}
}


func Alias(f string, r Type) Function {
	return Function{
		Exists: true,
		Inline: true,
		Data: "RUN "+f,
		Returns: []Type{r},
		Import: f,
	}
}

func (ic *Compiler) Builtin() {
	ic.DefinedFunctions["number"] = Method(Number, true, "PUSH 0")
	ic.DefinedFunctions["array"] = Method(Array, true, "PUSH 0\nMAKE")
	ic.DefinedFunctions["binary"] = Method(Number, true, "PUSH 0")
	ic.DefinedFunctions["binary"] = Method(Number, true, "PUSH 0")
	
	ic.DefinedFunctions["random"] = Method(Number, true, "PUSH 0")
	
	ic.DefinedFunctions["delete_m_text"] = Method(Text, true, "DELETE")
	
	ic.DefinedFunctions["load"] = Method(Text, true, "PUSH 0")
	ic.DefinedFunctions["text"] = Method(Number, false, `
FUNCTION text
	ARRAY a
	SHARE a
RETURN
`)
	
	//ic.DefinedFunctions["random_m_decimal"] = Alias("random_m_number", Decimal)

	ic.DefinedFunctions["copy"] = Method(Undefined, true, "")
	
	
	ic.DefinedFunctions["collect_m_text"] = BlankMethod(Number)
	ic.DefinedFunctions["number_m_letter"] = BlankMethod(Number)
	ic.DefinedFunctions["number_m_set"] = BlankMethod(Number)
	ic.DefinedFunctions["text_m_text"] = BlankMethod(Text)
	
	ic.DefinedFunctions["load"] = Method(Undefined, true, "")
	ic.DefinedFunctions["sort"] = Method(Undefined, true, "")
	ic.DefinedFunctions["open"] = Method(Undefined, true, "")
	ic.DefinedFunctions["trim"] = Method(Undefined, true, "")
	
	ic.DefinedFunctions["execute"] = InlineFunction([]Type{Text}, "EXECUTE", nil)
	ic.DefinedFunctions["delete"] = InlineFunction([]Type{Text}, "DELETE", nil)
	ic.DefinedFunctions["rename"] = InlineFunction([]Type{Text, Text}, "MOVE", nil)
	
	ic.DefinedFunctions["read"] = Method(Text, true, "PUSH 0\nSTDIN")
	
	ic.DefinedFunctions["link"] = InlineFunction([]Type{Text, Number}, "LINK", nil)
	ic.DefinedFunctions["connect"] = InlineFunction([]Type{Text}, "CONNECT", []Type{Number})
	
	ic.DefinedFunctions["load_m_text"] = Method(Text, true, "LOAD")
	
	ic.DefinedFunctions["text_m_number"] = Method(Text, true, "PUSH 10\nRUN i_base_number", "i_base_number")
	
	ic.DefinedFunctions["binary_m_number"] = Method(Text, true, "PUSH 2\nRUN i_base_number", "i_base_number")
	
	ic.DefinedFunctions["binary_m_text"] = Method(Text, true, "PUSH 2\nRUN i_base_string", "i_base_string")
	
	ic.DefinedFunctions["number_m_text"] = Method(Number, true, "PUSH 10\nRUN i_base_string", "i_base_string")
	
	ic.DefinedFunctions["inbox"] = InlineFunction(nil, "INBOX", []Type{Text})
	ic.DefinedFunctions["outbox"] = InlineFunction([]Type{Text}, "OUTBOX", nil)
	
	ic.DefinedFunctions["output"] = InlineFunction([]Type{Text}, "STDOUT", nil)
	ic.DefinedFunctions["output_m_pipe"] = InlineFunction([]Type{Text}, "OUT", nil)
	
	ic.DefinedFunctions["len"] = Method(Undefined, true, "")
	ic.DefinedFunctions["len_m_list"] = InlineFunction([]Type{Array}, "LEN", nil)
	ic.DefinedFunctions["len_m_text"] = InlineFunction([]Type{Text}, "LEN", nil)
	ic.DefinedFunctions["len_m_number"] = Method(Array, true, "MAKE")
	ic.DefinedFunctions["array_m_number"] = Method(Array, true, "MAKE")
	
	ic.DefinedFunctions["len_m_pipe"] = Function{Exists:true, Returns:[]Type{Number}, Data: `
FUNCTION len_m_pipe
	ARRAY a
	SHARE a
	STAT
	GRAB s
	PLACE s
	POP n
	PUSH n
RETURN
`}
		
	ic.DefinedFunctions["load_m_number"] = Function{Exists:true, Returns:[]Type{Text}, Data: `
FUNCTION load_m_number
	ARRAY a
	PULL b
	PUT b
	SHARE a
	LOAD
RETURN
`}

	ic.DefinedFunctions["load_m_letter"] = Function{Exists:true, Returns:[]Type{Text}, Data: `
FUNCTION load_m_letter
	ARRAY a
	PULL b
	MUL b b -1
	PUT b
	SHARE a
	LOAD
RETURN
`}

	ic.DefinedFunctions["text_m_letter"] = Function{Exists:true, Returns:[]Type{Text}, Data: `
FUNCTION text_m_letter
	ARRAY a
	PULL letter
	PUT letter
	SHARE a
RETURN
`}

ic.DefinedFunctions["open_m_letter"] = Function{Exists:true, Returns:[]Type{Number}, Data: `
FUNCTION open_m_letter
	ARRAY a
	PULL b
	MUL b b -1
	PUT b
	SHARE a
	ADD ERROR 0 0
	LOAD
	IF ERROR
		PUSH 0
	ELSE
		PUSH 1
	END
RETURN
`}
	
	ic.DefinedFunctions["read_m_letter"] = Function{Exists:true, Returns:[]Type{Text}, Data: `
FUNCTION read_m_letter
	PULL delim
	MUL delim delim -1
	PUSH delim
	STDIN
RETURN
`}

ic.DefinedFunctions["reada_m_pipe"] = Function{Exists:true, Returns:[]Type{Text}, Data:`
FUNCTION reada_m_pipe
	PULL delim
	MUL delim delim -1
	PUSH delim
	IN
RETURN
`}

ic.DefinedFunctions["text_m_something"] = Function{Exists:true, Returns:[]Type{Text}, Data: `
FUNCTION text_m_something
	GRAB something
	PUSH #something
	PUSH 3
	SHARE something
	SLICE
RETURN
`}
	
	ic.DefinedFunctions["i_base_string"] = Function{Exists:true, Data: `
FUNCTION i_base_string
	PULL base
	GRAB string
	
	VAR num
	VAR exp
	ADD exp 0 1
	
	VAR end 
	ADD end 0 #string
	SUB end end 1
	
	VAR i
	VAR __first
	VAR __toobig
	VAR __toosmall
	VAR __invalid
	ADD i 0 end
	LOOP
		VAR __condition
		IF __first
			SUB i i 1
		ELSE
			ADD __first 0 1
		END
		SGE __condition i 0
		IF __condition
			ADD __condition 0 0
		ELSE
			BREAK
		END
		
		PLACE string
			PUSH i
			GET exp*i
		SEQ __condition exp*i 45
		IF __condition 
			MUL num num -1
			BREAK
		END
		
		
		SGT __toobig exp*i 57
		SLT __toosmall exp*i 46
		ADD __invalid __toobig __toosmall
		IF __invalid
			ERROR 1
			PUSH 0
			RETURN
		END
		
		SLT __condition exp*i 48
		
		IF __condition
			ADD num 0 0
			ADD exp 0 1
		ELSE
		
			#Convert from unicode.
			SUB exp*i exp*i 48
			MUL exp*i exp exp*i
	
			ADD num num exp*i
	
			MUL exp exp base
		END
	REPEAT
	
	PUSH num
RETURN
`}

ic.DefinedFunctions["i_base_number"] = Function{Exists:true, Data: `
FUNCTION i_base_number
	PULL base
	PULL num
	ARRAY txt
	
	VAR test
	
	SEQ test num 0
	IF test
		PLACE txt
			PUT 48
		SHARE txt
		RETURN
	END
	
	VAR exp
	VAR exp>num
	VAR num<0
	
	ADD exp exp 1
	
	SLT num<0 num 0
	IF num<0
		PLACE txt
			PUT 45
		MUL num num -1
	END
	
	#What is the highest power to 10 which fits in num.
	LOOP
		SGT exp>num exp num
		IF exp>num 
			DIV exp exp base
			BREAK
		END
		
		MUL exp exp base
	REPEAT
	
	VAR num/exp
	VAR exp*(num/exp)
	VAR exp<=0
	
	#Find each digit.
	LOOP
		SLE exp<=0 exp 0
		IF exp<=0  
			BREAK
		END
		DIV num/exp num exp
		MUL exp*(num/exp) exp num/exp
		SUB num num exp*(num/exp)
		
		ADD num/exp num/exp 48
		PLACE txt
			PUT num/exp
		
		DIV exp exp base
	REPEAT
	SHARE txt
RETURN
`}

	ic.DefinedFunctions["copy_m_text"] = Alias("copy_m_numberlist", Text)
	
	ic.DefinedFunctions["text_m_numbers"] = SimpleMethod(Text, `
FUNCTION text_m_number
	PULL num
	ARRAY txt
	
	VAR test
	
	SEQ test num 0
	IF test
		PLACE txt
			PUT 48
		SHARE txt
		RETURN
	END
	
	VAR tens
	VAR tens>num
	VAR num<0
	
	ADD tens tens 1
	
	SLT num<0 num 0
	IF num<0
		PLACE txt
			PUT 45
		MUL num num -1
	END
	
	#What is the highest power to 10 which fits in num.
	LOOP
		SGT tens>num tens num
		IF tens>num 
			DIV tens tens 10
			BREAK
		END
		
		MUL tens tens 10
	REPEAT
	
	VAR num/tens 
	VAR tens*(num/tens)
	VAR tens<=0
	
	#Find each digit.
	LOOP
		SLE tens<=0 tens 0
		IF tens<=0  
			BREAK
		END
		DIV num/tens num tens
		MUL tens*(num/tens) tens num/tens
		SUB num num tens*(num/tens)
		
		ADD num/tens num/tens 48
		PLACE txt
			PUT num/tens
		
		DIV tens tens 10
	REPEAT
	SHARE txt
RETURN
`)
	
	//ic.DefinedFunctions["text_m_letter"] = Alias("text_m_number", Text)

	/*ic.DefinedFunctions["split_m_text"] = Function{Exists:true, Args:[]Type{GetType("letter"), Number}, Returns:[]Type{Text.MakeList()}, Data: `
FUNCTION split_m_text
	PULL amount
	PULL char
	GRAB string
	ARRAY i_array1
	SHARE i_array1
	GRAB result
	ARRAY i_string2
	SHARE i_string2
	GRAB segment
	PUSH 1
	PULL split
	
	IF 1
	ARRAY i_delete6
	VAR i_i4
	VAR i_backup5
	LOOP
		VAR i_in3
		ADD i_i4 0 i_backup5
		SGE i_in3 i_i4 #string
		IF i_in3
			BREAK
		END
		PLACE string
		PUSH i_i4
		GET c
		ADD i_backup5 i_i4 1
	
		VAR i_operator7
		SEQ i_operator7 c char
		VAR i_operator8
		MUL i_operator8 i_operator7 split
		IF i_operator8
			SHARE segment
			PUSH 0
			HEAP
			PULL i_index9
			PLACE result
			PUT i_index9
			ARRAY i_string10
			PLACE i_string10
			RENAME segment
			SUB amount amount 1
			VAR i_operator12
			DIV i_operator12 amount 0
			IF i_operator12
				ADD split 0 0
			END
		ELSE
			PLACE segment
			PUT c
		END
	REPEAT
	
		VAR ii_i8
		VAR ii_backup9
		LOOP
			VAR ii_in7
			ADD ii_i8 0 ii_backup9
			SGE ii_in7 ii_i8 #i_delete6
			IF ii_in7
				BREAK
			END
			PLACE i_delete6
			PUSH ii_i8
			GET i_v
			ADD ii_backup9 ii_i8 1
	
			VAR ii_operator11
			SUB ii_operator11 #string 1
			PLACE string
			PUSH ii_operator11
			GET ii_index12
			PLACE string
			PUSH i_v
			SET ii_index12
			PLACE string
			POP n
			ADD n 0 0
		REPEAT
							
	END
	VAR i_operator14
	SGT i_operator14 #segment 0
	IF i_operator14
		SHARE segment
		PUSH 0
		HEAP
		PULL i_index15
		PLACE result
		PUT i_index15
	END
	SHARE result
RETURN
	`}*/
	
	ic.DefinedFunctions["strings.equal"] = Function{Exists:true, Args:[]Type{Text}, Returns:[]Type{Text}, Data: `

FUNCTION collect_m_something
	GRAB something
	PUSH 1
	PLACE something
	GET garbage
	IF garbage
		PUSH 0
		GET address
		MUL address -1 address
		SUB garbage garbage 1
		IF garbage
			SUB garbage garbage 1
			IF garbage
				#Collect user type TODO
				ADD garbage 0 0
			ELSE
				PUSH address
				#HEAPIT
			END
		ELSE
			PUSH address
			HEAP
		END
	END
RETURN

FUNCTION i_unhash
	PULL exp
	PULL num
	ARRAY txt
	
	VAR test
	
	SEQ test num 0
	IF test
		SHARE txt
		RETURN
	END
	
	VAR tens
	VAR tens>num
	
	ADD tens tens 1
	
	#What is the highest power to 10 which fits in num.
	LOOP
		SGT tens>num tens num
		IF tens>num 
			DIV tens tens exp
			BREAK
		END
		
		MUL tens tens exp
	REPEAT
	
	VAR num/tens 
	VAR tens*(num/tens)
	VAR tens<=0
	
	#Find each digit.
	LOOP
		SLE tens<=0 tens 0
		IF tens<=0  
			BREAK
		END
		DIV num/tens num tens
		MUL tens*(num/tens) tens num/tens
		SUB num num tens*(num/tens)
	
		PLACE txt
			PUT num/tens
		
		DIV tens tens exp
	REPEAT
	SHARE txt
RETURN

FUNCTION i_hash
	GRAB text
	PULL exp
	
	VAR number
	VAR tens
	ADD tens 0 1
	
	VAR end 
	ADD end 0 #text
	SUB end end 1
	
	VAR i
	VAR __first
	ADD i 0 end
	LOOP
		VAR __condition
		IF __first
			SUB i i 1
		ELSE
			ADD __first 0 1
		END
		SGE __condition i 0
		IF __condition
			ADD __condition 0 0
		ELSE
			BREAK
		END
		
		PLACE text
			PUSH i
			GET tens*i
		
		MUL tens*i tens tens*i
		
		ADD number number tens*i
		
		MUL tens tens exp
	REPEAT
	
	PUSH number
RETURN
	
FUNCTION strings.equal
	GRAB str1
	GRAB str2
	
	VAR len(str1)!=len(str2)
	SNE len(str1)!=len(str2) #str1 #str2
	IF len(str1)!=len(str2)
		PUSH 0
		RETURN
	END
	
	VAR iterator
	VAR i>=len(str1)
	VAR char1!=char2
	LOOP
		SGE i>=len(str1) iterator #str1
		IF i>=len(str1)
			PUSH 1
			RETURN
		END
		
		PLACE str1
			PUSH iterator
			GET char1
			
		PLACE str2
			PUSH iterator
			GET char2
		
		SNE char1!=char2 char1 char2
		IF char1!=char2 
			PUSH 0
			RETURN
		END
		
		ADD iterator iterator 1
	REPEAT
RETURN
`}

ic.DefinedFunctions["sort_m_numberlist"] = Function{Exists:true, Args:[]Type{Array}, Data:`
FUNCTION sort_m_numberlist
	GRAB a
	PUSH #a
	MAKE
	GRAB b
	PUSH #a
	PULL num
	PUSH 0
	PULL rght
	PUSH 0
	PULL rend
	PUSH 0
	PULL i
	PUSH 0
	PULL j
	PUSH 0
	PULL m
	PUSH 1
	PULL k
	LOOP
		VAR i_operator2
		SLT i_operator2 k num
		VAR i_operator3
		DIV i_operator3 i_operator2 0
		IF i_operator3
			BREAK
		END
		PUSH 0
		PULL left
		LOOP
			VAR i_operator4
			ADD i_operator4 left k
			VAR i_operator5
			SLT i_operator5 i_operator4 num
			VAR i_operator6
			DIV i_operator6 i_operator5 0
			IF i_operator6
				BREAK
			END
			VAR i_operator7
			ADD i_operator7 left k
			ADD rght 0 i_operator7
			VAR i_operator8
			ADD i_operator8 rght k
			ADD rend 0 i_operator8
			VAR i_operator9
			SGT i_operator9 rend num
			IF i_operator9
				ADD rend 0 num
			END
			ADD m 0 left
			ADD i 0 left
			ADD j 0 rght
			LOOP
				VAR i_operator10
				SLT i_operator10 i rght
				VAR i_operator12
				SLT i_operator12 j rend
				VAR i_operator11
				MUL i_operator11 i_operator10 i_operator12
				VAR i_operator13
				DIV i_operator13 i_operator11 0
				IF i_operator13
					BREAK
				END
				PLACE a
				PUSH i
				GET i_index14
				PLACE a
				PUSH j
				GET i_index16
				VAR i_operator15
				SLE i_operator15 i_index14 i_index16
				IF i_operator15
					PLACE a
					PUSH i
					GET i_index17
					PLACE b
					PUSH m
					SET i_index17
					ADD i i 1
				ELSE
					PLACE a
					PUSH j
					GET i_index19
					PLACE b
					PUSH m
					SET i_index19
					ADD j j 1
				END
				ADD m m 1
			REPEAT
			LOOP
				VAR i_operator22
				SLT i_operator22 i rght
				VAR i_operator23
				DIV i_operator23 i_operator22 0
				IF i_operator23
					BREAK
				END
				PLACE a
				PUSH i
				GET i_index24
				PLACE b
				PUSH m
				SET i_index24
				ADD i i 1
				ADD m m 1
			REPEAT
			LOOP
				VAR i_operator27
				SLT i_operator27 j rend
				VAR i_operator28
				DIV i_operator28 i_operator27 0
				IF i_operator28
					BREAK
				END
				PLACE a
				PUSH j
				GET i_index29
				PLACE b
				PUSH m
				SET i_index29
				ADD j j 1
				ADD m m 1
			REPEAT
			ADD m 0 left
			LOOP
				VAR i_operator32
				SLT i_operator32 m rend
				VAR i_operator33
				DIV i_operator33 i_operator32 0
				IF i_operator33
					BREAK
				END
				PLACE b
				PUSH m
				GET i_index34
				PLACE a
				PUSH m
				SET i_index34
				ADD m m 1
			REPEAT
			VAR i_operator37
			MUL i_operator37 k 2
			ADD left left i_operator37
		REPEAT
		MUL k k 2
	REPEAT
RETURN
`}


ic.DefinedFunctions["watch"] = Function{Exists:true, Args:[]Type{Text}, Data:`
FUNCTION watch
	GRAB id
	
	ARRAY i+tmp+1
	SHARE i+tmp+1
	OPEN
	TAKE i+output+2
	RELAY i+output+2
	TAKE grabserver
	
	RELOAD grabserver
	RELAY grabserver
	ARRAY i+tmp+3
		PUT 87
		PUT 65
		PUT 84
		PUT 67
		PUT 72
		PUT 32
	ARRAY i+tmp+6
	PUT 10
	ARRAY i+shunt+5
	JOIN i+shunt+5 id i+tmp+6
	ARRAY i+shunt+4
	JOIN i+shunt+4 i+tmp+3 i+shunt+5
	SHARE i+shunt+4
	OUT
	RELAY grabserver
RETURN`}


	ic.DefinedFunctions["grab"] = Function{Exists:true, Args:[]Type{Text}, Returns:[]Type{Text}, Data:`
FUNCTION grab
	GRAB id
	
	ARRAY i+tmp+1
	SHARE i+tmp+1
	OPEN
	TAKE i+output+2
	RELAY i+output+2
	TAKE grabserver
	
	RELOAD grabserver
	ERROR 0
	RELAY grabserver
	ARRAY i+tmp+9
		PUT 71
		PUT 82
		PUT 65
		PUT 66
		PUT 32
	ARRAY i+tmp+12
	PUT 10
	ARRAY i+shunt+11
	JOIN i+shunt+11 id i+tmp+12
	ARRAY i+shunt+10
	JOIN i+shunt+10 i+tmp+9 i+shunt+11
	SHARE i+shunt+10
	OUT
	RELAY grabserver
	RELAY grabserver
	PUSH 1
	RUN reada_m_pipe
	RELAY grabserver
	GRAB i+output+13
	PUSH -1
	RUN reada_m_pipe
	GRAB a
	SHARE i+output+13
RETURN
`}

	ic.DefinedFunctions["gui"] = Function{Exists:true, Args:[]Type{Text}, Data:`
FUNCTION gui
	GRAB design
	ARRAY i+tmp+14
	PUT 116
	PUT 99
	PUT 112
	PUT 58
	PUT 47
	PUT 47
	PUT 108
	PUT 111
	PUT 99
	PUT 97
	PUT 108
	PUT 104
	PUT 111
	PUT 115
	PUT 116
	PUT 58
	PUT 50
	PUT 50
	PUT 50
	PUT 50
	SHARE i+tmp+14
	OPEN
	TAKE i+output+15
	RELAY i+output+15
	TAKE server
	IF ERROR
		
		ARRAY i+tmp+16
			PUT 67
			PUT 111
			PUT 117
			PUT 108
			PUT 100
			PUT 32
			PUT 110
			PUT 111
			PUT 116
			PUT 32
			PUT 99
			PUT 114
			PUT 101
			PUT 97
			PUT 116
			PUT 101
			PUT 32
			PUT 103
			PUT 117
			PUT 105
			PUT 33
		SHARE i+tmp+16
		STDOUT
		SHARE i_newline
		STDOUT
		EXIT
	END
	RELAY server
	ARRAY i+tmp+17
		PUT 68
		PUT 65
		PUT 84
		PUT 65
		PUT 32
	ARRAY i+tmp+20
	PUT 10
	ARRAY i+shunt+19
	JOIN i+shunt+19 design i+tmp+20
	ARRAY i+shunt+18
	JOIN i+shunt+18 i+tmp+17 i+shunt+19
	SHARE i+shunt+18
	OUT
	RELAY server
RETURN
`}

	ic.DefinedFunctions["edit"] = Function{Exists:true, Args:[]Type{Text, Text}, Data:`
FUNCTION edit
	GRAB txt
	GRAB id
	ARRAY i+tmp+1
	SHARE i+tmp+1
	OPEN
	TAKE i+output+2
	RELAY i+output+2
	TAKE grabserver
	RELOAD grabserver
	RELAY grabserver
	ARRAY i+tmp+3
		PUT 69
		PUT 68
		PUT 73
		PUT 84
		PUT 32
	ARRAY i+tmp+6
		PUT 32
	ARRAY i+tmp+9
		PUT 10
		PUT 1
		PUT 10
	ARRAY i+shunt+8
	JOIN i+shunt+8 txt i+tmp+9
	ARRAY i+shunt+7
	JOIN i+shunt+7 i+tmp+6 i+shunt+8
	ARRAY i+shunt+5
	JOIN i+shunt+5 id i+shunt+7
	ARRAY i+shunt+4
	JOIN i+shunt+4 i+tmp+3 i+shunt+5
	SHARE i+shunt+4
	OUT
	RELAY grabserver
RETURN
`}

ic.DefinedFunctions["trim_m_text"] = Function{Exists:true, Args:[]Type{Text}, Returns:[]Type{Text}, Data: `
FUNCTION trim_m_text
	GRAB s
	ARRAY i_string1
	SHARE i_string1
	GRAB result
	PUSH 0
	PULL done
	
	IF 1
	ARRAY i_delete5
	VAR i_i3
	VAR i_backup4
	LOOP
		VAR i_in2
		ADD i_i3 0 i_backup4
		SGE i_in2 i_i3 #s
		IF i_in2
			BREAK
		END
		PLACE s
		PUSH i_i3
		GET char
		ADD i_backup4 i_i3 1
	
		VAR i_operator6
		SNE i_operator6 char 32
		VAR i_operator8
		SNE i_operator8 char 9
		VAR i_operator7
		MUL i_operator7 i_operator6 i_operator8
		VAR i_operator10
		SNE i_operator10 char 10
		VAR i_operator9
		MUL i_operator9 i_operator7 i_operator10
		VAR i_operator12
		SNE i_operator12 char 13
		VAR i_operator11
		MUL i_operator11 i_operator9 i_operator12
		VAR i_operator13
		ADD i_operator13 i_operator11 done
		IF i_operator13
			ADD done 0 1
			PLACE result
			PUT char
		END
	REPEAT
	
		VAR ii_i8
		VAR ii_backup9
		LOOP
			VAR ii_in7
			ADD ii_i8 0 ii_backup9
			SGE ii_in7 ii_i8 #i_delete5
			IF ii_in7
				BREAK
			END
			PLACE i_delete5
			PUSH ii_i8
			GET i_v
			ADD ii_backup9 ii_i8 1
	
			VAR ii_operator11
			SUB ii_operator11 #s 1
			PLACE s
			PUSH ii_operator11
			GET ii_index12
			PLACE s
			PUSH i_v
			SET ii_index12
			PLACE s
			POP n
			ADD n 0 0
		REPEAT
							
	END
	LOOP
		VAR i_operator15
		SUB i_operator15 0 1
		PLACE result
		PUSH i_operator15
		GET i_index16
		VAR i_operator17
		SEQ i_operator17 i_index16 32
		VAR i_operator19
		SUB i_operator19 0 1
		PLACE result
		PUSH i_operator19
		GET i_index20
		VAR i_operator21
		SEQ i_operator21 i_index20 9
		VAR i_operator23
		SUB i_operator23 0 1
		PLACE result
		PUSH i_operator23
		GET i_index24
		VAR i_operator25
		SEQ i_operator25 i_index24 13
		VAR i_operator27
		SUB i_operator27 0 1
		PLACE result
		PUSH i_operator27
		GET i_index28
		VAR i_operator29
		SEQ i_operator29 i_index28 10
		VAR i_operator26
		ADD i_operator26 i_operator25 i_operator29
		VAR i_operator22
		ADD i_operator22 i_operator21 i_operator26
		VAR i_operator18
		ADD i_operator18 i_operator17 i_operator22
		IF i_operator18
			PLACE result
			POP i_tmp31
		ELSE
			BREAK
		END
	REPEAT
	SHARE result
RETURN
`}

ic.DefinedFunctions["random_m_number"] = Function{Exists:true, Args:[]Type{Number}, Returns:[]Type{Number}, Data:`
FUNCTION random_m_number
	PULL limit
	PUSH limit
	PUSH 255
	RUN i_unhash
	GRAB i_operator1
	
	SHARE i_operator1
	GRAB bits
	IF 1
	VAR i
	VAR i_backup3
	ADD i_backup3 0 0
	ADD i 0 0
	LOOP
		VAR i_over2
		SNE i_over2 i #bits
		ADD i 0 i_backup3
		IF i_over2
			SLT i_over2 i #bits
			IF i_over2
				ADD i_backup3 i 1
			ELSE
				SUB i_backup3 i 1
			END
			SEQ i_over2 0 #bits
			IF i_over2
				BREAK
			END
		ELSE
			SEQ i_over2 0 #bits
			IF i_over2
				ADD i i 1
	       ELSE
				BREAK
			END
		END
	
		VAR i_operator4
		DIV i_operator4 0 0
		PLACE bits
		PUSH i
		SET i_operator4
	REPEAT
	END
	SHARE bits
	PUSH 255
	RUN i_hash
	PULL i_operator5
	
	VAR i_operator6
	MOD i_operator6 i_operator5 limit
	PUSH i_operator6
RETURN
`}

ic.DefinedFunctions["text_m_decimal"] = Function{Exists:true, Args:[]Type{Number}, Import:"i_base_number", Returns:[]Type{Text}, Data:`
FUNCTION text_m_decimal
	PULL value
	VAR i_operator2
	MOD i_operator2 value 1000000
	
	VAR test
	VAR test2
	SLT test value 0
	IF test
		SUB i_operator2 1000000 i_operator2 
	END
	
	SGT test2 value -1000000
	IF test
	IF test2
		ADD value 0 0
	END
	END
	
	PUSH i_operator2
	PUSH 10
	RUN i_base_number
	GRAB i_result3
	SHARE i_result3
	GRAB decimal
	VAR i_operator4
	SNE i_operator4 #decimal 6
	IF i_operator4
		IF 1
		VAR each
		VAR i_backup6
		ADD i_backup6 0 #decimal
		ADD each 0 #decimal
		LOOP
			VAR i_over5
			SNE i_over5 each 6
			ADD each 0 i_backup6
			IF i_over5
				SLT i_over5 each 6
				IF i_over5
					ADD i_backup6 each 1
				ELSE
					SUB i_backup6 each 1
				END
				SEQ i_over5 #decimal 6
				IF i_over5
					BREAK
				END
			ELSE
				SEQ i_over5 #decimal 6
				IF i_over5
					ADD each each 1
		       ELSE
					BREAK
				END
			END
		
			ARRAY i_string7
			PUT 48
			ARRAY i_operator8
			JOIN i_operator8 i_string7 decimal
			PLACE i_operator8
			RENAME decimal
		REPEAT
		END
	END
	VAR i_operator9
	SUB i_operator9 #decimal 1
	IF 1
	VAR i
	VAR i_backup11
	ADD i_backup11 0 i_operator9
	ADD i 0 i_operator9
	LOOP
		VAR i_over10
		SNE i_over10 i 0
		ADD i 0 i_backup11
		IF i_over10
			SLT i_over10 i 0
			IF i_over10
				ADD i_backup11 i 1
			ELSE
				SUB i_backup11 i 1
			END
			SEQ i_over10 i_operator9 0
			IF i_over10
				BREAK
			END
		ELSE
			SEQ i_over10 i_operator9 0
			IF i_over10
				ADD i i 1
	       ELSE
				BREAK
			END
		END
	
		PLACE decimal
		PUSH i
		GET i_index12
		PUSH 48
		
		PULL i_result14
		VAR i_operator13
		SEQ i_operator13 i_index12 i_result14
		IF i_operator13
			PLACE decimal
			POP i_tmp16
			ADD i_tmp16 0 0
		ELSE
			BREAK
		END
	REPEAT
	END
	VAR i_operator17
	SGT i_operator17 #decimal 0
	IF i_operator17
		ARRAY i_string18
		PUT 46
		ARRAY i_operator19
		JOIN i_operator19 i_string18 decimal
		PLACE i_operator19
		RENAME decimal
	END
	VAR i_operator20
	DIV i_operator20 value 1000000
	PUSH i_operator20
	PUSH 10
	RUN i_base_number
	GRAB i_result21
	IF test
	IF test2
		POP n
		ADD n 0 0
		PUT 45
		PUT 48
	END
	END
	ARRAY i_operator22
	JOIN i_operator22 i_result21 decimal
	SHARE i_operator22
RETURN
`}

ic.DefinedFunctions["replace_m_text"] = Function{Exists:true, Args:[]Type{Text, Text}, Returns:[]Type{Text}, Data:`
FUNCTION replace_m_text
	GRAB b
	GRAB a
	GRAB s
	ARRAY i_string1
	SHARE i_string1
	GRAB result
	PUSH 0
	PULL found
	PUSH 0
	PULL skip
	
	IF 1
	ARRAY i_delete4
	VAR i
	VAR i_backup3
	LOOP
		VAR i_in2
		ADD i 0 i_backup3
		SGE i_in2 i #s
		IF i_in2
			BREAK
		END
		PLACE s
		PUSH i
		GET char
		ADD i_backup3 i 1
	
		IF skip
			SUB skip skip 1
		ELSE
			PLACE a
			PUSH 0
			GET i_index7
			VAR i_operator6
			SEQ i_operator6 char i_index7
			IF i_operator6
				ADD found 0 1
				
				IF 1
				ARRAY i_delete10
				VAR j
				VAR i_backup9
				LOOP
					VAR i_in8
					ADD j 0 i_backup9
					SGE i_in8 j #a
					IF i_in8
						BREAK
					END
					PLACE a
					PUSH j
					GET lettr
					ADD i_backup9 j 1
				
					VAR i_operator11
					ADD i_operator11 i j
					PLACE s
					PUSH i_operator11
					GET i_index12
					VAR i_operator13
					SNE i_operator13 i_index12 lettr
					IF i_operator13
						ADD found 0 0
					END
				REPEAT
				
					VAR ii_i8
					VAR ii_backup9
					LOOP
						VAR ii_in7
						ADD ii_i8 0 ii_backup9
						SGE ii_in7 ii_i8 #i_delete10
						IF ii_in7
							BREAK
						END
						PLACE i_delete10
						PUSH ii_i8
						GET i_v
						ADD ii_backup9 ii_i8 1
				
						VAR ii_operator11
						SUB ii_operator11 #a 1
						PLACE a
						PUSH ii_operator11
						GET ii_index12
						PLACE a
						PUSH i_v
						SET ii_index12
						PLACE a
						POP n
						ADD n 0 0
					REPEAT
										
				END
				IF found
					VAR i_operator14
					SUB i_operator14 #a 1
					ADD skip 0 i_operator14
					JOIN result result b
				ELSE
					PLACE result
					PUT char
				END
			ELSE
				PLACE result
				PUT char
			END
		END
	REPEAT
	
	IF 1
		VAR ii_i8
		VAR ii_backup9
		LOOP
			VAR ii_in7
			ADD ii_i8 0 ii_backup9
			SGE ii_in7 ii_i8 #i_delete4
			IF ii_in7
				BREAK
			END
			PLACE i_delete4
			PUSH ii_i8
			GET i_v
			ADD ii_backup9 ii_i8 1
	
			VAR ii_operator11
			SUB ii_operator11 #s 1
			PLACE s
			PUSH ii_operator11
			GET ii_index12
			PLACE s
			PUSH i_v
			SET ii_index12
			PLACE s
			POP n
			ADD n 0 0
		REPEAT
	END			
	END
	SHARE result
RETURN
`}

//TODO fix this function, it needs to ignore the distinction between upper and lowercase letters.
ic.DefinedFunctions["strings.compare"] = Function{Exists:true, Args:[]Type{Text, Text}, Returns:[]Type{Number}, Data:`
FUNCTION strings.compare
	GRAB b
	GRAB a
	PUSH 0
	PULL cmp
	
	IF 1
	ARRAY i_delete3
	VAR i
	VAR i_backup2
	LOOP
		VAR i_in1
		ADD i 0 i_backup2
		SGE i_in1 i #a
		IF i_in1
			BREAK
		END
		PLACE a
		PUSH i
		GET char
		ADD i_backup2 i 1
	
		VAR i_operator4
		SLT i_operator4 #b i
		IF i_operator4
			ADD cmp 0 1
			BREAK
		END
		PLACE b
		PUSH i
		GET i_index5
		PUSH i_index5
		
		PULL i_result6
		PUSH char
		
		PULL i_result8
		VAR i_operator7
		SNE i_operator7 i_result6 i_result8
		IF i_operator7
			PLACE b
			PUSH i
			GET i_index9
			PUSH i_index9
			
			PULL i_result10
			PUSH char
			
			PULL i_result12
			VAR i_operator11
			SLT i_operator11 i_result10 i_result12
			IF i_operator11
				ADD cmp 0 1
			ELSE
				VAR i_operator13
				SUB i_operator13 0 1
				ADD cmp 0 i_operator13
			END
			BREAK
		END
	REPEAT
	
		VAR ii_i8
		VAR ii_backup9
		LOOP
			VAR ii_in7
			ADD ii_i8 0 ii_backup9
			SGE ii_in7 ii_i8 #i_delete3
			IF ii_in7
				BREAK
			END
			PLACE i_delete3
			PUSH ii_i8
			GET i_v
			ADD ii_backup9 ii_i8 1
	
			VAR ii_operator11
			SUB ii_operator11 #a 1
			PLACE a
			PUSH ii_operator11
			GET ii_index12
			PLACE a
			PUSH i_v
			SET ii_index12
			PLACE a
			POP n
			ADD n 0 0
		REPEAT
							
	END
	PUSH cmp
RETURN
`}

/*ic.DefinedFunctions["sort_m_textarray"] = Function{Exists:true, List:true, Args:[]Type{Text.MakeList()}, Data:`
FUNCTION sort_m_textarray
	GRAB a
	PUSH #a
	MAKE
	GRAB b
	PUSH #a
	PULL num
	PUSH 0
	PULL rght
	PUSH 0
	PULL rend
	PUSH 0
	PULL i
	PUSH 0
	PULL j
	PUSH 0
	PULL m
	PUSH 1
	PULL k
	LOOP
		VAR i_operator7
		SLT i_operator7 k num
		VAR i_operator8
		DIV i_operator8 i_operator7 0
		IF i_operator8
			BREAK
		END
		PUSH 0
		PULL left
		LOOP
			VAR i_operator9
			ADD i_operator9 left k
			VAR i_operator10
			SLT i_operator10 i_operator9 num
			VAR i_operator11
			DIV i_operator11 i_operator10 0
			IF i_operator11
				BREAK
			END
			VAR i_operator12
			ADD i_operator12 left k
			ADD rght 0 i_operator12
			VAR i_operator13
			ADD i_operator13 rght k
			ADD rend 0 i_operator13
			VAR i_operator14
			SGT i_operator14 rend num
			IF i_operator14
				ADD rend 0 num
			END
			ADD m 0 left
			ADD i 0 left
			ADD j 0 rght
			LOOP
				VAR i_operator15
				SLT i_operator15 i rght
				VAR i_operator17
				SLT i_operator17 j rend
				VAR i_operator16
				MUL i_operator16 i_operator15 i_operator17
				VAR i_operator18
				DIV i_operator18 i_operator16 0
				IF i_operator18
					BREAK
				END
				PLACE a
				PUSH i
				GET i_index20
				IF i_index20
				PUSH i_index20
				HEAP
				ELSE
				ARRAY i_listdex19
				SHARE i_listdex19
				END
				GRAB i_listdex19
				PLACE a
				PUSH j
				GET i_index23
				IF i_index23
				PUSH i_index23
				HEAP
				ELSE
				ARRAY i_listdex22
				SHARE i_listdex22
				END
				GRAB i_listdex22
				SHARE i_listdex19
				SHARE i_listdex22
				RUN strings.compare
				PULL i_operator21
				SLE i_operator21 i_operator21 0
				IF i_operator21
					PLACE a
					PUSH i
					GET i_index25
					IF i_index25
					PUSH i_index25
					HEAP
					ELSE
					ARRAY i_listdex24
					SHARE i_listdex24
					END
					GRAB i_listdex24
					PLACE b
					PUSH m
					PLACE b
					PUSH m
					GET i_index26
					IF i_index26
					PUSH i_index26
					HEAP
					
					MUL i_index26 i_index26 -1
					PUSH i_index26
					HEAP
					END
					SHARE i_listdex24
					PUSH 0
					HEAP
					PULL i_index27
					PLACE b
					PUSH m
					SET i_index27
					ADD i i 1
				ELSE
					PLACE a
					PUSH j
					GET i_index30
					IF i_index30
					PUSH i_index30
					HEAP
					ELSE
					ARRAY i_listdex29
					SHARE i_listdex29
					END
					GRAB i_listdex29
					PLACE b
					PUSH m
					PLACE b
					PUSH m
					GET i_index31
					IF i_index31
					PUSH i_index31
					HEAP
					
					MUL i_index31 i_index31 -1
					PUSH i_index31
					HEAP
					END
					SHARE i_listdex29
					PUSH 0
					HEAP
					PULL i_index32
					PLACE b
					PUSH m
					SET i_index32
					ADD j j 1
				END
				ADD m m 1
			REPEAT
			LOOP
				VAR i_operator35
				SLT i_operator35 i rght
				VAR i_operator36
				DIV i_operator36 i_operator35 0
				IF i_operator36
					BREAK
				END
				PLACE a
				PUSH i
				GET i_index38
				IF i_index38
				PUSH i_index38
				HEAP
				ELSE
				ARRAY i_listdex37
				SHARE i_listdex37
				END
				GRAB i_listdex37
				PLACE b
				PUSH m
				PLACE b
				PUSH m
				GET i_index39
				IF i_index39
				PUSH i_index39
				HEAP
				
				MUL i_index39 i_index39 -1
				PUSH i_index39
				HEAP
				END
				SHARE i_listdex37
				PUSH 0
				HEAP
				PULL i_index40
				PLACE b
				PUSH m
				SET i_index40
				ADD i i 1
				ADD m m 1
			REPEAT
			LOOP
				VAR i_operator43
				SLT i_operator43 j rend
				VAR i_operator44
				DIV i_operator44 i_operator43 0
				IF i_operator44
					BREAK
				END
				PLACE a
				PUSH j
				GET i_index46
				IF i_index46
				PUSH i_index46
				HEAP
				ELSE
				ARRAY i_listdex45
				SHARE i_listdex45
				END
				GRAB i_listdex45
				PLACE b
				PUSH m
				PLACE b
				PUSH m
				GET i_index47
				IF i_index47
				PUSH i_index47
				HEAP
				
				MUL i_index47 i_index47 -1
				PUSH i_index47
				HEAP
				END
				SHARE i_listdex45
				PUSH 0
				HEAP
				PULL i_index48
				PLACE b
				PUSH m
				SET i_index48
				ADD j j 1
				ADD m m 1
			REPEAT
			ADD m 0 left
			LOOP
				VAR i_operator51
				SLT i_operator51 m rend
				VAR i_operator52
				DIV i_operator52 i_operator51 0
				IF i_operator52
					BREAK
				END
				PLACE b
				PUSH m
				GET i_index54
				IF i_index54
				PUSH i_index54
				HEAP
				ELSE
				ARRAY i_listdex53
				SHARE i_listdex53
				END
				GRAB i_listdex53
				PLACE a
				PUSH m
				PLACE a
				PUSH m
				GET i_index55
				IF i_index55
				PUSH i_index55
				HEAP
				
				MUL i_index55 i_index55 -1
				PUSH i_index55
				HEAP
				END
				SHARE i_listdex53
				PUSH 0
				HEAP
				PULL i_index56
				PLACE a
				PUSH m
				SET i_index56
				ADD m m 1
			REPEAT
			VAR i_operator59
			MUL i_operator59 k 2
			ADD left left i_operator59
		REPEAT
		MUL k k 2
	REPEAT
RETURN
`}*/

}
