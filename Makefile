COMMIT:=$(shell git log -1 --pretty=format:'%H')
BRANCH:=$(TRAVIS_BRANCH)

ifeq ($(strip $(BRANCH)),)
	BRANCH:=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
endif

all: clean iagod

clean:

	rm -rf ./bin
	rm -rf ./release

iagod:

	mkdir -p ./bin
	
	go build main.go
	mv main ./bin/iagod

release: clean iagod

	mkdir release
	cd bin && zip -r ../dist.zip .

	cp dist.zip release/$(COMMIT).zip
	cp dist.zip release/$(BRANCH).zip

	rm dist.zip

.PHONY: clean iagod
