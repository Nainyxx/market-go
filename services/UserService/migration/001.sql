CREATE TABLE IF NOT EXISTS users (
    user_id     UUID         PRIMARY KEY,
    name        VARCHAR(64)  NOT NULL,
    surname     VARCHAR(64)  NOT NULL DEFAULT '',
    image_url   TEXT         NOT NULL DEFAULT '',
    rate        REAL         NOT NULL DEFAULT 0,
    rates_count INT          NOT NULL DEFAULT 0,
    balance     BIGINT       NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL,
    updated_at  TIMESTAMPTZ  NOT NULL,

    CONSTRAINT users_balance_non_negative CHECK (balance >= 0),
    CONSTRAINT users_rate_range           CHECK (rate >= 0 AND rate <= 5),
    CONSTRAINT users_rates_count_non_neg  CHECK (rates_count >= 0)
);
