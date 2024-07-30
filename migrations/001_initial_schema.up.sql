CREATE TABLE IF NOT EXISTS links (
	id serial PRIMARY KEY,
	url varchar(255) not null,
	code varchar(255) not null,
	created_at timestamp default now() not null
);
