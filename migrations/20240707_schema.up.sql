CREATE TABLE IF NOT EXISTS users(
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	first_name varchar(100),
	second_name varchar(100),
	birthdate varchar(20),
	gender varchar(10),
	biography varchar(1000),
	city varchar(100)
);

CREATE TABLE IF NOT EXISTS posts(
	id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id uuid NOT NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	content varchar
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

CREATE TABLE IF NOT EXISTS public.friends
(
    user_id uuid NOT NULL,
    friend_user_id uuid NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT now(),
    CONSTRAINT friends_pkey PRIMARY KEY (user_id, friend_user_id),
    CONSTRAINT friend_user_id FOREIGN KEY (friend_user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID,
    CONSTRAINT user_id FOREIGN KEY (user_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
)