//Interface definition.

type String { ""string }

method text(String) {
	return string
}

type SpecialString is String

type Printable ?{
	text() ""
}

software {
	var s = list(Printable)
	s += SpecialString{"This is a special string"}
	s += String{"This is a normal string"}
	
	for item in s
		print(item)
	end
}
