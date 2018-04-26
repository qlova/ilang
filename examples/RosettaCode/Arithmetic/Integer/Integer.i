software {
	a = number(read(' '))
	b = number(read('\n'))
	
	print("Sum: "		, a + b)
	print("Difference: ", a - b)
	print("Product: "	, a * b)
	print("Quotient: "	, a / b) // rounds towards zero
	print("Modulus: "	, a % b) // same sign as first operand
	print("Exponent: "	, a ^ b)
}
