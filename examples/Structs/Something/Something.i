software {
	var a is Something
	
	a = [97]
	
	print(a) //This is slow.
	if a.type = number
		print("a is a number!")
		print(a.number) //This is fast.
	end
	if a.type = text
		print("a is text")
	end
	if a.type = array
		print("a is array")
	end
}
