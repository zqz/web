.PHONY: test coverage coverage-html clean

test:
	go test -v ./...

coverage:
	go test -coverpkg=./... -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

coverage-html: coverage
	go tool cover -html=coverage.out -o coverage.html
	open coverage.html

clean:
	rm -f coverage.out coverage.html
