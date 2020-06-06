create extension "uuid-ossp" if not exists;

-- 1. Add your table to drop if exists

DROP TABLE IF EXISTS project_has_theme;
DROP TABLE IF EXISTS user_interested_theme;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS medias;
DROP TABLE IF EXISTS discussion_has_media;
DROP TABLE IF EXISTS post_has_media;
DROP TABLE IF EXISTS user_react_post;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS contributing;
DROP TABLE IF EXISTS discussions;
DROP TABLE IF EXISTS interested;
DROP TABLE IF EXISTS users;

-- 2. Create your table
CREATE TABLE users
(
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_token        TEXT NOT NULL UNIQUE,
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    image           TEXT,
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
    FOREIGN KEY (user_id) REFERENCES users (id)
);

CREATE TABLE themes
(
    name        TEXT PRIMARY KEY,
    colour      TEXT NOT NULL,
    logo        TEXT NOT NULL,
    description TEXT
);

CREATE TABLE medias
(
    url             TEXT PRIMARY KEY,
    file_name       TEXT,
    uploader        uuid,
    date_uploaded   TIMESTAMP DEFAULT current_timestamp,
    FOREIGN KEY (uploader) REFERENCES users(id)
);

CREATE TABLE discussions
(
    proj_id    uuid,
    disc_num   SERIAL,
    creator    uuid NOT NULL,
    creation_date TIMESTAMP DEFAULT current_timestamp,
    title      TEXT,
    text       TEXT,
    closed     BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (proj_id) REFERENCES projects(id),
    FOREIGN KEY (creator) REFERENCES users(id),
    PRIMARY KEY (proj_id, disc_num)
);

CREATE TABLE posts
(
    proj_id     uuid,
    disc_num    INTEGER,
    post_num    SERIAL,
    creator     uuid NOT NULL,
    creation_date TIMESTAMP DEFAULT current_timestamp,
    text        TEXT,
    pinned      BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (proj_id, disc_num) REFERENCES discussions(proj_id, disc_num),
    FOREIGN KEY (creator)  REFERENCES users(id),
    PRIMARY KEY (proj_id, disc_num, post_num)
);

CREATE TABLE discussion_has_media
(
    proj_id    uuid,
    disc_num   INTEGER,
    media_url  TEXT,
    FOREIGN KEY (proj_id, disc_num) REFERENCES discussions(proj_id, disc_num),
    FOREIGN KEY (media_url) REFERENCES medias(url),
    PRIMARY KEY (proj_id, disc_num, media_url)
);

CREATE TABLE post_has_media
(
    proj_id    uuid,
    disc_num   INTEGER,
    post_num   INTEGER,
    media_url  TEXT,
    FOREIGN KEY (proj_id, disc_num) REFERENCES discussions(proj_id, disc_num),
    FOREIGN KEY (media_url) REFERENCES medias(url),
    PRIMARY KEY (proj_id, disc_num, media_url)
);

CREATE TABLE user_react_post
(
    user_id     uuid,
    proj_id     uuid,
    disc_num    INTEGER,
    post_num    INTEGER,
    reaction    TEXT NOT NULL,
    FOREIGN KEY (proj_id, disc_num, post_num) REFERENCES posts(proj_id, disc_num, post_num),
    PRIMARY KEY (user_id, proj_id, disc_num, post_num)
);

CREATE TABLE follows
(
    follower_id   uuid,
    followed_id   uuid,
    FOREIGN KEY (follower_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (followed_id) REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (follower_id,followed_id)
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
INSERT INTO users (first_name, last_name, email, image, profile_id, deactivated, banned)
VALUES ('Alexander', 'Bergholm', 'bergholm.alexander@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, profile_id, deactivated, banned)
VALUES ('John', 'Zhang', 'projincubator@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later_1', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, profile_id, deactivated, banned)
VALUES ('Kenrick', 'Yap', 'dicksaresocute69@gmail.com', 'someurllater.com', 'pass', 'some_uuid_later_2', FALSE, FALSE);
INSERT INTO users (first_name, last_name, email, image, profile_id, deactivated, banned)
VALUES ('Test', 'Testing', 'testperson@gmail.com', 'sometest.com', 'test', 'test', FALSE, FALSE);
-- Examples Projects
-- Examples Follows
-- Examples Contributing
-- Examples Interested
-- 4. Check that they now exist in the database
select * from users
