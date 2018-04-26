
concept add(a, b) {
	return a+b
}

software {
	a = 50
	b = 50
	
	e = 2
	
	c = [a*e, 32, add(a, b), 32, 98]

	for value in c
		send(value, " ")
	end
	send("\n")
}
