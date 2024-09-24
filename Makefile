
generate-cover:
	go test -v -coverpkg=github.com/ybru-tech/georm/... -coverprofile=cover.out ./...

go-cover-func: generate-cover
	@go tool cover -func cover.out
	@rm cover.out

go-cover-html: generate-cover
	go tool cover -html=cover.out
	@rm cover.out
