
test:
	go test -v ./...

vet:
	go vet ./...

format-fix:
	gofmt -w ${GOFMT_FLAGS} .

import-fix:
	goimports -w .