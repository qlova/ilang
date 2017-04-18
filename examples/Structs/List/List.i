type Item {
	value
}

method text(Item) "" {
	return text(value)
}

type Bank {
	Items()
}

method clear(Bank) {
	items = Items()
}

software {
	var list = Items()
	list += Item{22}
	
	var b = Bank()
	b.items = list
	
	list[0] = Item{44}
	
	list += Item{33}
	
	print(list[0])
	print(b.items[1])
	
	for item in list
		print(item)
	end
	
	clear(b)
	print(len(b.items))
}
