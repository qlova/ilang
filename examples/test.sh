#! /bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd $DIR

if [ "$1" = "" ]; then
	LANGUAGE=py
else
	LANGUAGE=$1
fi

function runit {
	case $LANGUAGE in
		py)
			python3 $1.py <<< $(echo -e "$2")
		;;
		go)
			go run $1.go <<< $(echo -e "$2")
		;;
		bash) 
			./$1.bash <<< $(echo -e "$2")
		;;
		java) 
			javac $1.java && java $1 <<< $(echo -e "$2")
		;;
		cs) 
			mcs -nowarn:414 /r:System.Numerics.dll $1.cs && mono $1.exe <<< $(echo -e "$2")
		;;
		rb)
			ruby $1.rb <<< $(echo -e "$2")
		;;
		lua)
			lua $1.lua <<< $(echo -e "$2")
		;;
	esac
}

function BasicTest {
	cd ../$1/ && ic $1.i && uct -$LANGUAGE $1.u
	if [ "$?" -eq "1" ]; then
		echo -e "$1 \e[31mFAILED!\e[0m to compile D:"
		exit 1
	fi
	local OUTPUT=$(runit $1 "$3")
	local DEFINED=$(echo -e "$2")
	if [ "$OUTPUT" = "$DEFINED" ]; then
		echo -e "$1 \e[32mPASSED!\e[0m"
	else
		echo -e "$1 \e[31mFAILED!\e[0m Got:"
		echo	 "$OUTPUT"
		echo "(Expecting)"
		echo	 "$DEFINED"
		exit 1
	fi
}

cd Plus
if [ "$2" != "" ]; then
	for l in py bash go lua rb java cs; do
		LANGUAGE=$l
		echo $l
		BasicTest $1 $2 
	done
	exit
fi

BasicTest Plus "4" "2 2\n"
BasicTest Arithmetic "Sum: 8\nDifference: -2\nProduct: 15\nQuotient: 0\nModulus: 3\nExponent: 243" "3 5\n"
BasicTest Arrays "2\n4"
BasicTest Chars "97\na"
BasicTest Concat "This string is joined!"
BasicTest Functions "hi\nhi\ncba"
BasicTest HelloWorld "Hello World"
BasicTest Length "2"
BasicTest OrderOfOperation "405"
BasicTest FileExists "input.txt exists\n/input.txt does not exist\ndocs exists\n/docs does not exist"
BasicTest Maths "d d b"
BasicTest Conditionals "3=3\n3!=2\nverified"
BasicTest Copy "2\n1"
BasicTest Variables "2"
