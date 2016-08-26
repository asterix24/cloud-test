all: deploy


clean:
	rm -rf "$(GOPATH)/pkg/darwin_amd64/fues3"

cross:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o testboard_daemon

run:
	godep go run `find . -name "*_test.go" -prune -o -path "./vendor" -prune -o -name "*.go" -print`

deploy:
	godep go build -o testboard_daemon



.PHONY: *
