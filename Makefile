version = 0.0.1
buildtags  = ""
imagetag = latest

branch=${CI_COMMIT_BRANCH}

dockerTarget = registry.gitlab.com/coollision/synology-videostation-reindexer

# region: stuff for pipelines
Testing:
	echo 'not needed'

Build: buildArm

Dockerize: dockerizeArm publishDockerArm

# endregion

BuildLocal:
	go build -v -tags "$buildtags" -ldflags="-X main.version=${version}" -gcflags "all=-N -l"
	./synology-videostation-reindexer

buildAndPushLocal: buildArm dockerizeArm publishDockerArm

buildArm:
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -v -tags "$buildtags" -ldflags="-X main.version=${version}" -gcflags "all=-N -l"

dockerizeArm:
	docker build . -t image

publishDockerArm:
	echo $(branch)
ifeq ($(branch), master)
	docker tag synology-videostation-reindexer:latest ${dockerTarget}:latest
	docker tag synology-videostation-reindexer:latest ${dockerTarget}:${branch}
	docker tag synology-videostation-reindexer:latest ${dockerTarget}:${version}
	docker push ${dockerTarget}:latest
	docker push ${dockerTarget}:${branch}
	docker push ${dockerTarget}:${version}
else
	docker tag image:latest ${dockerTarget}:${branch}-${version}
	docker push ${dockerTarget}:${branch}-${version}
endif
