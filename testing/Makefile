.PHONY: build build-openzeppelin openzeppelin
build: build-openzeppelin

# Build openzeppelin docker files with docker-compose
openzeppelin-docker-compose:
	cd openzeppelin && ./run-openzeppelin-tests.sh

build-openzeppelin:
	cd openzeppelin && docker build -t qtum/janus-openzeppelin-test -f ./openzeppelin-test.Dockerfile .
	cd openzeppelin && docker build -t qtum/janus-openzeppelin -f ./openzeppelin.Dockerfile .

init:
	git submodule init
	git submodule update
	cd openzeppelin/openzeppelin-contracts && yarn install

openzeppelin:
	cd openzeppelin/openzeppelin-contracts && truffle test --network testing

openzeppelin-docker: build
	docker run --rm qtum/janus-openzeppelin