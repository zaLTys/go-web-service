# go-web-service

# db stuff
   
docker run --name entity-db-container -e POSTGRES_PASSWORD=mysecretpassword -d -p 5432:5432 postgres

docker exec -it entity-db-container psql -U postgres      

postgres=# CREATE DATABASE entities;
CREATE DATABASE

postgres=# create role entityuser with login password 'password';
CREATE ROLE

postgres=# \c entities
You are now connected to database "entities" as user "postgres".

entities=# CREATE TABLE IF NOT EXISTS entities (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    labels text[] NOT NULL,
    version integer NOT NULL DEFAULT 1,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE
entities=# grant select,insert,update,delete on entities to entityuser;
GRANT

entities=# grant usage, select on sequence entities_id_seq to entityuser;
GRANT


[System.Environment]::SetEnvironmentVariable(
    "ENTITIES_DB_DSN",
    "postgres://entityuser:password@localhost/entities?sslmode=disable",
    "User"
)


go run .\cmd\api\
2025/11/16 09:12:09 db conn pool established
2025/11/16 09:12:09 Starting dev server on :4000   


### RUN
go run .\cmd\api\ -db-dsn="postgres://entityuser:password@localhost:5432/entities?sslmode=disable"