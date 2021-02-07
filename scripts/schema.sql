SELECT 'CREATE DATABASE audioconverter'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'audioconverter');
\gexec

\c audioconverter

CREATE SCHEMA IF NOT EXISTS converter;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'format') THEN
        CREATE TYPE format as ENUM ('MP3', 'WAV');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
        CREATE TYPE status as ENUM ('queued', 'processing','done');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS converter."user"(
id UUID PRIMARY KEY,
username TEXT NOT NULL,
password TEXT NOT NULL,
email TEXT,
created TIMESTAMP NOT NULL,
updated TIMESTAMP
);

CREATE TABLE IF NOT EXISTS converter.audio (
id UUID PRIMARY KEY,
user_id UUID NOT NULL,
name TEXT NOT NULL,
format format NOT NULL,
location TEXT NOT NULL,
created TIMESTAMP NOT NULL,
updated TIMESTAMP,
FOREIGN KEY (user_id) REFERENCES converter."user" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS converter.request (
id UUID PRIMARY KEY,
original_id UUID NOT NULL,
converted_id UUID,
created TIMESTAMP NOT NULL,
updated TIMESTAMP,
status status,
FOREIGN KEY (original_id) REFERENCES converter.audio (id) ON DELETE CASCADE,
FOREIGN KEY (converted_id) REFERENCES converter.audio (id) ON DELETE CASCADE
);



