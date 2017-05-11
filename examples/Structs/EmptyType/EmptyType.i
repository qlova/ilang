type Empty {}

method text(Empty) "" {
	return "nothing"
}

software {
    var e = Empty()
    
    var breakloop = 100
    loop {
         e = Empty() //This should optimise away. No arrays should be allocated.
         breakloop--
         if breakloop <= 0 
              break
         end
    }
    print(e)
}
