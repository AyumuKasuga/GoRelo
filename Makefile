make: clean
	mkdir -p ./builds
	go build -o builds/gorelo -v

all: clean
	mkdir -p ./builds
	GOOS=linux go build -o builds/gorelo -v
	GOOS=darwin go build -o builds/gorelo_darwin -v
	GOOS=windows go build -o builds/gorelo -v

clean:
	rm -rf ./builds