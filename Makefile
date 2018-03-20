all:
	cd ./src/it && go build -o ../../it

install:
	cp ./it /usr/bin/it

windows:
	cd ./src/it && GOOS=windows go build -o ../../it.exe

gedit:
	cp ./doc/i.lang /usr/share/gtksourceview-3.0/language-specs/i.lang
