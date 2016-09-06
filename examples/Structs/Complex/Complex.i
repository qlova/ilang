

type Complex {
	real, imag
}

Complex * Complex {
	c  = ab - aibi
	ci = abi + bai
}


Complex + Complex {
	c  = a + b
	ci = ai + bi
}

Complex - Complex { 
	c  = a - b
	ci = ai - bi
}

Complex / Complex {
	var d = b*b + bi*bi
	c  = (ab + aibi) / d
	ci = (aib - abi) / d
}

method text() [] {
	return text(Complex.real)+" + "+text(Complex.imag)+"i"
}


software {

	var n is Complex(3, 6)
	
	var m is Complex
	m.real = 2
	m.imag = 1
	
	var b = n+m
	
	print(b)
	
	
}
