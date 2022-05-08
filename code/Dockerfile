FROM golang:1.18 AS build
COPY . /go/src/api
RUN cd src/api \
    && CGO_ENABLED=0 go build -o /go/bin/api main.go

FROM debian:bullseye
ENV PSQLVer 13
ENV DB_HOST 0.0.0.0
ENV POSTGRES_USER subd_project
ENV POSTGRES_PASSWORD 123azxv
ENV POSTGRES_DB technopark

COPY --from=build /go/bin/api ./
RUN chmod +x api
RUN apt update && apt install -y tzdata
RUN ln -snf /usr/share/zoneinfo/Russia/Moscow /etc/localtime && echo Russia/Moscow > /etc/timezone
RUN apt update && apt install postgresql-$PSQLVer -y
RUN chmod -R u=rwx /var/lib/postgresql/$PSQLVer/main/
RUN chmod -R 0700 /etc/postgresql/$PSQLVer/main
USER postgres
RUN    /etc/init.d/postgresql start &&\
    psql --command "CREATE USER $POSTGRES_USER WITH SUPERUSER PASSWORD '$POSTGRES_PASSWORD';" \
       && createdb -O $POSTGRES_USER $POSTGRES_DB\
       && /etc/init.d/postgresql stop


RUN echo "host all  all 0.0.0.0/0  md5" >> /etc/postgresql/$PSQLVer/main/pg_hba.conf
RUN ls -la /etc/postgresql/
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PSQLVer/main/postgresql.conf
#COPY postgresql.conf /etc/postgresql/$PSQLVer/main/postgresql.conf
COPY DBScript DBScript

CMD service postgresql start && ./api


#