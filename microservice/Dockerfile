FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /src

COPY ./infra ./infra
COPY ./internal ./internal
COPY ./utils ./utils
COPY ./main.go .
COPY ./go.mod .
COPY ./go.sum .

RUN go mod download
RUN go build -buildvcs=false -o main.exe .

FROM alpine:3 AS runtime

RUN addgroup -S nonroot \
    && adduser -S nonroot -G nonroot

USER nonroot

COPY --from=builder src/main.exe .

EXPOSE 8082

ENTRYPOINT [ "./main.exe" ]
