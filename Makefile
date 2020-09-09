APP=tweets-en-vivo

build:
	GOOS=linux go build -o ${APP}
	docker build --tag ${APP} .
	rm -f ${APP}
