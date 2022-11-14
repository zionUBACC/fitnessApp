-- Filename: migrations/000002_create_tempsteps_table.up.sql

CREATE TABLE IF NOT EXISTS tempsteps(
    id bigserial PRIMARY KEY,
    user_id integer NOT NULL,
    steps integer NOT NULL
);
