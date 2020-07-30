.PHONY: all clean

OUTPUT=ymprotoc

all: clean
	go build -o ./${OUTPUT} main.go

clean:
	rm -f ./${OUTPUT}