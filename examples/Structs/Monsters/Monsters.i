type Vector {
	X 
	Y
	
	convert text {
		return "("+text(X)+","+text(Y)+")"
	}
}

type Monster {
	Name = "Unknown"
	HP = 100
	
	Pos = Vector{
		X = 20
	}
	
	convert text {
		return Name
	}
}

software {
	m = Monster()
	print(m.Name)
	
	print(m.HP)
	print(m.Pos)
	print(m.Pos.X)
}
