create extension "uuid-ossp" if not exists;

-- 1. Add your table to drop if exists
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS contributing;
DROP TABLE IF EXISTS interested;

-- 2. Create your table
CREATE TABLE users
(
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    image           TEXT,
    password        TEXT NOT NULL,
    profile_id      uuid UNIQUE      DEFAULT uuid_generate_v4(),
    deactivated     BOOLEAN DEFAULT FALSE,
    banned          BOOLEAN DEFAULT FALSE
);

CREATE TABLE projects
(
    id            uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title         TEXT NOT NULL,
    state         TEXT NOT NULL,
    /* tags       TEXT,*/
    user_id       uuid NOT NULL,
    start_date    TIMESTAMPTZ      NOT NULL,
    end_date      TIMESTAMPTZ      DEFAULT NULL,
    oneliner      TEXT,
    discussion_id TEXT,
    Logo          TEXT,
    CoverPhoto    TEXT,
    /* Media         Text*/
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE themes
(
    name        TEXT PRIMARY KEY,
    colour      TEXT NOT NULL,
    logo        TEXT NOT NULL,
    description TEXT
);

CREATE TABLE discussion
(
    proj_id    uuid,
    disc_num   SERIAL,
    creator    uuid NOT NULL,
    creation_date DATE,
    title      TEXT,
    text       TEXT,
    closed     BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (proj_id) REFERENCES projects(id),
    FOREIGN KEY (creator) REFERENCES users(id),
    PRIMARY KEY (proj_id, disc_num)
);

CREATE TABLE post
(
    proj_id     uuid,
    disc_num    INTEGER,
    post_num    SERIAL,
    creator     uuid NOT NULL,
    creation_date DATE,
    text        TEXT,
    pinned      BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (proj_id)  REFERENCES projects(id),
    FOREIGN KEY (disc_num) REFERENCES discussion(disc_num),
    FOREIGN KEY (creator)  REFERENCES users(id),
    PRIMARY KEY (proj_id, disc_num, post_num)
);

CREATE TABLE reaction
(
    reaction_type TEXT PRIMARY KEY,
    reaction_icon TEXT NOT NULL
);

CREATE TABLE follows
(
    follower_id   uuid,
    followed_id   uuid,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (follower_id,followed_id)
);

CREATE TABLE user_react_post
(
    user_id     uuid,
    proj_id     uuid,
    disc_num    INTEGER,
    post_num    INTEGER,

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

CREATE TABLE interested
(
    user_id       uuid,
    project_id    uuid,
    notifications BOOLEAN DEFAULT true,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id,project_id)
);

CREATE TABLE project_has_theme
(
    project_id    uuid,
    theme_name    TEXT,
    primary_theme BOOLEAN DEFAULT false,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    FOREIGN KEY (theme_name) REFERENCES themes(name) ON DELETE CASCADE,
    PRIMARY KEY (project_id, theme_name)
);

CREATE TABLE user_interested_theme
(
    user_id     uuid,
    theme_name  TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (theme_name) REFERENCES themes(name) ON DELETE CASCADE,
    PRIMARY KEY (user_id, theme_name)
);

-- 3. Give a few examples to be added into your table
-- Examples User
INSERT INTO users (first_name, last_name, email, image, password, profile_id, deactivated, banned)
VALUES ('Alexander', 'Bergholm', 'bergholm.alexander@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, password, profile_id, deactivated, banned)
VALUES ('John', 'Zhang', 'projincubator@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later_1', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, password, profile_id, deactivated, banned)
VALUES ('Kenrick', 'Yap', 'dicksaresocute69@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later_2', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, password, profile_id, deactivated, banned)
VALUES ('Test', 'Testing', 'testperson@gmail.com', 'sometest.com', 'test', 'test', FALSE, FALSE);
-- Examples Projects
-- Examples Follows
-- Examples Contributing
-- Examples Interested
-- 4. Check that they now exist in the database
select * from users
