SELECT 'CREATE DATABASE audioconverter'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'audioconverter');
\gexec

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname='acuser') THEN
        CREATE USER acuser WITH ENCRYPTED PASSWORD 'AcUser!';
        GRANT ALL PRIVILEGES ON DATABASE audioconverter TO acuser;
    END IF;
END$$;

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
id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
username TEXT UNIQUE NOT NULL,
password TEXT NOT NULL,
created TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
updated TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS converter.audio (
id UUID  DEFAULT gen_random_uuid() PRIMARY KEY,
name TEXT NOT NULL,
format format NOT NULL,
location TEXT NOT NULL,
created TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
updated TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL
);

CREATE TABLE IF NOT EXISTS converter.request (
id UUID  DEFAULT gen_random_uuid() PRIMARY KEY,
user_id UUID NOT NULL,
original_id UUID NOT NULL,
converted_id UUID,
created TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
updated TIMESTAMP WITHOUT TIME ZONE DEFAULT NOW() NOT NULL,
status status NOT NULL,
FOREIGN KEY (user_id) REFERENCES converter."user" (id) ON DELETE CASCADE,
FOREIGN KEY (original_id) REFERENCES converter.audio (id) ON DELETE CASCADE,
FOREIGN KEY (converted_id) REFERENCES converter.audio (id) ON DELETE CASCADE
);





