-- +migrate Up
CREATE TABLE IF NOT EXISTS user_subs (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price INTEGER NOT NULL,
    user_id UUID NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL
);

-- Добавляем индексы для улучшения производительности
CREATE INDEX IF NOT EXISTS idx_user_subs_user_id ON user_subs(user_id);
CREATE INDEX IF NOT EXISTS idx_user_subs_service_name ON user_subs(service_name);
CREATE INDEX IF NOT EXISTS idx_user_subs_start_date ON user_subs(start_date);
CREATE INDEX IF NOT EXISTS idx_user_subs_end_date ON user_subs(end_date);

-- +migrate Down
DROP TABLE user_subs;