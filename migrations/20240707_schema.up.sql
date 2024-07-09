CREATE TABLE IF NOT EXISTS users(
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	first_name varchar(100),
	second_name varchar(100),
	birthdate varchar(20),
	gender varchar(10),
	biography varchar(1000),
	city varchar(100)
);

CREATE TABLE IF NOT EXISTS user_data(
	user_id uuid NOT NULL,
	password_hash varchar NOT NULL,
	salt varchar NOT NULL
);

CREATE TABLE IF NOT EXISTS tokens(
	access_token uuid PRIMARY KEY NOT NULL,
	user_id uuid NOT NULL
);