-- +goose Up
-- +goose StatementBegin
CREATE TABLE cars
(
    id SERIAL PRIMARY KEY,
    reg_num VARCHAR(20) UNIQUE NOT NULL,
    mark VARCHAR(100) NOT NULL,
    model VARCHAR(100) NOT NULL,
    year INTEGER,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    owner_name VARCHAR(100) NOT NULL,
    owner_surname VARCHAR(100) NOT NULL,
    owner_patronymic VARCHAR(100)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cars;
-- +goose StatementEnd