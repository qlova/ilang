all:
	cd ./src && go build -o ../ic
	cd ./src/it && go build -o ../../it

install:
	cp ./ic /usr/bin/ic
	cp ./it /usr/bin/it
	cp ./doc/i.lang /usr/share/gtksourceview-3.0/language-specs/i.lang
