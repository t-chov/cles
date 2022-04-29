FROM golang:1.18 as build-env

COPY . /usr/local/go/src/cles/
WORKDIR /usr/local/go/src/cles/
RUN CGO_ENABLED=0 go build -o /bin/cles .

FROM gcr.io/distroless/static-debian10

COPY --from=build-env /bin/cles /bin/cles
CMD ["cles"]