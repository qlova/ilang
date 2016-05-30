package main

import "io"

func builtin(output io.Writer) {
	//Inbuilt output function.
	output.Write([]byte(
`
SUBROUTINE output
	POPSTRING data
	PUSHSTRING data
	STDOUT 
END
`	))
	functions["output"] = Function{Exists:true, Args:[]bool{true}}

	//Inbuilt output function.
	output.Write([]byte(
`
SUBROUTINE len
	POPSTRING data
	PUSH #data
END
`	))
	functions["len"] = Function{Exists:true, Args:[]bool{true}, Returns:[]bool{false}}

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
	functions["reada"] = Function{Exists:true, Args:[]bool{false}, Returns:[]bool{true}}
	
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
		
		SUB tens*i tens*i 48 #Convert from unicode.
		MUL tens*i tens tens*i
		
		ADD number number tens*i
		
		MUL tens tens 10
	REPEAT
	
	PUSH number
END
`	))
	functions["num"] = Function{Exists:true, Args:[]bool{true}, Returns:[]bool{false}}

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
	functions["text"] = Function{Exists:true, Args:[]bool{false}, Returns:[]bool{true}}
	
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
	//functions["hash"] = Function{Exists:true, Args:[]bool{true}, Returns:[]bool{false}}
}
