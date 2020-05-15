CREATE DATABASE telebot;
CREATE TABLE users (
  id 						SERIAL PRIMARY KEY,
  user_id 			INT NOT NULL,
  user_name 		VARCHAR(255),
  chat_id 			INT NOT NULL,
	first_name   	VARCHAR(255),
	last_name    	VARCHAR(255),
	language_code VARCHAR(255),
	is_bot       	BOOLEAN,
	created_at    	TIMESTAMP NOT NULL
);

CREATE TABLE sites (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  name VARCHAR(255),
  url VARCHAR(2083),
  request_timeout INT NOT NULL,
  response_status INT NOT NULL,
  description VARCHAR(1000),
  created_at TIMESTAMP NOT NULL
);

CREATE TABLE rules (
                     id SERIAL PRIMARY KEY,
                     interval INT NOT NULL,
                     status VARCHAR(255),
                     description VARCHAR(1000),
                     create_at TIMESTAMP NOT NULL
);
