create extension "uuid-ossp" if not exists;
-- 1. Add your table to drop if exists
DROP TABLE IF EXISTS users;
-- 2. Create your table
CREATE TABLE users
(
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    email           TEXT DEFAULT '' UNIQUE,
    image           TEXT,
    password        TEXT,
    profile_id      TEXT UNIQUE,
    deactivated     BOOLEAN DEFAULT FALSE,
    banned          BOOLEAN DEFAULT FALSE
);
-- 3. Give a few examples to be added into your table
INSERT INTO users (first_name, last_name, email, image, password, profile_id, deactivated, banned)
VALUES ('Alexander', 'Bergholm', 'bergholm.alexander@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, password, profile_id, deactivated, banned)
VALUES ('John', 'Zhang', 'projincubator@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later_1', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, password, profile_id, deactivated, banned)
VALUES ('Kenrick', 'Yap', 'dicksaresocute69@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later_2', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, password, profile_id, deactivated, banned)
VALUES ('Test', 'Testing', 'testperson@gmail.com', 'sometest.com', 'test', 'test', FALSE, FALSE);
-- 4. Check that they now exist in the database
select * from users