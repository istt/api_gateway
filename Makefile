# note: call scripts from /scripts
.DEFAULT_GOAL := all

GOBUILD=GOOS=linux GOARCH=amd64 go build 
GOSTATICBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"' 
# Clean all binaries
clean:
	rm -f build/package/*

# update dependencies
go-update:
	go get -u ./...

# === application specific commands ===
api-gateway:
	$(GOBUILD) -o build/package/api-gateway cmd/api-gateway/api-gateway.go

# === front-end tasks ===
# Install presiquition
npm:
	cd website; npm install; cd -

# Build angular application and copy it into /web folder
ng:
	cd website; npm run build; cd -

# Generate a pkged.go under project root folder for embed into the app
static: 
	pkger

# A combination of compile and static
ng-go: ng static
