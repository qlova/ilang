//The number type in 'I' can store any integar value. 
software {

	//numbers support all the basic mathematical operations.
	n = number()
	n += 1
	n *= 6
	n /= 2
	n -= 1
	print(n) //-> 2
	
	//numbers also support power operations.
	n = n^2
	n = n^n
	print(n) //-> 256
	
	//a modulo operation is also supported.
	print(n % 200) //-> 56
	
	//Increment and decrement operators are also available.
	n++
	print(n) //-> 257
	n--
	print(n) //-> 256
	
	//There are no booleans in 'i'. Instead, numbers are used for boolean operations and conditions.
	a, b = 0, 0
	print(a = b) //-> 1
	print(a < b) //-> 0
	print(a > b) //-> 0
	print(a - b) //-> 0
}
