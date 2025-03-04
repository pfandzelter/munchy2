ROOT=github.com/pfandzelter/munchy2
GOFILES=$(shell find . -name "*.go")

.PHONY: deploy plan clean test

deploy: go-eat.zip munchy.zip main.tf init.done
	terraform apply
	touch $@

plan: go-eat.zip munchy.zip main.tf init.done
	terraform plan
	touch $@
	rm -f go-eat.zip go-eat
	rm -f munchy.zip munchy

init.done:
	terraform init
	touch $@

go-eat.zip: go-eat
	chmod +x go-eat
	zip -j $@ $<

munchy.zip: munchy
	chmod +x munchy
	zip -j $@ $<

go-eat: ${GOFILES}
	go get ${ROOT}/cmd/eat
	GOOS=linux GOARCH=amd64 go build -ldflags="-d -s -w" -o $@ ${ROOT}/cmd/eat

dev-eat: ${GOFILES}
	go get ${ROOT}/cmd/eat
	go build -o $@ ${ROOT}/cmd/eat

munchy: ${GOFILES}
	go get ${ROOT}/cmd/munchy
	GOOS=linux GOARCH=amd64 go build -ldflags="-d -s -w" -o $@ ${ROOT}/cmd/munchy

clean:
	terraform destroy
	rm -f init.done deploy.done go-eat.zip go-eat munchy.zip munchy

test:
	rm test.log || true
	go test -v ./...