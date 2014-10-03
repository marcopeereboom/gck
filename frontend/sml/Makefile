all: tokenizer.go lang.go
	go test -v

lang.go: lang.y
	go tool yacc -o lang.go.u lang.y
	gofmt lang.go.u > lang.go
	rm lang.go.u

tokenizer.go: tokenizer.l
	golex -t tokenizer.l | gofmt > tokenizer.go

clean:
	rm -rf tokenizer.go lang.go y.output lang.go.u
