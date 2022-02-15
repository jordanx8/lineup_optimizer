FROM golang:latest AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY ./ ./
RUN go mod download


RUN go build -o /lineup-optimizer .

FROM chromedp/headless-shell:latest

WORKDIR /app

RUN apt-get update; apt install dumb-init -y

# RUN apt update && apt -y upgrade 
# RUN apt -y install chromium

ENTRYPOINT ["dumb-init", "--"]

COPY templates ./
COPY static ./
COPY --from=build /lineup-optimizer ./

EXPOSE 8080

CMD [ "./lineup-optimizer" ]