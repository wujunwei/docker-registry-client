module github.com/heroku/docker-registry-client

go 1.16

require (
	github.com/docker/distribution v0.0.0
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/gorilla/mux v1.7.3 // indirect
	github.com/opencontainers/go-digest v1.0.0
	github.com/stretchr/testify v1.4.0 // indirect
)

replace github.com/docker/distribution => github.com/docker/distribution v0.0.0-20191216044856-a8371794149d
