type Test {
	a,b
}

method Swap(Test) {
	swap(a, b)
}

software {
	var s = Test{2, 4}
	
	Swap(s)
	
	var a, b = 2, 4
	swap(a, b)
	
	print(s.a, " ", s.b)
	print(a, " ", b)
}
