all: deploy


clean:
	rm -rf "$(GOPATH)/pkg/darwin_amd64/fues3"

deploy:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o goboard

run:
	go run `find . -name "*_test.go" -prune -o -path "./vendor" -prune -o -name "*.go" -print`

test:
	go build -o goboard

.PHONY: *
