//A pure implementation of 'i' has no limit of recursion.
//However, pure implemetations of 'i' cannot exist therefore
//  this will crash on different target platforms at some point. 

concept FindLimitOfRecursion(counter) {
	counter++
	FindLimitOfRecursion(counter)
}

software {
	FindLimitOfRecursion(0)
}
