//Interface definition.

type Printable ?{
	text() ""
}

type Test {}
method text(Test) {
	return ("Test")
}

software {
	var s = list(Printable)
	s += 4
	s += "lol"
	s += 'a'
	s += Test()
	
	for item in s
		print(item)
	end
}
