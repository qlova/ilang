//The rational type in 'I' can represent any rational number.
software {
	//A rational number can be created by using a single backslash character.
	var r = 1\2
	
	//All the standard arithmatic operations are available, + - * /
	r *= 3\1
	r /= 1\4
	
	r += 1\2
	r -= 1\1
	
	print(r) //-> 11/2
	
	//Rational numbers can be cast to traditional numbers. This will truncate the rational number.
	print(number(r)) //-> 5
	
	//Numbers can also be cast back into rationals.
	print(rational(6)) //-> 6/1
}
