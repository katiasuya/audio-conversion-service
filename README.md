# Audio-converter

Audio-converter is a service that exposes a RESTful API to convert WAV to MP3 and vice versa.   

## Architecture Diagram

![diagram](docs/architecture.jpeg)

## Configuration

To set the configuration parameters for the application, create a .env file at the root of the repository  
with the following environment variables:

[1]  

    AUDIO-CONVERTER_HOST=your_host(default localhost)  
    AUDIO-CONVERTER_PORT=your_port(default 5432)  
    AUDIO-CONVERTER_USERNAME=your_username  
    AUDIO-CONVERTER_PASSWORD=your_password  
    AUDIO-CONVERTER_DBNAME=audioconverter  
    AUDIO-CONVERTER_SSLMODE=disable  

[2]  

    AUDIO-CONVERTER_PRIVATEKEYPATH="your_private_key_path"  
    AUDIO-CONVERTER_PUBLICKEYPATH="your_public_key_path"  

[3]  

    AUDIO-CONVERTER_ACCESSKEYID=your_access_key_id  
    AUDIO-CONVERTER_SECRETACCESSKEY=your_secret_access_key  
    AUDIO-CONVERTER_REGION=your_region  
    AUDIO-CONVERTER_BUCKET=your_bucket_name  

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
