CREATE TABLE users (
    user_id       UUID PRIMARY KEY,
    phone         VARCHAR(20)  NOT NULL UNIQUE,
    email         VARCHAR(254) NOT NULL UNIQUE,
    password_hash TEXT         NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT now()
);
