run:
	./build
docker:
	docker run -it --rm -p 3000:80  -v "$(PWD)":/go/src/gogit -w /go/src/gogit golang:1.9.3-alpine3.7 /go/src/gogit/build.sh
dokcer_image:
	docker run -it --rm -p 3000:80 gogit
