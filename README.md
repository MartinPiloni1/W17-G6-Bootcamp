# W17-G6-Bootcamp

## Fresh Products API

## Sprint 1

| Nombre           | Feature    |
|------------------|------------|
| Pedro Crespo     |    #1      |
| Laura Mosquera   |    #2      |
| Daniel Villegas  |    #3      |
| Martin Piloni    |    #4      |
| Tomas Bardin     |    #5      |
| Agustin Aguero   |    #6      |


## Sprint 2

| Nombre           | Feature    |
|------------------|------------|
| Pedro Crespo     |    #1      |
| Laura Mosquera   |    #2      |
| Daniel Villegas  |    #3      |
| Martin Piloni    |    #4      |
| Tomas Bardin     |    #5      |
| Agustin Aguero   |    #6      |

### Use in local

using the local mysql (not in docker) user host and password

you need to export these env variables to your terminal or create a .env file like this

```
ADDRESS
DB_HOST
DB_PORT
DB_NAME=fresh  # must be this name for the migration scripts
DB_USER
DB_PASS
```
dbName cannot be changed, because how the migration is implemented

install dependencies
`go mod tidy`

then create the tables
`go run cmd/migrate/main.go`

then insert regiters to the table
`go run cmd/seed/main.go`


### Use in docker-compose

with docker and docker-compose, that reads a `.env` file or the root of the project with this structure
if not provided it will asume the default values in the compose 

```
ADDRESS
DB_PORT
DB_USER
DB_PASS
```

run the project with

`docker-compose up --build -d`