FROM ubuntu:22.04
WORKDIR /app
COPY ./bin/app_prod ./app
CMD ["./app"]
