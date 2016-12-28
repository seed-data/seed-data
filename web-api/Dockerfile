FROM golang:1.7.4

COPY . /go/src/github.com/seed-data/seed-data/web-api
WORKDIR /go/src/github.com/seed-data/seed-data/web-api

RUN go get -v && go build -v

CMD [
  "./main",
  "-v",
  "--auto-migrate"
  "--backfill-from", "https://seed_data.s3.amazonaws.com/seed_data.csv.gz",
  "--import-workers", "5",
  "--import-interval", "900s"
]
