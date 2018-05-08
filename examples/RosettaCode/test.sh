TESTING $1 $2 $3
BasicTest A+B "4" "2 2\n"
BasicTest Arithmetic/Integer "Sum: 8\nDifference: -2\nProduct: 15\nQuotient: 0\nModulus: 3\nExponent: 243" "3 5\n"
cd ..
BasicTest "Array concatenation" "[1,2,3,4,5,6]\n"
BasicTest "Array length" "2\n"
BasicTest "Arrays" "2\n4\n7\n"
BasicTest "Boolean values" "this prints\nthis prints\n"
BasicTest "Call a function" "\nInput a number!\n1234567890\nIt was: a\nmyprint is a not a real function\n[DEBUG] partial function!\nDebugPrint is a real function\n" "a\n"
BasicTest "Character codes" "97\na\n"
BasicTest "Check that file exists" "input.txt exists!\n/input.txt does not exist!\ndocs exists!\n/docs does not exist!\ndocs/Abdu'l-Bahá.txt exists!\n"
BasicTest "Conditional structures" "a = three\n"
BasicTest "Copy a string" "Hello Worlds\nHello World\n"

#Special Test
rm -r "Create a file/docs"
rm "Create a file/output.txt"
IgnoreTest "Create a file" "\n"
if [ -f  "Create a file/output.txt" ]; then
	if [ -d  "Create a file/docs" ]; then
		Passed
	else
		Failed
	fi
else
	Failed
fi

BasicTest "Detect division by zero" "5/0 is a division by zero.\n5/2 is not division by zero.\n0/0 is a division by zero.\n"
BasicTest "Determine if a string is numeric" "1200  is numeric :)\n3.14  is numeric :)\n3/4  is numeric :)\nabcdefg  is not numeric!\n1234test  is not numeric!\n"
BasicTest "Empty program" ""
BasicTest "Empty string" "Empty string!"
BasicTest "Enviroment variables" "$HOME\n$USER\n$PATH\n"
BasicTest "Factorial" "0\n1\n1\n2\n6\n1124000727777607680000\n"

#Special Test
rm "File input/output/output.txt"
IgnoreTest "File input/output" "\n"
cd ..
if [ -f  "File input/output/output.txt" ]; then
	if [ "$(cat 'File input/output/output.txt')" = "This is input text!" ]; then
		Passed
	else
		Failed
	fi
else
	Failed
fi

FakeTest "Find limit of recursion"

BasicTest "FizzBuzz" "1\n2\nFizz\n4\nBuzz\nFizz\n7\n8\nFizz\nBuzz\n11\nFizz\n13\n14\nFizzBuzz\n16\n17\nFizz\n19\nBuzz\nFizz\n22\n23\nFizz\nBuzz\n26\nFizz\n28\n29\nFizzBuzz\n31\n32\nFizz\n34\nBuzz\nFizz\n37\n38\nFizz\nBuzz\n41\nFizz\n43\n44\nFizzBuzz\n46\n47\nFizz\n49\nBuzz\nFizz\n52\n53\nFizz\nBuzz\n56\nFizz\n58\n59\nFizzBuzz\n61\n62\nFizz\n64\nBuzz\nFizz\n67\n68\nFizz\nBuzz\n71\nFizz\n73\n74\nFizzBuzz\n76\n77\nFizz\n79\nBuzz\nFizz\n82\n83\nFizz\nBuzz\n86\nFizz\n88\n89\nFizzBuzz\n91\n92\nFizz\n94\nBuzz\nFizz\n97\n98\nFizz\nBuzz\n"
BasicTest "Formatted numeric output" "00007.125\n-0007.125\n"
BasicTest "Function definition" ""
BasicTest "Greatest element of a list" "35757\n"

GraphicsTest "Hello world/Graphical"
BasicTest "Hello world/Text" "Hello world!"
cd ..
BasicTest "Increment a numerical string" "2\n"
BasicTest "Input loop" ""
BasicTest "Integer comparison" "7 is greater than 2\n" "7 2\n" 
BasicTest "Jump anywhere" "This will print\nHello there\n"
BasicTest "Integer comparison" "7 is greater than 2" '7 2\n'
IgnoreTest "Keyboard input/Flush the keyboard buffer" "\n"
