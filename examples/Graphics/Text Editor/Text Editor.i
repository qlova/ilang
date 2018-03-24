import grate

type Graphics {
	""string
}

method update(Graphics) {
	for char in keys.pressed()
		if char
			string += char
		else
			string--
		end
	end
}

method draw(Graphics) {
	set.scale(10, 10)
	display(string, 0, 0)
}
