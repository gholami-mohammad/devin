run:
	./build.sh
docker:
	docker run -it --rm -p 13000:13000  -v "$(PWD)":/go/src/devin -w /go/src/devin golang:1.9.3-alpine3.7 /go/src/devin/build.sh

docker-go:
	docker run -it --rm -p 13000:13000 devin:go

docker-light:
	env GOOS=linux GOARCH=amd64 go build -v devin
	docker run -it -p 13000:13000 --rm devin:light

jwt_rsa_keys:
	openssl genrsa -out auth/keys/jwt.rsa 4096
	openssl rsa -in auth/keys/jwt.rsa -pubout > auth/keys/jwt.rsa.pub
	echo "package keys\n" > auth/keys/jwt_rsa.go
	echo "const JWT_RSA_PRIVATE string = \` " >> auth/keys/jwt_rsa.go
	cat auth/keys/jwt.rsa >> auth/keys/jwt_rsa.go
	echo "\`\n" >> auth/keys/jwt_rsa.go
	echo "const JWT_RSA_PUBLIC string = \` " >> auth/keys/jwt_rsa.go
	cat auth/keys/jwt.rsa.pub >> auth/keys/jwt_rsa.go
	echo "\`\n" >> auth/keys/jwt_rsa.go

test_user:
	go test -v --coverprofile=cover.out devin/modules/user/controllers
	go tool cover --html=cover.out

test_user_profile_update:
	go test -v --coverprofile=cover.out devin/modules/user/controllers -run=TestUpdateProfile
	go tool cover --html=cover.out

test_middlewares:
	go test -v --coverprofile=cover.out devin/middlewares
	go tool cover --html=cover.out

test_crypto:
	go test -v --coverprofile=cover.out devin/crypto
	go tool cover --html=cover.out
