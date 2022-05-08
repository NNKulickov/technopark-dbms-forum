FROM golang:1.18 AS build
COPY . /go/src/api
RUN cd src/api \
    && CGO_ENABLED=0 go build -o /go/bin/api main.go

FROM postgres:bullseye
COPY --from=build /go/bin/api ./
RUN chmod +x api
COPY DBScript DBScript
ENV POSTGRES_USER postgres
ENV POSTGRES_PASSWORD 123azxv
ENV POSTGRES_DB technopark
EXPOSE 5000/tcp
CMD ["./api"]




#