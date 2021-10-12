# Build the helm-passwordsafe-binary
FROM golang:1.17.1-alpine3.14 as builder
WORKDIR /go/src/github.com/snooyen/animal-facts/factscraper
ADD . .
ARG TARGETARCH
ARG TARGETOS
ARG LDFLAGS
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -a -ldflags "${LDFLAGS}" -o bin/fact-scraper ./cmd/

FROM golang:1.17.1-alpine3.14
WORKDIR /app
COPY --from=builder /go/src/github.com/snooyen/animal-facts/factscraper/bin/fact-scraper /bin/
RUN chmod a+x /bin/fact-scraper
