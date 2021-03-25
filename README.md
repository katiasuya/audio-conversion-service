# Audio-converter

Audio-converter is a service that exposes a RESTful API to convert WAV to MP3 and vice versa. 

## Architecture Diagram

![diagram](docs/architecture.jpeg)

## DataBase

First, download PostgreSQL server of the version 13.x, install it on your system and run it.

To create a PostgreSQL user, database, schema and tables needed for the service,

* go to the scripts folder of this repository as `cd scripts/`;
* set execute permission on the script as `chmod +x create_db.sh`;
* run the script as `./create_db.sh` and enter postgres password and user's name and password when asked.

## Configuration

To set the configuration parameters for the application, create a .env file in the root of the repository
and set the following environment variables:

AUDIO-CONVERTER_HOST=your_host(default localhost)
AUDIO-CONVERTER_PORT=your_port(default 5432)
AUDIO-CONVERTER_USERNAME=your_username
AUDIO-CONVERTER_PASSWORD=your_password
AUDIO-CONVERTER_DBNAME=audioconverter
AUDIO-CONVERTER_SSLMODE=disable
AUDIO-CONVERTER_STORAGE_PATH=your_storage_path
AUDIO-CONVERTER_SECRETKEY=your_secret_key

## Conversion

The service uses `ffmpeg` multimedia framework for audio conversion, so it needs to be installed.
Go to `https://www.ffmpeg.org/download.html` and follow the instructions to download it.
