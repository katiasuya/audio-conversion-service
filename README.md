# Audio-converter

Audio-converter is a service that exposes a RESTful API to convert WAV to MP3 and vice versa.

## Architecture Diagram

![diagram](docs/architecture.jpeg)

## Configuration

To set the configuration parameters for the application, set the following environment variables:

[1]  
```bash
CONVERTER_HOST=your_host 
CONVERTER_PORT=your_port 
CONVERTER_USER=your_user 
CONVERTER_PASSWORD=your_password  
CONVERTER_DB=audioconverter  
CONVERTER_SSLMODE=disable  
```
[2]  
```bash
CONVERTER_PRIVATEKEY="`cat your_private_key_path`"
CONVERTER_PUBLICKEY="`cat your_public_key_path`" 
```
[3]  
```bash
CONVERTER_ACCESSKEYID=your_access_key_id  
CONVERTER_SECRETACCESSKEY=your_secret_access_key  
CONVERTER_REGION=your_region  
CONVERTER_BUCKET=your_bucket_name  
```
[4]  
```bash
CONVERTER_URI=your_ampq_uri
CONVERTER_QUEUENAME=your_queue_name 
```
[5]*
```bash
POSTGRES_HOST=postgres.local
POSTGRES_PORT=your_port
POSTGRES_USER=your_user
POSTGRES_PASSWORD=your_password
POSTGRES_DB=audioconverter
POSTGRES_SSLMODE=disable
```

## DataBase

First, set environment variables from group [1].

Download PostgreSQL server of the version 13.x, install it on your system  
and run it with the corresponding configuration data.  

To create a PostgreSQL user, database, schema and tables needed for the service,  

* go to the scripts folder of this repository as `cd scripts/`;  
* set execute permission on the script as `chmod +x create_db.sh`;  
* run the script as `./create_db.sh` and enter postgres password and user's name and password when asked.  

## Authorization

Private and public keys that are used to implement authorization must be stored in `.pem` files.  
To create them, download OpenSSL library for your OS and run  
`openssl genpkey -algorithm RSA -out private_key_filename.pem -pkeyopt rsa_keygen_bits:2048`  
to generate a private key;  
`openssl rsa -pubout -in private_key_filename.pem -out public_key_filename.pem`  
to generate a public key from the given private key.  

Then set corresponding environment variables from group [2].  

## Storage

To store original and converted files for the service, AWS Simple Storage Service (Amazon S3) is used.  
For that, configure the credentials of the user with access to the bucket and set corresponding  
environment variables from group [3].  

## Conversion

The service uses `ffmpeg` multimedia framework for audio conversion, so it needs to be installed.  
Go to `https://www.ffmpeg.org/download.html` and follow the instructions to download it for your OS.

## Queuing

To use request queuing in the application, RabbitMQ is used.  
For that, set corresponding environment variables from group [4].  

## Docker

To run your application in docker, create an `.env` file at the root of the directory  
with all the values, decribed in Configuration section. 
Variables from group [5] are special for docker.
For more info read `https://hub.docker.com/_/postgres`.

Lastly, run `docker-compose up`.  
