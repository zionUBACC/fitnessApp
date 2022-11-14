-- Filename: migrations/000002_create_dailyfitness_table.up.sql

CREATE TABLE IF NOT EXISTS dailyfitness(
    id bigserial PRIMARY KEY,
    user_id INTEGER NOT NULL,
    steps INTEGER NOT NULL,
    cups INTEGER NOT NULL,
    date timestamp(0) with time zone NOT NULL DEFAULT NOW()
);