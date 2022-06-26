FROM golang:1.18-alpine as builder

WORKDIR /app

COPY . ./

RUN go mod download

RUN go build -o /draw

FROM alpine:latest

COPY --from=builder ./draw ./draw

EXPOSE 3000

CMD [ "./draw" ]

