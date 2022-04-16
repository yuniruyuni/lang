TARGET = lang
BINDIR = bin
TMPDIR = .tmp

.DEFAULT_GOAL = build

clean:
	rm -rf $(BINDIR)
	rm -rf $(TMPDIR)

build:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(TARGET) .

test: clean build
	# the `./` is needed for this test,
	# but I don't know why so please let me know
	# if you reader have opinion.
	go test ./...
	./test.sh

bench: clean build
	go test -bench Compile -o pprof/compile.bin -cpuprofile pprof/compile.out ./

showprof:
	go tool pprof -http=":8081" pprof/compile.bin pprof/compile.out
