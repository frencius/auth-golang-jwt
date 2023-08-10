/**
  This is the SQL script that will be used to initialize the database schema.
  We will evaluate you based on how well you design your database.
  1. How you design the tables.
  2. How you choose the data types and keys.
  3. How you name the fields.
  In this assignment we will use PostgreSQL as the database.
  */

/** This is test table. Remove this table and replace with your own tables. */


CREATE TABLE "user"(
	"id"                    UUID PRIMARY KEY,
	full_name               VARCHAR (100) NOT NULL,
  phone_number            VARCHAR (13) UNIQUE NOT NULL,
  "password"              VARCHAR (255) NOT NULL,
  created_at              timestamptz		NOT NULL DEFAULT now(),
	updated_at              timestamptz		NOT NULL DEFAULT now(),
	created_by              varchar(100)	NOT NULL DEFAULT 'system'::character varying,
	updated_by              varchar(100)	NOT NULL DEFAULT 'system'::character varying
);

CREATE TABLE login(
	"id"                    UUID PRIMARY KEY,
	user_id                 UUID NOT NULL UNIQUE,
  success_counter         int NOT NULL DEFAULT 0,
  last_login              timestamptz		NOT NULL DEFAULT now(),
  created_at              timestamptz		NOT NULL DEFAULT now(),
	updated_at              timestamptz		NOT NULL DEFAULT now(),
	created_by              varchar(100)	NOT NULL DEFAULT 'system'::character varying,
	updated_by              varchar(100)	NOT NULL DEFAULT 'system'::character varying,
  CONSTRAINT fk_user_id FOREIGN KEY(user_id) REFERENCES "user"(id) ON DELETE NO ACTION
);