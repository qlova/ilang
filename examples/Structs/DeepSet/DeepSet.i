type Pos {x, y}

type Monster {
	{}pos
}

software {
	var m = Monster()
	m.pos = Pos()
	
	m.pos.x = 3
	m.pos.x++
	m.pos.x--
	
	print(m.pos.x)
}
