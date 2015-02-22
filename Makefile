all: clean iagod

clean:

	rm -rf ./bin

iagod: clean

	mkdir ./bin
	
	go build iagod/main.go
	mv main ./bin/iagod

.PHONY: clean iagod
