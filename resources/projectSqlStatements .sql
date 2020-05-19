create extension "uuid-ossp";
-- 1. Add your table to drop if exists
DROP TABLE IF EXISTS projects;
-- 2. Create your table
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
   /* members       Text,*/
    Logo          Text,
    CoverPhoto    Text,
   /* Media         Text*/
    FOREIGN KEY (user_id) REFERENCES users (id)
);
-- 3. Give a few examples to be added into your table
INSERT INTO projects (title, state, user_id, created_date, end_date, oneliner, discussion_id, Logo, CoverPhoto)
VALUES ('human rights', 'onGoing', 'cbb73ab2-e9d0-403d-b4e7-12624900a6ac', '2016-06-22 19:10:25-07', '2016-06-22 19:10:25-07', 'SomeOneliner',
 'discussionID','logo','CoverPhoto');

-- 4. Check that they now exist in the database
select *
from users