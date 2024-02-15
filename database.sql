/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE IF NOT EXISTS "users" (
  "id" BIGSERIAL NOT NULL PRIMARY KEY,
  "name" VARCHAR NOT NULL,
  "phone" VARCHAR NOT NULL,
  "password" VARCHAR NOT NULL,
  "created_at" TIMESTAMPTZ(0),
  "updated_at" TIMESTAMPTZ(0),
  UNIQUE ("phone")
);

CREATE TABLE IF NOT EXISTS "user_tokens" (
  "id" BIGSERIAL NOT NULL PRIMARY KEY,
  "user_id" BIGINT NOT NULL,
  "token" VARCHAR NOT NULL,
  "count_login" INT NOT NULL,
  "created_at" TIMESTAMPTZ(0),
  "updated_at" TIMESTAMPTZ(0),
  FOREIGN KEY ("user_id") REFERENCES "users" ("id")
);
