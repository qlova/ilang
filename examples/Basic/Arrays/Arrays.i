software {
	a = []
	a += 2
	
	print(a[0]) //Outputs 2

	a[0] = 4

	print(a[0]) //Outputs 4
	
	print(a[1]) //Outputs 4, arrays wrap.
}
