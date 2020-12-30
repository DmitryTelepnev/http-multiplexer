FROM golang:1.14-stretch as builder

WORKDIR /srv/app
COPY ./ .
RUN make build

FROM scratch

WORKDIR /srv
COPY --from=builder /srv/app/bin/app /srv/app
COPY --from=builder /srv/app/config.yml /srv/config.yml
EXPOSE 8080

ENTRYPOINT ["/srv/app"]