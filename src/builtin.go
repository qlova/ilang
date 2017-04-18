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
	ic.DefinedFunctions["letter"] = Method(Letter, true, "PUSH 0")
	ic.DefinedFunctions["binary"] = Method(Number, true, "PUSH 0")
	ic.DefinedFunctions["binary"] = Method(Number, true, "PUSH 0")
	
	ic.DefinedFunctions["random"] = Method(Number, true, "PUSH 0")
	
	ic.DefinedFunctions["load"] = Method(Text, true, "PUSH 0")
	ic.DefinedFunctions["text"] = Method(Number, false, `
FUNCTION text
	ARRAY a
	SHARE a
RETURN
`)
	
	ic.DefinedFunctions["random_m_decimal"] = Alias("random_m_number", Decimal)

	ic.DefinedFunctions["copy"] = Method(Undefined, true, "")
	
	ic.DefinedFunctions["number_m_letter"] = BlankMethod(Number)
	ic.DefinedFunctions["letter_m_number"] = BlankMethod(Letter)
	ic.DefinedFunctions["text_m_text"] = BlankMethod(Text)
	ic.DefinedFunctions["text_m_array"] = BlankMethod(Text)
	
	ic.DefinedFunctions["load"] = Method(Undefined, true, "")
	ic.DefinedFunctions["open"] = Method(Undefined, true, "")
	
	ic.DefinedFunctions["open_m_text"] = Method(Pipe, true, "OPEN")
	ic.DefinedFunctions["execute"] = InlineFunction([]Type{Text}, "EXECUTE", nil)
	ic.DefinedFunctions["delete"] = InlineFunction([]Type{Text}, "DELETE", nil)
	ic.DefinedFunctions["rename"] = InlineFunction([]Type{Text, Text}, "MOVE", nil)
	
	ic.DefinedFunctions["read"] = Method(Text, true, "PUSH 0\nSTDIN")
	ic.DefinedFunctions["read_m_pipe"] = InlineFunction(nil, "PUSH 0\nIN", []Type{Text})
	
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
	ic.DefinedFunctions["len_m_array"] = InlineFunction([]Type{Array}, "LEN", nil)
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
	
	ic.DefinedFunctions["read_m_letter"] = Function{Exists:true, Args:[]Type{Letter}, Returns:[]Type{Text}, Data: `
FUNCTION read_m_letter
	PULL delim
	MUL delim delim -1
	PUSH delim
	STDIN
RETURN
`}

ic.DefinedFunctions["reada_m_pipe"] = Function{Exists:true, Args:[]Type{Letter}, Returns:[]Type{Text}, Data:`
FUNCTION reada_m_pipe
	PULL delim
	MUL delim delim -1
	PUSH delim
	IN
END
`}


ic.DefinedFunctions["close"] = Function{Exists:true, Args:[]Type{Pipe}, Data: `
FUNCTION close
	TAKE file
	CLOSE file
RETURN
`}

ic.DefinedFunctions["text_m_Something"] = Function{Exists:true, Returns:[]Type{Text}, Data: `
FUNCTION text_m_Something
	GRAB something
	PUSH #something
	PUSH 3
	SHARE something
	SLICE
RETURN
`}

	
	
	ic.DefinedFunctions["number_m_array"] = Alias("number_m_text", Text)
	
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
END
`}
	ic.DefinedFunctions["copy_m_array"] = Function{Exists:true, Args:[]Type{Array}, Returns:[]Type{Array}, Data: `
#Compiled with IC.
FUNCTION copy_m_array
	GRAB array
	ARRAY c
	
	VAR i
	LOOP
		VAR i+shunt+1
		SGE i+shunt+1 i #array
		IF i+shunt+1
			SHARE c
			RETURN
		END
		PLACE array 
			PUSH i 
			GET i+shunt+3
		VAR v
		ADD v 0 i+shunt+3
		PLACE c
			PUT v
		ADD i i 1
	REPEAT
END
`}

	ic.DefinedFunctions["copy_m_text"] = Alias("copy_m_array", Text)
	
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
	
	ic.DefinedFunctions["strings.equal"] = Function{Exists:true, Args:[]Type{Text}, Returns:[]Type{Text}, Data: `

FUNCTION collect_m_Something
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

FUNCTION unhash
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

FUNCTION hash
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

ic.DefinedFunctions["reada_m_text"] = Function{Exists:true, Args:[]Type{Letter}, Returns:[]Type{Text}, Data:`
FUNCTION reada_m_text
	GRAB s
	PULL n

	ARRAY i+string+1
	SHARE i+string+1
	GRAB result
	VAR i
	ADD i 0 0
	LOOP
		VAR i+shunt+2
		SLT i+shunt+2 i #s
		IF i+shunt+2
			ERROR 0
		ELSE
			BREAK
		END
		PLACE s
		PUSH i
		GET i+shunt+4
		VAR i+shunt+5
		SEQ i+shunt+5 i+shunt+4 n
		IF i+shunt+5
			VAR j
			ADD j 0 1
			LOOP
				VAR i+shunt+6
				SLT i+shunt+6 j #s
				IF i+shunt+6
					ERROR 0
				ELSE
					BREAK
				END
				VAR i+shunt+8
				SUB i+shunt+8 j 1
				VAR i+shunt+9
				ADD i+shunt+9 j i
				PLACE s
				PUSH i+shunt+9
				GET i+shunt+10
				PLACE s
				PUSH i+shunt+8
				SET i+shunt+10
				VAR i+shunt+11
				ADD i+shunt+11 j 1
				ADD j 0 i+shunt+11
			REPEAT
			ADD j 0 0
			LOOP
				VAR i+shunt+12
				SLE i+shunt+12 j i
				IF i+shunt+12
					ERROR 0
				ELSE
					BREAK
				END
				PLACE s
				POP z
				MUL z z 0
				VAR i+shunt+13
				ADD i+shunt+13 j 1
				ADD j 0 i+shunt+13
			REPEAT
			SHARE result
			RETURN
		END
		PLACE s
		PUSH i
		GET i+shunt+14
		PLACE result
		PUT i+shunt+14
		VAR i+shunt+15
		ADD i+shunt+15 i 1
		ADD i 0 i+shunt+15
	REPEAT
	ERROR 1
	SHARE result
RETURN
`}

ic.DefinedFunctions["sort"] = Function{Exists:true, Args:[]Type{Array}, Data:`
FUNCTION i_part
PULL back
PULL start
GRAB alist
PLACE alist
PUSH back
GET i+shunt+1
VAR pivot
ADD pivot 0 i+shunt+1
VAR border
ADD border 0 start
VAR i+shunt+2
SLT i+shunt+2 start back
IF i+shunt+2
VAR i
ADD i 0 start
LOOP
VAR i+shunt+4
ADD i+shunt+4 back 1
VAR i+shunt+3
SNE i+shunt+3 i i+shunt+4
IF i+shunt+3
ERROR 0
ELSE
BREAK
END
PLACE alist
PUSH i
GET i+shunt+5
VAR i+shunt+6
SLE i+shunt+6 i+shunt+5 pivot
IF i+shunt+6
PLACE alist
PUSH i
GET i+shunt+7
VAR ab
ADD ab 0 i+shunt+7
PLACE alist
PUSH border
GET i+shunt+8
VAR ai
ADD ai 0 i+shunt+8
PLACE alist
PUSH i
SET ai
PLACE alist
PUSH border
SET ab
VAR i+shunt+9
SNE i+shunt+9 i back
IF i+shunt+9
VAR i+shunt+10
ADD i+shunt+10 border 1
ADD border 0 i+shunt+10
END
END
VAR i+shunt+12
ADD i+shunt+12 back 1
VAR i+shunt+11
SLT i+shunt+11 i i+shunt+12
IF i+shunt+11
VAR i+shunt+13
ADD i+shunt+13 i 1
ADD i 0 i+shunt+13
ELSE
VAR i+shunt+15
ADD i+shunt+15 back 1
VAR i+shunt+14
SGT i+shunt+14 i i+shunt+15
IF i+shunt+14
VAR i+shunt+16
SUB i+shunt+16 i 1
ADD i 0 i+shunt+16
END
END
REPEAT
SHARE alist
PUSH start
VAR i+shunt+17
SUB i+shunt+17 border 1
PUSH i+shunt+17
RUN i_part
SHARE alist
VAR i+shunt+18
ADD i+shunt+18 border 1
PUSH i+shunt+18
PUSH back
RUN i_part
END
RETURN
FUNCTION sort
GRAB alist
VAR i+shunt+20
SLE i+shunt+20 #alist 1
IF i+shunt+20
RETURN
END
SHARE alist
PUSH 0
VAR i+shunt+22
SUB i+shunt+22 #alist 1
PUSH i+shunt+22
RUN i_part
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
	
	ic.DefinedFunctions["print_m_array"] = Function{Exists:true, Args:[]Type{Array}, Data:`
FUNCTION print_m_array
	GRAB a
	
	IF 1
	ARRAY i_delete4
	VAR i_i2
	VAR i_backup3
	LOOP
		VAR i_in1
		ADD i_i2 0 i_backup3
		SGE i_in1 i_i2 #a
		IF i_in1
			BREAK
		END
		PLACE a
		PUSH i_i2
		GET i
		ADD i_backup3 i_i2 1
	
		PUSH i
		PUSH 10
		RUN i_base_number
		GRAB i_result5
		ARRAY i_string7
		PUT 32
		ARRAY i_operator6
		JOIN i_operator6 i_result5 i_string7
		SHARE i_operator6
		STDOUT
	REPEAT
	
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

ic.DefinedFunctions["random_m_number"] = Function{Exists:true, Args:[]Type{Number}, Returns:[]Type{Number}, Data:`
FUNCTION random_m_number
	PULL limit
	PUSH limit
	PUSH 255
	RUN unhash
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
	RUN hash
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
	SHARE result
RETURN
`}
}
