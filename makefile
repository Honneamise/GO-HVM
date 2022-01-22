all: clean compile

clean:
	rm -f *.bin

compile:
	go build -o hvm.bin hvm/main.go
	go build -o hvmc.bin hvmc/main.go