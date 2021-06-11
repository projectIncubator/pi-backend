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