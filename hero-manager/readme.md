# Hero Manager

## Postgres Notes

```txt
docker run --name pgsrv --network psql -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres:alpine
docker run --name pgsql -it --rm --network psql postgres:alpine psql -h pgsrv -U postgres
create database heroes;
\c heroes;
docker run --name pgsql -it --rm --network psql postgres:alpine psql -h pgsrv -U postgres --dbname=heroes
postgres://postgres:mysecretpassword@localhost/heroes?sslmode=disable
export POSTGRES_DSN=postgres://postgres:mysecretpassword@localhost/heroes?sslmode=disable
```
