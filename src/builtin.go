package main

import "io"

//TODO optimise some of the functions to be inline.
func builtin(output io.Writer) {
	//Inbuilt output function.
	output.Write([]byte(
`
SUBROUTINE output
	STDOUT 
END
`	))
	functions["output"] = Function{Exists:true, Args:[]int{STRING}}

	//Inbuilt output function.
	output.Write([]byte(
`
SUBROUTINE load
	LOAD
END
`	))
	functions["load"] = Function{Exists:true, Args:[]int{STRING}, Returns:[]int{STRING}}
	
	//Inbuilt output function.
	output.Write([]byte(
`
SUBROUTINE open
	OPEN file
	POP status
	
	IF status
		ERROR 1
	END
	
	PUSHIT file
END
`	))
	functions["open"] = Function{Exists:true, Args:[]int{STRING}, Returns:[]int{FILE}}
	
	output.Write([]byte(	
`
#Compiled with IC.
SUBROUTINE copy
	POPSTRING array
	STRING c
	VAR i 0
	LOOP
		PUSHSTRING array
		RUN len
		POP i+output+2
		VAR i+shunt+1
		SGE i+shunt+1 i i+output+2
		IF i+shunt+1
			PUSHSTRING c
			RETURN
		END
		INDEX array i i+shunt+3
		VAR v i+shunt+3
		PUSH v c
		ADD i i 1
	REPEAT
END
`	))
	functions["copy"] = Function{Exists:true, Args:[]int{STRING}, Returns:[]int{STRING}}

	output.Write([]byte(	
`	
SUBROUTINE output_m_3
	POPSTRING text
	POPIT self
	PUSHSTRING text
	OUT self
	POP status
	IF status
		ERROR 1
	END
END
`))
	functions["output_m_3"] = Function{Exists:true, Args:[]int{STRING}}
	methods["output"] = true

	
	
	//Inbuilt output function.
	output.Write([]byte(
`
SUBROUTINE close
	POPIT file
	CLOSE file
END
`	))
	functions["close"] = Function{Exists:true, Args:[]int{FILE}}

	//Inbuilt output function.
	output.Write([]byte(
`
SUBROUTINE len
	POPSTRING data
	PUSH #data
END
`	))
	functions["len"] = Function{Exists:true, Args:[]int{STRING}, Returns:[]int{NUMBER}}

	//Inbuilt reada function.
	output.Write([]byte(
`
SUBROUTINE reada
	POP delim
	STRING input
	VAR canbreak
	LOOP
		PUSH 1
		STDIN
		POP byte
		
	
		VAR byte==delim
		SEQ byte==delim byte delim
		IF byte==delim
			IF canbreak
				BREAK
			END
		ELSE
			ADD canbreak 0 1
			PUSH byte input
		END
	REPEAT
	PUSHSTRING input
END
`	))
	functions["reada"] = Function{Exists:true, Args:[]int{NUMBER}, Returns:[]int{STRING}}
	
	//Inbuilt num function.
	output.Write([]byte(
`
SUBROUTINE num
	POPSTRING text
	
	VAR number
	VAR tens 1
	
	VAR end #text
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
		
		
		INDEX text i tens*i
		SEQ __condition tens*i 45
		IF __condition 
			MUL number number -1
			BREAK
		END
		
		
		SGT __toobig tens*i 57
		SLT __toosmall tens*i 46
		ADD __invalid __toobig __toosmall
		IF __invalid
			ERROR 1
		END
		SUB tens*i tens*i 48 #Convert from unicode.
		MUL tens*i tens tens*i
		
		ADD number number tens*i
		
		MUL tens tens 10
	REPEAT
	
	PUSH number
END
`	))
	functions["num"] = Function{Exists:true, Args:[]int{STRING}, Returns:[]int{NUMBER}}

	//Inbuilt text function.
	output.Write([]byte(
`
SUBROUTINE text
	POP num
	STRING txt
	
	VAR test
	
	SEQ test num 0
	IF test
		PUSH 48 txt
		PUSHSTRING txt
		RETURN
	END
	
	VAR tens
	VAR tens>num
	VAR num<0
	
	ADD tens tens 1
	
	SLT num<0 num 0
	IF num<0
		PUSH 45 txt
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
		PUSH num/tens txt
		
		DIV tens tens 10
	REPEAT
	PUSHSTRING txt
END
`	))
	functions["text"] = Function{Exists:true, Args:[]int{NUMBER}, Returns:[]int{STRING}}
	
	//Hash function.
	output.Write([]byte(
`
SUBROUTINE hash
	POPSTRING text
	POP exp
	
	VAR number
	VAR tens 1
	
	VAR end #text
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
		
		
		INDEX text i tens*i
		
		MUL tens*i tens tens*i
		
		ADD number number tens*i
		
		MUL tens tens exp
	REPEAT
	
	PUSH number
END
`	))

//Hash function.
	output.Write([]byte(
`
SUBROUTINE unhash
	POP exp
	POP num
	STRING txt
	
	VAR test
	
	SEQ test num 0
	IF test
		PUSHSTRING txt
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

		PUSH num/tens txt
		
		DIV tens tens exp
	REPEAT
	PUSHSTRING txt
END
`	))
	//functions["hash"] = Function{Exists:true, Args:[]int{STRING}, Returns:[]int{NUMBER}}
}
