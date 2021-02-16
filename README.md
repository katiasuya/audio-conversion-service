# Audio-converter

Audio-converter is a service that exposes a RESTful API to convert WAV to MP3 and vice versa. 

## Architecture Diagram

![diagram](docs/architecture.jpeg)

## DataBase

First, download PostgreSQL server of the version 13.x, install it on your system and run it.

To create a PostgreSQL user, database, schema and tables needed for the service, run 
`scripts/create_db.sh`.

