.PHONY: build build-docker publish-docker clean format lint vet staticcheck gosec govulncheck gofumpt


lint: vet staticcheck gosec govulncheck gofumpt

vet:
	go vet -v ./...
staticcheck:
	go install honnef.co/go/tools/cmd/staticcheck@latest
	staticcheck -checks all ./...
gosec:
	go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude=G304,G114 ./...
govulncheck:
	go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck ./...
gofumpt:
	go install mvdan.cc/gofumpt@latest
	test -z "$$(gofumpt -d -e . | tee /dev/stderr)"

build:
	go get -d -v
	CGO_ENABLED=0 go build -a -ldflags '-extldflags "-static"' -o dnsupd

format:
	go install mvdan.cc/gofumpt@latest
	gofumpt -w .

build-docker:
	docker build --build-arg GIT_COMMIT=$$CIRCLE_SHA1 -t docker.x-way.org/xway/dnsupd:latest .

publish-docker: build-docker
	echo $$DOCKER_ACCESS_TOKEN | docker login -u $$DOCKER_USERNAME --password-stdin docker.x-way.org
	docker push docker.x-way.org/xway/dnsupd:latest

clean:
	rm -f dnsupd
