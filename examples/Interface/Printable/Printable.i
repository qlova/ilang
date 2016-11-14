//Interface definition.

interface Printable {
	text() ""
}

software {
	var s has Printable(s)
	s += 4
	s += "lol"
	s += 'a'
	
	for item in s
		print(item)
	end
}
