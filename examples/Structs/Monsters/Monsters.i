type Vector {
	X, Y
}

method text(Vector) "" {
	return "("+text(X)+","+text(Y)+")"
}

type Monster {
	""Name
	  HP
	{}Pos
}

method new(Monster) {
	Name = "Unknown"
	HP   = 100
	Pos  = Vector{2}
}

software {
	var m = new(Monster)
	print(m.Name)
	
	print(m.HP)
	print(m.Pos)
	print(m.Pos.X)
}
