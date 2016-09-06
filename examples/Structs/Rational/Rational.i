
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

Rational * Rational {
	c  = ab
	ci = aibi
	
	var g = gcd(c, ci)
	
	c = c/g
	ci = ci/g
}

Rational + Rational {
	c  = abi + bai
	ci = aibi
	
	var g = gcd(c, ci)
	
	c = c/g
	ci = ci/g
}

Rational - Rational {
	c  = abi - bai
	ci = aibi
	
	var g = gcd(c, ci)
	
	c = c/g
	ci = ci/g
}

Rational / Rational {
	c  = abi
	ci = aib
	
	var g = gcd(c, ci)
	
	c = c/g
	ci = ci/g
}

method text() [] {
	return text(Rational.numer)+"/"+text(Rational.denom)
}

software {
	var r is Rational
	
	r.numer = 8
	r.denom = 16
	
	var n is Rational(2, 5)
	
	print(r)
	
	r = r*r
	
	print(r)
	
	r = r + n
	
	print(r)
}
