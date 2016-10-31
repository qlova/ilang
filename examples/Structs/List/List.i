type Item {
	value
}

method text(Item) "" {
	return text(value)
}

type Bank {
	{..}Items
}

method new(Bank) {
	Items has Item(s)
}

method clear(Bank) {
	Items has Item(s)
}

software {
	var list has Item(s)
	list & Item{22}
	
	var b is Bank
	b.Items = list
	
	list[0] = Item{44}
	
	list & Item{33}
	
	print(list[0])
	print(b.Items[1])
	
	for item in list
		print(item)
	end
	
	clear(b)
	print(len(b.Items))
}
