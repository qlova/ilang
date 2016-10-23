# i
i is a hardware-agnostic maths-based cross-platform programming language in heavy development.

Hello World:

		software {
			print("Hello World")
		}
		
## Design
i is a language which is meant to be as clean and visually pleasing as possible, too many languages obfuscate code (Looking at you Java).
So I created i to solve this.

Here is a contrast between different languages and how visually pleasing they are at common tasks.

**File IO:**
Java:
```java
//This is ugly.
for (String filename : args) {
    try (FileReader fr = new FileReader(filename);BufferedReader br = new BufferedReader(fr)){
        String line;
        int lineNo = 0;
        while ((line = br.readLine()) != null) {
            processLine(++lineNo, line);
        }
    }
    catch (Exception x) {
        x.printStackTrace();
    }
}
```
Go:
```go
//Marginaly better
inputFile, err := os.Open("byline.go")
if err != nil {
	log.Fatal("Error opening input file:", err)
}
defer inputFile.Close()

scanner := bufio.NewScanner(inputFile)

for scanner.Scan() {
	fmt.Println(scanner.Text())
}

if err := scanner.Err(); err != nil {
	log.Fatal(scanner.Err())
}
```
Lua:
```lua
--Now we're talking! Still kinda obscure.
filename = "input.txt"
fp = io.open( filename, "r" )
 
for line in fp:lines() do
    print( line )
end
 
fp:close()
```
I:
```
	//Much better, well structured, clean and understandable.
	software {
		var file = open("input.txt")
		loop {
			print(read(file))
			issues {
				exit
			}
		}
	}
```

So you can hopefully see some of the philosophy of i.

##Features

* Big integars.
* Clean.
* Simple.
* Statically typed.
* Feels dynamic.
* Cross-platform.
* Has a Gui library.
* Has a **WIP** game engine.

##Documentation

Documentation is minimal but take a look at the [Wiki](https://github.com/Qlova/ilang/wiki).
The Rosetta Code [page](http://rosettacode.org/wiki/Category:I) may be helpful. 
Flick me an email (splizard @ splizard dot com) if you have any questions =)

##Types
There are 8 types in the i language.
```
number()
letter()	''
array() 	[]
text() 		""

pipe 		||
type		<>
function 	()

Something() {}
```

**number**
Numbers in the i language can be integers of any size.

**letter**
Letters correspond to a character in a piece of text, for example "a", "!" etc.

**array**
An array is a list of numbers.

**text**
Text is a string of letters, these can form words, sentences and the like.

**pipe**
A pipe is any object with an input and an output. These can be files, the internet etc.

**type**
A type is a representation of a type in the I language.

**function**
A function is a special form of pipe, one which you can call.

**Something** (WIP)
Something is the only builtin struture type, it can hold any value and needs to be type checked when used.  
Any types which are added in a program are "things".

##DOWNLOAD AND INSTALL

There are no official releases yet as the language is in a alpha state but you can grab the source and start hacking!  
Hacking instructions: (Linux and Mac)
	
		#OPTIONAL GUI SUPPORT
		git clone https://github.com/Qlova/grab
		cd grab && make && sudo make install
		grab #Needs to be started manually for now.
	
		#You need UCT.
		git clone https://github.com/Qlova/uct
		cd uct && make && sudo make install
		
		cd ../
		git clone https://github.com/Qlova/ilang
		cd ilang && make && sudo make install
		
		#Compile examples.
		cd examples && ./test.sh
		
		#You can now play around with the examples.
		#Standard building looks like this:
		ic File.i && uct -go File.u && go build File.go

Please be aware that many features are missing or incomplete in i!
		
