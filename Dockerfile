FROM alpine

RUN apk update

# copy all project files (including the required .env file)
COPY . /bin/app
WORKDIR /bin/app

EXPOSE 5000

ENTRYPOINT ["/bin/app/tweets-en-vivo"]
