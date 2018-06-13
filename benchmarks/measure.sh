#! /bin/bash
cd "$1"

rm -rf ./.it/
echo
echo -n "Old Go target"
mv "$1.oi" "swap"
mv "$1.i" "$1.oi"
mv "swap" "$1.i"
~/go/bin/it build go
mv "$1.oi" "swap"
mv "$1.i" "$1.oi"
mv "swap" "$1.i"
/usr/bin/time -v "./.it/$1.gob" > /dev/null
rm -rf ./.it/

it build go
echo -n "Go target"
/usr/bin/time -v "./.it/$1.gob" > /dev/null
rm -rf ./.it/

echo
echo -n "Go native"
go build "$1.go"
/usr/bin/time -v "./$1" > /dev/null

echo
echo -n "Old Python target"
mv "$1.oi" "swap"
mv "$1.i" "$1.oi"
mv "swap" "$1.i"
~/go/bin/it build py
mv "$1.oi" "swap"
mv "$1.i" "$1.oi"
mv "swap" "$1.i"
/usr/bin/time -v "python3" "./.it/$1.py" > /dev/null
rm -rf ./.it/

echo
echo -n "New Python target"
it build py
/usr/bin/time -v "python3" "./.it/main.py" > /dev/null
rm -rf ./.it/

echo
echo -n "Python native"
/usr/bin/time -v "python3" "$1.py" > /dev/null

echo -n "Old Java target"
mv "$1.oi" "swap"
mv "$1.i" "$1.oi"
mv "swap" "$1.i"
~/go/bin/it build java
mv "$1.oi" "swap"
mv "$1.i" "$1.oi"
mv "swap" "$1.i"
cd .it
/usr/bin/time -v "java" "$1" > /dev/null
cd ..
rm -rf ./.it/

echo
echo -n "New Java target"
it build java
cd .it
/usr/bin/time -v "java" "Runtime" > /dev/null
cd ..
rm -rf ./.it/

echo
echo -n "Java native"
javac "$1.java"
/usr/bin/time -v "java" "$1" > /dev/null

echo -n "Old Javascript target"
mv "$1.oi" "swap"
mv "$1.i" "$1.oi"
mv "swap" "$1.i"
~/go/bin/it build js
mv "$1.oi" "swap"
mv "$1.i" "$1.oi"
mv "swap" "$1.i"
cd .it
/usr/bin/time -v "nodejs" "$1.js" > /dev/null
cd ..
rm -rf ./.it/

echo
echo -n "Javascript native"
/usr/bin/time -v "nodejs" "$1.js" > /dev/null
