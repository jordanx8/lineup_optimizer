FROM golang:1.16.5 AS build-env

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY ./ ./
RUN go mod download


RUN go build -o /lineup-optimizer .

FROM golang:1.16.5

WORKDIR /app

RUN apt update && apt -y upgrade 
RUN apt -y install chromium

COPY templates ./
COPY static ./
COPY --from=build-env /lineup-optimizer ./

EXPOSE 8080

CMD [ "./lineup-optimizer" ]