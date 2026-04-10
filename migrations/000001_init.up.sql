CREATE SCHEMA todoapp;

CREATE TABLE todoapp.users (
    id           SERIAL                PRIMARY KEY,
    version      BIGINT       NOT NULL DEFAULT 1,
    full_name    VARCHAR(100) NOT NULL CHECK(char_length(full_name) BETWEEN 3 AND 100),
    phone_number VARCHAR(15)           CHECK(
        phone_number ~ '^\+[0-9]+$'
        AND
        char_length(phone_number) BETWEEN 10 AND 15
    )
);

CREATE TABLE todoapp.tasks (
    id           SERIAL                  PRIMARY KEY,
    version      BIGINT         NOT NULL DEFAULT 1,
    title        VARCHAR(100)   NOT NULL CHECK(char_length(title) BETWEEN 1 AND 100),
    description  VARCHAR(1000)           CHECK(char_length(description) BETWEEN 1 AND 1000),
    complited    BOOLEAN        NOT NULL,
    is_important BOOLEAN        NOT NULL,
    is_in_my_day BOOLEAN        NOT NULL,
    created_at   TIMESTAMPTZ    NOT NULL,
    due_date     TIMESTAMPTZ,
    complited_at TIMESTAMPTZ,

    CHECK (
        (complited=FALSE AND complited_at IS NULL)
        OR
        (complited=TRUE AND complited_at IS NOT NULL AND complited_at >= created_at)    
    ),

    author_user_id INTEGER       NOT NULL REFERENCES todoapp.users(id)
);

CREATE TABLE todoapp.lists (
    id             SERIAL                PRIMARY KEY,
    version        BIGINT       NOT NULL DEFAULT 1,
    name           VARCHAR(100) NOT NULL CHECK(char_length(name) BETWEEN 1 and 100),
    author_user_id INTEGER      NOT NULL REFERENCES todoapp.users(id)
)