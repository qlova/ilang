type Item {
	value
	
	convert text {
		return text(value)
	}
	
	content number {
		value = number
	}
}

type Bank {
	items = list.Item()
	
	concept clear() {
		items = []
	}
	
	concept add() {
		items += Item(64)
	}
}

software {
	l = []
	l += Item(22)
	
	print(l[1])
	
	b = Bank()

	b.items = l
	
	b.items[0] = Item(12)

	l[0] = Item(44)
	
	l += Item(33)
	
	b.add()

	
	for item in l
		print(item)
	end

	b.clear()
	print(b.items.size())
}
