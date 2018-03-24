type Item {
	value
}

method text(Item) "" {
	return text(value)
}

type Bank {
	..{Item} items
}

method clear(Bank) {
	items = []
}

method add(Bank) {
	items += Item{64}
}

software {
	var list = []
	list += Item{22}
	
	var b = Bank()
	b.items = list
	
	list[0] = Item{44}
	
	list += Item{33}
	
	print(list[0])
	print(b.items[1])
	
	add(b)
	
	for item in list
		print(item)
	end
	
	clear(b)
	print(len(b.items))
}
