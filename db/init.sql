-- Создание таблицы Users
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       amount INT DEFAULT 1000
);
CREATE INDEX IF NOT EXISTS idx_username ON users (username);

-- Создание таблицы Items
CREATE TABLE items (
                       id SERIAL PRIMARY KEY,
                       name VARCHAR(255) NOT NULL,
                       price INT NOT NULL
);

-- Создание таблицы Inventory
CREATE TABLE inventory (
                           id SERIAL PRIMARY KEY,
                           user_id INT NOT NULL,
                           item_id INT NOT NULL,
                           quantity INT DEFAULT 0,
                           FOREIGN KEY (user_id) REFERENCES users(id),
                           FOREIGN KEY (item_id) REFERENCES items(id)
);
CREATE UNIQUE INDEX idx_user_item ON inventory (user_id, item_id);


-- Создание таблицы CoinHistory
CREATE TABLE coin_history (
                             id SERIAL PRIMARY KEY,
                             --type VARCHAR(50) NOT NULL CHECK (type IN ('received', 'sent')),
                             from_user INT,
                             to_user INT,
                             amount INT NOT NULL,
                             FOREIGN KEY (from_user) REFERENCES users(id),
                             FOREIGN KEY (to_user) REFERENCES users(id)
);

INSERT INTO items(name, price) VALUES ('t-shirt', 80);
INSERT INTO items(name, price) VALUES ('cup', 20);
INSERT INTO items(name, price) VALUES ('book', 50);
INSERT INTO items(name, price) VALUES ('pen', 10);
INSERT INTO items(name, price) VALUES ('powerbank', 200);
INSERT INTO items(name, price) VALUES ('hoody', 300);
INSERT INTO items(name, price) VALUES ('umbrella', 200);
INSERT INTO items(name, price) VALUES ('socks', 10);
INSERT INTO items(name, price) VALUES ('wallet', 50);
INSERT INTO items(name, price) VALUES ('pink-hoody', 500);

-- Триггер и функция для него
-- CREATE OR REPLACE FUNCTION add_user_coins()
--     RETURNS TRIGGER AS $$
-- BEGIN
--     INSERT INTO coins (user_id, amount)
--     VALUES (NEW.id, 1000);
--     RETURN NEW;
-- END;
-- $$ LANGUAGE plpgsql;
--
-- CREATE TRIGGER trigger_add_user_coins
--     AFTER INSERT ON Users
--     FOR EACH ROW
-- EXECUTE FUNCTION add_user_coins();
--

