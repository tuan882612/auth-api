# How to run api

## Requirements

- [Golang](https://golang.org/)
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Running the API

1. Fill the `.env` file with the correct values in `assets/.env.example` and rename it to `.env`

2. Run the following command to start the API:

    ```bash
    # The API will be available at `http://localhost:8080`
    docker-compose up
    ```

    a. Run the sql script to create the database and tables

    ```bash
    docker exec -i auth-api_db_1 psql -U postgres -d postgres -a < assets/db.sql
    ```
