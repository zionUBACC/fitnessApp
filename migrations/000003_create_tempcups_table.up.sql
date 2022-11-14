-- Filename: migrations/000004_create_tempcups_table.up.sql

CREATE TABLE IF NOT EXISTS tempcups(
    id bigserial PRIMARY KEY,
    user_id integer NOT NULL,
    cups integer NOT NULL
);