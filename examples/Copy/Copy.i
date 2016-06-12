software {
	var []a = [1, 2, 3]
	var []b = copy(a)
	
	a.0 = 2
	
	output(text(a.0)&"\n")
	output(text(b.0)&"\n")
}
