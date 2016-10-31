
type Rational {
	numer, denom
}

function gcd(a, b) n {
	if b = 0
		return a
	else
		return gcd(b, a mod b)
	end
}

method Rational * Rational {
	c.numer  = a.numer*b.numer
	c.denom  = a.denom*b.denom
	
	var g = gcd(c.numer, c.denom)
	
	c.numer  = c.numer/g
	c.denom  = c.denom/g
}

method Rational + Rational {
	c.numer  = a.numer*b.denom + b.numer*a.denom
	c.denom  = a.denom*b.denom
	
	var g = gcd(c.numer, c.denom)
	
	c.numer  = c.numer/g
	c.denom  = c.denom/g
}

method Rational - Rational {
	c.numer  = a.numer*b.denom - b.numer*a.denom
	c.denom  = a.denom*b.denom
	
	var g = gcd(c.numer, c.denom)
	
	c.numer  = c.numer/g
	c.denom  = c.denom/g
}

method Rational / Rational {
	c.numer  = a.numer*b.denom
	c.denom  = a.denom*b.numer
	
	var g = gcd(c.numer, c.denom)
	
	c.numer  = c.numer/g
	c.denom  = c.denom/g
}

method text(Rational) "" {
	return text(numer)+"/"+text(denom)
}

software {
	var r = Rational()
	
	r.numer = 8
	r.denom = 16
	
	var n = Rational{2, 5}
	
	print(r)
	
	r = r * r
	
	print(r)
	
	r = r + n
	
	print(r)
}
