create extension "uuid-ossp" if not exists;
-- 1. Add your table to drop if exists
DROP TABLE IF EXISTS users;
-- 2. Create your table
CREATE TABLE users
(
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    image           TEXT,
    password        TEXT NOT NULL,
    profile_id      TEXT UNIQUE,
    deactivated     BOOLEAN DEFAULT FALSE,
    banned          BOOLEAN DEFAULT FALSE
);

CREATE TABLE projects
(
    id            uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title         TEXT NOT NULL,
    state         TEXT NOT NULL,
    /* tags          TEXT,*/
    user_id       uuid NOT NULL,
    created_date  TIMESTAMPTZ      NOT NULL,
    end_date      TIMESTAMPTZ      DEFAULT NULL,
    oneliner      TEXT,
    discussion_id TEXT,
    Logo          Text,
    CoverPhoto    Text,
    /* Media         Text*/
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE follows
(
    followed_id   uuid,
    follower_id   uuid,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (followed_id,follower_id)
);

CREATE TABLE contributing
(
    user_id     uuid,
    project_id  uuid,
    is_admin    BOOLEAN DEFAULT false,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id,project_id)
);

CREATE TABLE intrested
(
    user_id       uuid,
    project_id    uuid,
    notifications BOOLEAN DEFAULT true,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id,project_id)
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