VERSION=`git describe --tags`

bin:
	mkdir bin

# NOTE Avoids conflict with markov directory
.PHONY: markov
markov: bin
	go build -o bin/ -ldflags "-X main.version=${VERSION}" cmd/markov/markov.go

.PHONY: install
install: markov
	sudo cp ./bin/markov /usr/local/bin
