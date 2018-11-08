fmt:
	goimports -w ./..
	gofmt -w ./..

configure:
ifeq (, $(GOPATH))
	GOPATH=~/go
endif
	GO111MODULE=off; cd ~/go; go get -u golang.org/x/tools/cmd/goimports