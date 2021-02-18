#!/bin/bash

echo 'Please, enter postgres password:'
read  -s PASSWORD
echo 'Please, enter username:'
read  DB_USERNAME
echo 'Please, enter password:'
read  -s DB_PASSWORD

psql postgres://postgres:$PASSWORD@localhost:5432 \
	--variable=user_var=$DB_USERNAME \
	--variable=password_var=$DB_PASSWORD \
	-f schema.sql

