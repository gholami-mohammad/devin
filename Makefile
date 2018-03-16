run:
	./build.sh
docker:
	docker run -it --rm -p 8080:8080  -v "$(PWD)":/go/src/gogit -w /go/src/gogit golang:1.9.3-alpine3.7 /go/src/gogit/build.sh

docker-go:
	docker run -it --rm -p 8080:8080 gogit:go

docker-light:
	env GOOS=linux GOARCH=amd64 go build -v gogit
	docker run -it -p 8080:8080 --rm gogit:light


test_user:
	go test -v --coverprofile=cover.out devin/modules/user/controllers
	go tool cover --html=cover.out
