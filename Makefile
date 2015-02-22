COMMIT:=$(shell git log -1 --pretty=format:'%H')
BRANCH:=$(TRAVIS_BRANCH)

ifeq ($(strip $(BRANCH)),)
	BRANCH:=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
endif

all: clean iagod iago

clean:

	rm -rf ./bin
	rm -rf ./release

iagod:

	mkdir -p ./bin
	
	go build iagod/main.go
	mv main ./bin/iagod

iago:

	mkdir -p ./bin

	go build iago/main.go
	mv main ./bin/iago

release: iagod iago

	mkdir release
	cd bin && zip -r ../dist.zip .

	cp dist.zip release/$(COMMIT).zip
	cp dist.zip release/$(BRANCH).zip

	rm dist.zip

.PHONY: clean iagod iago
