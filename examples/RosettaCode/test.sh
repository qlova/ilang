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
BasicTest "Conditional structures" "a = three"
BasicTest "Copy a string" "Hello Worlds\nHello World"

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

