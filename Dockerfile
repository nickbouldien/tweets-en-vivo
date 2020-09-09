FROM alpine

RUN apk update

# copy all project files (including the required .env file)
ADD . .

EXPOSE 5000

ENTRYPOINT ["/tweets-en-vivo"]
