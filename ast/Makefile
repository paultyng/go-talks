all:
	rm -f ssadump
	rm -f astdump
	rm -f helloworld.*.txt
	go build -o astdump astdump.go
	go build -o ssadump ssadump.go
	cat helloworld.go | ./astdump > helloworld.ast.txt
	cat helloworld.go | ./ssadump > helloworld.ssa.txt
	rm -f ssadump
	rm -f astdump