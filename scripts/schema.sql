SELECT 'CREATE DATABASE audioconverter'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'audioconverter');
\gexec

\c audioconverter

CREATE SCHEMA IF NOT EXISTS converter;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'format') THEN
        CREATE TYPE format AS ENUM ('MP3', 'WAV');
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'status') THEN
        CREATE TYPE status AS ENUM ('queued', 'processing','done');
    END IF;
END$$;

CREATE TABLE IF NOT EXISTS converter."user"(
id UUID PRIMARY KEY,
username TEXT NOT NULL,
password TEXT NOT NULL,
created TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
updated TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS converter.audio (
id UUID PRIMARY KEY,
user_id UUID NOT NULL,
name TEXT NOT NULL,
format format NOT NULL,
location TEXT NOT NULL,
created TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
updated TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
FOREIGN KEY (user_id) REFERENCES converter."user" (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS converter.request (
id UUID PRIMARY KEY,
original_id UUID NOT NULL,
converted_id UUID,
created TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
updated TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
status status DEFAULT 'queued',
FOREIGN KEY (original_id) REFERENCES converter.audio (id) ON DELETE CASCADE,
FOREIGN KEY (converted_id) REFERENCES converter.audio (id) ON DELETE CASCADE
);



