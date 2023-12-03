FROM golang:1.21
WORKDIR /morph
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./out .

FROM alpine:3.18
WORKDIR /morph
COPY --from=0 /morph/out /bin/morph

ENV ENV=prod
CMD ["morph"]
