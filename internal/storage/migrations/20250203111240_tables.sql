-- +goose Up
CREATE TABLE promocodes (
                            id SERIAL PRIMARY KEY,
                            code VARCHAR(255) UNIQUE NOT NULL,
                            reward INT NOT NULL,
                            max_uses INT NOT NULL,
                            uses INT DEFAULT 0,
                            created_at TIMESTAMP

);
CREATE TABLE rewards (
                         id SERIAL PRIMARY KEY,
                         player_id INT NOT NULL,
                         promocode_id INT NOT NULL,
                         reward INT NOT NULL,
                         applied_at TIMESTAMP,
                         FOREIGN KEY (promocode_id) REFERENCES promocodes(id)
);
-- +goose Down
DROP TABLE IF EXISTS rewards;
DROP TABLE IF EXISTS promocodes;

