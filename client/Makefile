COMMIT:=$(shell git log -1 --pretty=format:'%H')
BRANCH:=$(TRAVIS_BRANCH)

ifeq ($(strip $(BRANCH)),)
	BRANCH:=$(shell git branch | sed -n -e 's/^\* \(.*\)/\1/p')
endif

all: clean milou

clean:

	rm -rf ./bin
	rm -rf ./release

milou:

	mkdir -p ./bin
	
	go build main.go
	mv main ./bin/milou

release: clean milou

	mkdir release
	cd bin && zip -r ../dist.zip .

	cp dist.zip release/$(COMMIT).zip
	cp dist.zip release/$(BRANCH).zip

	rm dist.zip

.PHONY: clean milou
