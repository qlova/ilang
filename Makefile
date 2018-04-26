all:
	go build -o ./it ./tools/it

install:
	cp ./it /usr/bin/it
