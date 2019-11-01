default: build

build:
	go build -o terraform-provider-macaddress-$(shell go env GOOS)-$(shell go env GOARCH)

testacc:
	TF_ACC=1 go test

clean:
	rm -f terraform-provider-macaddress-*
