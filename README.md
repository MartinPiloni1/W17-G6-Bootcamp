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

you can populate the .env file as you need

```
ADDRESS
DB_HOST
DB_PORT
DB_NAME=fresh  # must be this name for the migration scripts
DB_USER
DB_PASS
```
dbName cant change because how the migration is implemented

install dependencies
`go mod tidy`

then create the tables
`go run cmd/migrate/main.go`

then insert regiters to the table
`go run cmd/seed/main.go`


### Use in docker-compose

with docker and docker-compose, that reads a `.env` file or the root of the project with this structure

```
ADDRESS=8080
DB_PORT=3306
DB_USER=freshuser # can be changed
DB_PASS=freshpass # can be changed
```

run the project with

`docker-compose up --build`