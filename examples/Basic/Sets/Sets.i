software {
	//All these set operations are in theory very fast.
	var ingredients = set()
	ingredients += <chicken, flour, eggs, oats, butter, rice, pasta, cheese, milk, sugar>
	
	var biscuits = <flour, eggs, butter, sugar>
	var dinner = <steak>
	
	if biscuits <= ingredients
		print("You can make biscuits!")
	else
		print("You need to go shopping to make biscuits!")
	end
	
	if dinner <= ingredients
		print("You can have steak tonight!")
	else
		print("You need to go shopping to have steak tonight!")
	end
}
