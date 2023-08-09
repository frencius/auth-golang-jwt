/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */
CREATE TABLE test (
	id serial PRIMARY KEY,
	name VARCHAR ( 50 ) UNIQUE NOT NULL
);

INSERT INTO test (name) VALUES ('test1');
INSERT INTO test (name) VALUES ('test2');


CREATE TABLE user (
	id                      UUID PRIMARY KEY,
	full_name               VARCHAR (100) NOT NULL,
  phone_number            VARCHAR (15) UNIQUE NOT NULL,
  password                VARCHAR (255) NOT NULL,
  success_login_counter   int NOT NULL DEFAULT 0,
  created_at              timestamptz		NOT NULL DEFAULT now(),
	updated_at              timestamptz		NOT NULL DEFAULT now(),
	created_by              varchar(100)	NOT NULL DEFAULT 'system'::character varying,
	updated_by              varchar(100)	NOT NULL DEFAULT 'system'::character varying
);