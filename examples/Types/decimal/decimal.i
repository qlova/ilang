software {

	//Decimals support all of the basic operations.
	d = decimal()
	d += 1
	d -= 0.5
	d *= 3
	d /= 2
	print(d) //0.75
	
	//By default, decimals have a precision of six decimal places. This can be adjusted to any length.
	d12 = decimal.12()
	d12 += 0.123456789012
	print(d12) //0.123456789012
	
	d2 = decimal.2()
	d2 += 0.99
	d2 /= 2
	print(d2) //0.49
	
	d4 = decimal.4()
	d4 += 0.99
	d4 /= 2
	print(d4) //0.495
	
	//Decimal equations can have mixed precisions, the resulting value will have the smallest precision.
	print(d+d12) //0.873456
	print(d+d2) //1.24
	
	//If you need to specify the cast, index the decimal with the precision you would like.
	print(d.12) //0.75
	print(d.12+d12) //0.873456789012
	
	//Casts from rational and number types are allowed.
	print(decimal(12\200)) //0.06
	print(decimal(5)) //5
	
	//It works backwards too!
	print(number(5.0)) //5/1
	print(rational(0.56)) //14/25
	
}
