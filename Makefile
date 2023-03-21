cover:
	go test ./test/... -coverpkg=./service,./shared -coverprofile=test/coverage/cover.out
	go tool cover -html=test/coverage/cover.out -o test/coverage/coverage.html
coverage: cover
	go tool cover -func test/coverage/cover.out | grep total: