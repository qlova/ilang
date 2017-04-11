@nl

software {
	ver a = aantal(lezen(' '))
	ver b = aantal(lezen('\n'))
	
	afdrukken("Som: "		, a + b)
	afdrukken("Verschil: ", a - b)
	afdrukken("Artikel: "	, a * b)
	afdrukken("Quotiënt: "	, a / b)
	afdrukken("Modulus: "	, a mod b)
	afdrukken("Exponent: "	, a ^ b)
}
