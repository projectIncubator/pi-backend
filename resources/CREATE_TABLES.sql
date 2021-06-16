CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE USER_STATUS AS ENUM ('active', 'deactivated', 'banned');

CREATE TABLE users
(
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_token        TEXT NOT NULL UNIQUE,
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    image           TEXT DEFAULT 'placeholder_url',
    profile_id      TEXT UNIQUE, /* TODO: set = to id if null*/
    bio             TEXT NOT NULL DEFAULT '',
    status          USER_STATUS DEFAULT 'active'
);

CREATE TABLE projects
(
    id            uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title         TEXT NOT NULL,
    state         TEXT DEFAULT 'idea',
    /* tags       TEXT,*/
    creator       uuid NOT NULL,
    start_date    TIMESTAMP DEFAULT current_timestamp,
    end_date      TIMESTAMP DEFAULT NULL,
    oneliner      TEXT DEFAULT 'insert description here',
    logo          TEXT DEFAULT 'placeholder_url', /* TODO make logo be media */
    cover_photo   TEXT DEFAULT 'placeholder_url', /* TODO make logo be media */
    FOREIGN KEY (creator) REFERENCES users (id)
);

CREATE TABLE project_has_media
(
    project_id              uuid,
    media           TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id,media)
);

CREATE TABLE themes
(
    name        TEXT PRIMARY KEY,
    logo        TEXT NOT NULL,
    description TEXT DEFAULT 'Insert description here'
);

CREATE TABLE media
(
    url             TEXT PRIMARY KEY,
    file_name       TEXT, /* TODO: set = to id if null*/
    uploader        uuid NOT NULL ,
    date_uploaded   TIMESTAMP DEFAULT current_timestamp,
    FOREIGN KEY (uploader) REFERENCES users(id)
);

CREATE TABLE discussions
(
    proj_id    uuid,
    disc_num   SERIAL,
    creator    uuid NOT NULL,
    creation_date TIMESTAMP DEFAULT current_timestamp,
    title      TEXT NOT NULL,
    text       TEXT DEFAULT '',
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
    text        TEXT NOT NULL,
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
    FOREIGN KEY (media_url) REFERENCES media(url),
    PRIMARY KEY (proj_id, disc_num, media_url)
);

CREATE TABLE post_has_media
(
    proj_id    uuid,
    disc_num   INTEGER,
    post_num   INTEGER,
    media_url  TEXT,
    FOREIGN KEY (proj_id, disc_num) REFERENCES discussions(proj_id, disc_num),
    FOREIGN KEY (media_url) REFERENCES media(url),
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

CREATE TABLE sidebar_modules
(
    project_id      uuid,
    index           INTEGER,
    module_type     TEXT NOT NULL,
    content         TEXT NOT NULL,
    FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, index)
);

CREATE TABLE user_interested_theme
(
    user_id     uuid,
    theme_name  TEXT,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (theme_name) REFERENCES themes(name) ON DELETE CASCADE,
    PRIMARY KEY (user_id, theme_name)
);