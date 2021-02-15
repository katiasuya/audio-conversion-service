# Audio-converter

Audio-converter is a service that exposes a RESTful API to convert WAV to MP3 and vice versa. 

## Architecture Diagram

![diagram](docs/architecture.jpeg)

## DataBase

First, download PostgreSQL server of the version 13.x, install it on your system and run it.

To create a PostgreSQL user, database, schema and tables needed for the service, run 
`make create-database  [USERNAME=username] PASSWORD=password DB_USERNAME=db_username DB_PASSWORD=db_password [DB_HOST=db_host] [DB_PORT=db_port]` with your data,
where USERNAME is the name of the user with privileges to create databases, schemas and users, default postgres;
      PASSWORD is the password of the user above;     
      DB_USERNAME is the name of the user to be created for the service;
      DB_PASSWORD is the password of the user above;
      DB_HOST postgres host, default localhost;
      DB_PORT postgres host port, default 5432.
