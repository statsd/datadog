
test:
	@go test -cover .

bench:
	@go test -bench=. -benchmem

.PHONY: bench test
