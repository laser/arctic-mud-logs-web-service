all:
	go build -o ./target/parser utils/parser.go
.PHONY: all

parse:
	find logs/ -name "*.txt" -type f | xargs -n 1 -P 8 -I {} sh -c "cat {} | ./target/parser > {}.meta"
.PHONY: parse

clean:
	rm -f ./target/parser
.PHONY: clean

server:
	go run ./main.go
.PHONY: server
