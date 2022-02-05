TARGET = lang
BINDIR = bin

.DEFAULT_GOAL = build

clean:
	rm -rf $(BINDIR)

build:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/$(TARGET) .
