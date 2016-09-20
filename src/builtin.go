package main

import "io"

//TODO optimise some of the functions to be inline.
func builtin(output io.Writer) {
	
	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION output
	STDOUT 
END
`	))
	functions["output"] = Function{Exists:true, Args:[]TYPE{STRING}}
	
	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION execute
	EXECUTE 
END
`	))
	functions["execute"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{STRING}}
	
	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION link
	LINK 
END
`	))
	functions["link"] = Function{Exists:true, Args:[]TYPE{STRING, NUMBER}}
	
	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION delete
	DELETE
END
`	))
	functions["delete"] = Function{Exists:true, Args:[]TYPE{STRING}}
	
	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION connect
	CONNECT 
END
`	))
	functions["connect"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{NUMBER}}

	output.Write([]byte(
`
#Returns whether or not a string is equal.
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
` ))
	functions["strings.equal"] = Function{Exists:true, Args:[]TYPE{STRING, STRING}, Returns:[]TYPE{NUMBER}}

	output.Write([]byte(
`
DATA i_true "true"
DATA i_false "false"
FUNCTION bool
	PULL n
	IF n
		SHARE i_true
		RUN copy
		RETURN
	END
	SHARE i_false
	RUN copy
END
`	))
	functions["bool"] = Function{Exists:true, Args:[]TYPE{NUMBER}, Returns:[]TYPE{STRING}}

	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION load
	LOAD
END
`	))
	functions["load"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{STRING}}
	
	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION open
	OPEN
END
`	))
	functions["open"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{FILE}}
	
	output.Write([]byte(	
`
#Compiled with IC.
FUNCTION copy
	GRAB array
	ARRAY c
	
	SHARE array
	RUN len
	PULL i+output+2
	
	VAR i
	LOOP
		VAR i+shunt+1
		SGE i+shunt+1 i i+output+2
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
`	))
	functions["copy"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{STRING}}

	output.Write([]byte(	
`	
FUNCTION output_m_file
	GRAB text
	SHARE text
	OUT
END
`))
	functions["output_m_file"] = Function{Exists:true, Args:[]TYPE{STRING}}
	methods["output"] = true

	
	
	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION close
	TAKE file
	CLOSE file
END
`	))
	functions["close"] = Function{Exists:true, Args:[]TYPE{FILE}}

	//Inbuilt output function.
	output.Write([]byte(
`
FUNCTION len
	GRAB data
	PUSH #data
END
`	))
	functions["len"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{NUMBER}}

	//Inbuilt reada function.
	output.Write([]byte(
`
FUNCTION reada
	PULL delim
	MUL delim delim -1
	PUSH delim
	STDIN
END
`	))
	functions["reada"] = Function{Exists:true, Args:[]TYPE{NUMBER}, Returns:[]TYPE{STRING}}
	
	//Inbuilt reada function.
	output.Write([]byte(
`
FUNCTION read
	PUSH 0
	STDIN
END
`	))
	functions["read"] = Function{Exists:true, Returns:[]TYPE{STRING}}
	
	//Inbuilt reada function.
	output.Write([]byte(
`
FUNCTION reada_m_file
	PULL delim
	MUL delim delim -1
	PUSH delim
	IN
END
`	))
	functions["reada_m_file"] = Function{Exists:true, Args:[]TYPE{NUMBER}, Returns:[]TYPE{STRING}}
	methods["reada"] = true
	
	//Inbuilt reada function.
	output.Write([]byte(
`
FUNCTION read_m_file
	PUSH 0
	IN
END
`	))
	functions["read_m_file"] = Function{Exists:true, Returns:[]TYPE{STRING}}
	methods["read"] = true
	
	//Inbuilt reada function.
	output.Write([]byte(
`
FUNCTION input_m_file
	TAKE file
	PULL delim
	ARRAY input
	
	VAR i
	
	PUSH delim
	RELAY file
	IN
	
	VAR condition
	
	LOOP
		ADD i i 1
		SGT condition i delim
		IF condition
			BREAK
		END
		
		PULL byte
		
		VAR byte==n1000
		SEQ byte==n1000 byte -1000
		IF byte==n1000
			ERROR 1
			BREAK
		END
		
		PLACE input
				PUT byte
	REPEAT
	SHARE input
END
`	))
	functions["input_m_file"] = Function{Exists:true, Args:[]TYPE{NUMBER}, Returns:[]TYPE{STRING}}
	methods["input"] = true

	//Inbuilt reada function.
	output.Write([]byte(
`
FUNCTION reada_m_string
	GRAB s
	PULL n

	ARRAY i+string+1
	SHARE i+string+1
	GRAB result
	VAR i
	ADD i 0 0
	LOOP
		SHARE s
		RUN len
		PULL i+output+3
		VAR i+shunt+2
		SLT i+shunt+2 i i+output+3
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
				SHARE s
				RUN len
				PULL i+output+7
				VAR i+shunt+6
				SLT i+shunt+6 j i+output+7
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
`	))
	functions["reada_m_string"] = Function{Exists:true, Args:[]TYPE{NUMBER}, Returns:[]TYPE{STRING}}

	//Inbuilt reada function.
	output.Write([]byte(
`
FUNCTION info_m_file
	STAT
END
`	))
	functions["info_m_file"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{STRING}}
	methods["info"] = true
	
	//Inbuilt num function.
	output.Write([]byte(
`
FUNCTION number
	GRAB string
	
	VAR num
	VAR tens
	ADD tens 0 1
	
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
			GET tens*i
		SEQ __condition tens*i 45
		IF __condition 
			MUL num num -1
			BREAK
		END
		
		
		SGT __toobig tens*i 57
		SLT __toosmall tens*i 46
		ADD __invalid __toobig __toosmall
		IF __invalid
			ERROR 1
			PUSH 0
			RETURN
		END
		
		#Convert from unicode.
		SUB tens*i tens*i 48
		MUL tens*i tens tens*i
		
		ADD num num tens*i
		
		MUL tens tens 10
	REPEAT
	
	PUSH num
END
`	))
	functions["number"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{NUMBER}}

	//Inbuilt text function.
	output.Write([]byte(
`
FUNCTION text
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
END
`	))
	functions["text"] = Function{Exists:true, Args:[]TYPE{NUMBER}, Returns:[]TYPE{STRING}}
	methods["text"] = true
	
	//Inbuilt text function.
	output.Write([]byte(
`
FUNCTION binary
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
	
	VAR twos
	VAR twos>num
	VAR num<0
	
	ADD twos twos 1
	
	SLT num<0 num 0
	IF num<0
		PLACE txt
			PUT 45
		MUL num num -1
	END
	
	#What is the highest power to 10 which fits in num.
	LOOP
		SGT twos>num twos num
		IF twos>num 
			DIV twos twos 2
			BREAK
		END
		
		MUL twos twos 2
	REPEAT
	
	VAR num/twos
	VAR twos*(num/twos)
	VAR twos<=0
	
	#Find each digit.
	LOOP
		SLE twos<=0 twos 0
		IF twos<=0  
			BREAK
		END
		DIV num/twos num twos
		MUL twos*(num/twos) twos num/twos
		SUB num num twos*(num/twos)
		
		ADD num/twos num/twos 48
		PLACE txt
			PUT num/twos
		
		DIV twos twos 2
	REPEAT
	SHARE txt
END
`	))
	functions["binary"] = Function{Exists:true, Args:[]TYPE{NUMBER}, Returns:[]TYPE{STRING}}
	
	functions["text_m_string"] = Function{Exists:true, Args:[]TYPE{}, Returns:[]TYPE{STRING}, Ghost:true}
	
	//Hash function.
	output.Write([]byte(
`
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
END
`	))

	output.Write([]byte(
`
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
SHARE alist
RUN len
PULL i+output+19
VAR i+shunt+20
SLE i+shunt+20 i+output+19 1
IF i+shunt+20
RETURN
END
SHARE alist
PUSH 0
SHARE alist
RUN len
PULL i+output+21
VAR i+shunt+22
SUB i+shunt+22 i+output+21 1
PUSH i+shunt+22
RUN i_part
RETURN
`))
	functions["sort"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{}}

//Hash function.
	output.Write([]byte(
`
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
END
`	))
	//functions["hash"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{NUMBER}}

output.Write([]byte(
`
FUNCTION watch
	GRAB id
	
	ARRAY i+tmp+1
	SHARE i+tmp+1
	RUN open
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
	RUN output_m_file
	RELAY grabserver
RETURN`))
	functions["watch"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{}}

output.Write([]byte(
`
FUNCTION grab
	GRAB id
	
	ARRAY i+tmp+1
	SHARE i+tmp+1
	RUN open
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
	RUN output_m_file
	RELAY grabserver
	RELAY grabserver
	PUSH 1
	RUN reada_m_file
	RELAY grabserver
	GRAB i+output+13
	PUSH -1
	RUN reada_m_file
	GRAB a
	SHARE i+output+13
RETURN
`))
	functions["grab"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{STRING}}

output.Write([]byte(
`
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
	RUN open
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
	RUN output_m_file
	RELAY server
RETURN
`))
	functions["gui"] = Function{Exists:true, Args:[]TYPE{STRING}, Returns:[]TYPE{}}

	output.Write([]byte(`
FUNCTION edit
	GRAB txt
	GRAB id
	ARRAY i+tmp+1
	SHARE i+tmp+1
	RUN open
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
	RUN output_m_file
	RELAY grabserver
RETURN
`))
	functions["edit"] = Function{Exists:true, Args:[]TYPE{STRING, STRING}, Returns:[]TYPE{}}
}
