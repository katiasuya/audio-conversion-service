#!/bin/bash

echo 'Please, enter username:'
read  DB_USERNAME
echo 'Please, enter password:'
read  DB_PASSWORD

psql \
	--variable=user_var=$DB_USERNAME \
	--variable=password_var=$DB_PASSWORD \
	-f schema.sql