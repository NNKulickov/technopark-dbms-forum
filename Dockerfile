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
RUN echo "listen_addresses='*'\nsynchronous_commit = off\nfsync = off\nshared_buffers = 256MB\neffective_cache_size = 1536MB\n" >> /etc/postgresql/$PSQLVer/main/postgresql.conf
RUN echo "wal_buffers = 16MB\nmax_wal_size = 2GB\nwal_writer_delay = 50ms\nrandom_page_cost = 1.0\nmax_connections = 100\nwork_mem = 8MB\nmaintenance_work_mem = 128MB\ncpu_tuple_cost = 0.0030\ncpu_index_tuple_cost = 0.0010\ncpu_operator_cost = 0.0005" >> /etc/postgresql/$PSQLVer/main/postgresql.conf
RUN echo "full_page_writes = off" >> /etc/postgresql/$PSQLVer/main/postgresql.conf


#COPY postgresql.conf /etc/postgresql/$PSQLVer/main/postgresql.conf
COPY DBScript DBScript

CMD service postgresql start && ./api


#