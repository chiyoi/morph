FROM golang:1.21

WORKDIR /morph
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o ./out .

FROM alpine:3.18

WORKDIR /morph
COPY --from=0 /morph/out /bin/morph

ENV VERSION="v0.1.0"
ENV ENV="prod"
ENV ADDR=":http"
ENV ENDPOINT_AZURE_COSMOS="https://neko03cosmos.documents.azure.com:443/"
ENV ENDPOINT_AZURE_BLOB="https://neko03storage.blob.core.windows.net/"
ENV BLOB_CONTAINER_CERT_CACHE="morph-cert-cache"
ENV DATABASE="morph"
CMD ["morph"]
