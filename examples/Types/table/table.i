 //Tables allow strings to be used as an index instead of a number.
 //In the future, they may allow an arbitary numerical index too.
 software {
 
    //Any type can be stored within a table as long as the use of the type is consistent.
    var t = table()
    t["a"] = "apple"
    t["b"] = "banana"
    
    print(t["a"]) //-> apple
    
    //Table values can be iterated over.
    for word in t
        print(word) //-> apple banana
    end
    
    //Table indices can also be iterated over.
    for index over t
        print(index) //-> a b
    end
 }
