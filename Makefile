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
