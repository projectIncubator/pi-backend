ALTER TABLE medias RENAME TO media;

ALTER TABLE users DROP COLUMN deactivated;

ALTER TABLE users DROP COLUMN banned;

CREATE TYPE USER_STATUS AS ENUM ('active', 'deactivated', 'banned');

ALTER TABLE users ADD COLUMN status USER_STATUS;

UPDATE users SET status = 'active';

ALTER TABLE users ALTER COLUMN status SET DEFAULT 'active';