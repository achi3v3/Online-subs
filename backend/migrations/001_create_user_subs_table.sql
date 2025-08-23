-- +migrate Up
CREATE TABLE IF NOT EXISTS user_subs (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE user_subs;