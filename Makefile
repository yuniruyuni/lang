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
	go test
	./test.sh
