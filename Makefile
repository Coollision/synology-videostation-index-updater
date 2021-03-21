version = "testing_stuff"
buildtags  = ""


Build: BuildAndRunLocal

RunDocker:
	docker buildx build --load -t test:new --platform linux/amd64 --build-arg buldtags=$buildtags --build-arg version="${version}" .
	docker run test:new

BuildAndRunLocal:
	go build -v -tags "$buildtags" -ldflags="-X main.version=${version}" -gcflags "all=-N -l"
	./synology-videostation-reindexer

buildArm64:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -v -tags "$buildtags" -ldflags="-X main.version=${version}" -gcflags "all=-N -l" -o synology-videostation-reindexer-arm64

