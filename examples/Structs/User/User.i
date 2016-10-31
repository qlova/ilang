type User {
	""Name, Age
	""Address
	""Phone
}

method text(User) "" {
	return Name
}

software {
	var bob = User()
	bob.Name = "Bob Normal"
	bob.Address = "22 Some Road"
	bob.Phone = "555000555"
	bob.Age = 33
	
	print(bob)
	print(bob.Age)
	print(bob.Phone)
	print(bob.Address)
}
