This project is a ToDo application built using HTMX, Go, and PostgreSQL. I tried to make it as interactive as possible without writing a single line of js.


Setting Up the Project

Download the project. Create a .env file in the main directory and set the following fields:


DB_CONN_STR=

SMTP_HOST=

SMTP_PORT=

SMTP_USER=

SMTP_PASSWORD=

ID_ENCRYPTION=

CONFIRMATION_ENCRYPTION=


Database Setup: Ensure your PostgreSQL database includes the following tables: users, tasks, and email_confirmation. Use the SQL queries provided below to create these tables.


CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(32),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(32),
    registration_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    info TEXT,
    is_completed BOOLEAN NOT NULL DEFAULT false,
    has_deadline BOOLEAN NOT NULL DEFAULT false,
    deadline TIMESTAMP WITH TIME ZONE,
    importance SMALLINT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE email_confirmation (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(32),
    email VARCHAR(255) NOT NULL,
    password VARCHAR(32),
    confirmation_code VARCHAR(12)
);
