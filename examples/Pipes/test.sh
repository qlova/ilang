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
			go build && ./$1 <<< $(echo -e "$2")
		;;
		bash) 
			./$1.bash <<< $(echo -e "$2")
		;;
		java) 
			javac $1.java && java $1 <<< $(echo -e "$2")
		;;
		cs) 
			mcs -nowarn:414 /r:mscorlib.dll /r:System.Numerics.dll $1.cs && mono $1.exe <<< $(echo -e "$2")
		;;
		rb)
			ruby $1.rb <<< $(echo -e "$2")
		;;
		lua)
			lua $1.lua <<< $(echo -e "$2")
		;;
		js)nodejs $1.js <<< $(echo -e "$2")
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

cd CreateFile
if [ "$2" != "" ]; then
	for l in rb py go java lua rb js bash cs; do
		LANGUAGE=$l
		echo $l
		BasicTest "$1" "$2" "$3" 
	done
	exit
fi

BasicTest FileExists "input.txt exists\n/input.txt does not exist\ndocs exists\n/docs does not exist"
BasicTest CreateFile "output.txt created!\ndocs/ created!\nFailed to create /output.txt\nFailed to create /docs/"
BasicTest ReadFile "This is the contents of the file!"

