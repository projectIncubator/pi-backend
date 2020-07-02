CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- 1. Add your table to drop if exists

DROP TABLE IF EXISTS project_has_theme;
DROP TABLE IF EXISTS user_interested_theme;
DROP TABLE IF EXISTS themes;
DROP TABLE IF EXISTS discussion_has_media;
DROP TABLE IF EXISTS post_has_media;
DROP TABLE IF EXISTS project_has_media;
DROP TABLE IF EXISTS user_react_post;
DROP TABLE IF EXISTS medias;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS discussions;
DROP TABLE IF EXISTS contributing;
DROP TABLE IF EXISTS interested;
DROP TABLE IF EXISTS follows;
DROP TABLE IF EXISTS sidebar_modules;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;

-- 2. Create your table
CREATE TABLE users
(
    id              uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    id_token        TEXT NOT NULL UNIQUE,
    first_name      TEXT NOT NULL,
    last_name       TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    image           TEXT DEFAULT 'placeholder_url',
    profile_id      TEXT UNIQUE, /* TODO: set = to id if null*/
    deactivated     BOOLEAN DEFAULT FALSE,
    banned          BOOLEAN DEFAULT FALSE
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
    description TEXT DEFAULT 'insert description here'
);

CREATE TABLE medias
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

-- 3. Give a few examples to be added into your table

-- Examples User
INSERT INTO users (id_token, id, first_name, last_name, email, image, profile_id, deactivated, banned)
VALUES ('00000', 'beee0bf1-5b7e-4f21-bcea-17ae7c45b18c', 'Alexander', 'Bergholm', 'bergholm.alexander@gmail.com', 'someurllater.com', 'purity', FALSE, FALSE);
INSERT INTO users (id_token, id, first_name, last_name, email, image, profile_id, deactivated, banned)
VALUES ('11111', 'de8ccc40-9372-411d-8497-a0becf01eff0', 'John', 'Zhang', 'projincubator@gmail.com', 'someurllater.com', 'joker', FALSE, FALSE);
INSERT INTO users (id_token, id, first_name, last_name, email, image, profile_id, deactivated, banned)
VALUES ('10101', '591564c5-45e4-43aa-bc5f-c37a50a6e7d8', 'Kenrick', 'Yap', 'dicksaresocute69@gmail.com', 'someurllater.com', 'KayYep', FALSE, FALSE);
INSERT INTO users (id_token, id,  first_name, last_name, email, image, profile_id, deactivated, banned)
VALUES ('01010', 'c9208c94-3c0f-416e-997d-a4d23cb3016f', 'Test', 'Testing', 'testperson@gmail.com', 'sometest.com', 'test', FALSE, FALSE);

-- Examples Projects
INSERT INTO projects (id, title, creator, oneliner)
VALUES ('6e2f8633-8743-4214-9115-d4e65a76b113', 'COVID vaccine', 'beee0bf1-5b7e-4f21-bcea-17ae7c45b18c', 'We aim cure COVID.');
INSERT INTO projects (id, title, creator, oneliner)
VALUES ('f34130ab-a785-47f1-a6e2-3fae4d4b0d07','Street photography', 'de8ccc40-9372-411d-8497-a0becf01eff0', 'Revealing inequity through art');
INSERT INTO projects (id, title, creator, oneliner)
VALUES ('260a1370-c968-490c-a623-d0f6d130990f', 'Pancakes', '591564c5-45e4-43aa-bc5f-c37a50a6e7d8', 'Believe that pancakes can fly.');

-- Examples Themes
INSERT INTO Themes (name, logo)
VALUES ('No Poverty', 'placeholder_url');
INSERT INTO Themes (name, logo)
VALUES ('Zero Hunger', 'placeholder_url');
INSERT INTO Themes (name, logo)
VALUES ('Good Health and Well-being', 'placeholder_url');
INSERT INTO Themes (name, logo)
VALUES ('Quality Education', 'placeholder_url');

-- Examples Follows
/* John follows Alex */
INSERT INTO follows (follower_id, followed_id)
VALUES ('de8ccc40-9372-411d-8497-a0becf01eff0', 'beee0bf1-5b7e-4f21-bcea-17ae7c45b18c');

/* Kenrick follows Alex */
INSERT INTO follows (follower_id, followed_id)
VALUES ('591564c5-45e4-43aa-bc5f-c37a50a6e7d8', 'beee0bf1-5b7e-4f21-bcea-17ae7c45b18c');


/* Kenrick follows John */
INSERT INTO follows (follower_id, followed_id)
VALUES ('591564c5-45e4-43aa-bc5f-c37a50a6e7d8', 'de8ccc40-9372-411d-8497-a0becf01eff0');

/* John follows Kenrick */
INSERT INTO follows (follower_id, followed_id)
VALUES ('de8ccc40-9372-411d-8497-a0becf01eff0','591564c5-45e4-43aa-bc5f-c37a50a6e7d8');

-- Examples Contributing

/* John contributes to COVID vaccine */
INSERT INTO contributing (user_id, project_id, is_admin)
VALUES ('de8ccc40-9372-411d-8497-a0becf01eff0','6e2f8633-8743-4214-9115-d4e65a76b113', true);

/* John contributes to Pancakes */
INSERT INTO contributing (user_id, project_id, is_admin)
VALUES ('de8ccc40-9372-411d-8497-a0becf01eff0','260a1370-c968-490c-a623-d0f6d130990f', true);

/* Kenrick contributes to COVID vaccine */
INSERT INTO contributing (user_id, project_id, is_admin)
VALUES ('591564c5-45e4-43aa-bc5f-c37a50a6e7d8','6e2f8633-8743-4214-9115-d4e65a76b113', false);

-- Examples Interested

/* John interested in COVID vaccine */
INSERT INTO interested (user_id, project_id)
VALUES ('de8ccc40-9372-411d-8497-a0becf01eff0','6e2f8633-8743-4214-9115-d4e65a76b113');

/* John interested in Pancakes */
INSERT INTO interested (user_id, project_id)
VALUES ('de8ccc40-9372-411d-8497-a0becf01eff0','260a1370-c968-490c-a623-d0f6d130990f');

/* Kenrick interested COVID vaccine */
INSERT INTO interested (user_id, project_id)
VALUES ('591564c5-45e4-43aa-bc5f-c37a50a6e7d8','6e2f8633-8743-4214-9115-d4e65a76b113');

-- 4. Check that they now exist in the database