

type Complex {
	real, imag
}

method Complex * Complex {
	c.real = a.real*b.real - a.imag*b.imag
	c.imag = a.real*b.imag + b.real*a.imag
}


method Complex + Complex {
	c.real  = a.real + b.real
	c.imag  = a.imag + b.imag
}

method Complex - Complex { 
	c.real  = a.real - b.real
	c.imag  = a.imag - b.imag
}

method Complex / Complex {
	var d = b.real² + b.imag²
	c.real  = (a.real*b.real + a.imag*b.imag) / d
	c.imag = (a.imag*b.real - a.real*b.imag) / d
}

method text(Complex) "" {
	return text(real)+" + "+text(imag)+"i"
}


software {

	var n is Complex{3, 6}
	
	var m is Complex
	m.real = 2
	m.imag = 1
	
	var b = n+m
	
	print(b)
	
	
}
