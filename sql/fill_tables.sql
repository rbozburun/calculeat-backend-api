CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

--Fill roles table 
INSERT INTO roles(id,role_name) VALUES(1, 'dietician');
INSERT INTO roles(id,role_name) VALUES(2, 'client');


--Fill users table
INSERT INTO users (born, gender, role_id, email, first_name, last_name, age, height, weight, goal, diseases, password) VALUES
    ('2016-06-23','M', 2, 'johndoe@example.com', 'John', 'Doe', 28, 175, 85, 'Lose wight', '{"Asthma", "migraine"}', '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08'),
    ('2016-06-23','M', 1, 'arif@diet.com', 'Arif', 'Sezen', 23, 163, 52, NULL, NULL, '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92'),
    ('2016-06-23','M', 2, 'neymar@example.com', 'Neymar', 'Doe', 29, 171, 85, 'Lose wight', '{"Asthma", "migraine"}', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92'),
    ('2016-06-23','F', 1, 'alice_jen@diet.com', 'Alice', 'Jen', 23, 178, 83, NULL, NULL, '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08'),
    ('2016-06-23','O', 2, 'kevin@gmail.com', 'Kevin', 'Clark', 45, 170, 56, 'Gain weight', NULL, '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08'),
    ('2016-06-23','F', 2, 'Hatice@gmail.com', 'Hatice', 'Uslu', 20, 170, 60, NULL, NULL, '9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08'),
    ('2016-06-23','M', 2, 'jackson_jane@yahoo.com', 'Jane', 'Jackson', 44, 160, 75, 'Lose weight','{"Asthma", "migraine"}', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92'),
    ('2016-06-23','M', 2, 'mark@yahoo.com', 'Mark', 'Philp', 44, 160, 75, 'Lose weight','{"Asthma"}', '8d969eef6ecad3c29a3a629280e686cf0c3f5d5a86aff3ca12020c923adc6c92');

UPDATE users
SET profile_img_filepath = concat('/images/profile_images/',users.id,'/');

INSERT INTO food_objects 
(meal, name, user_id, imgleft_filepath, imgright_filepath, imgtop_filepath, calorie, protein, 
carbonhydrate, fat, magnesium, calcium, sodium, iron, vit_a, vit_c, 
sugar, saturated_fat, unsaturated_fat, gram, vit_d)
VALUES 
('breakfast','pizza', 1, '/images/food_images/', '/images/food_images/', '/images/food_images', 241, 13.4, 
34.2, 10.2, '2017-07-21T17:32:28.000Z', 2.3, NULL, 1.2, 1.1, 3, NULL, 
25.3, 6.0, 4.2, 262.0, 1.2),
('brunch','chicken', 2, '/images/food_images/', '/images/food_images/', '/images/food_images', 241, 23.4, 
34.2, 10.2, '2017-07-21T17:32:28.000Z', 2.3, NULL, 1.2, 1.1, 3, NULL, 
25.3, 6.0, 4.2, 103.0, 1.2),
('lunch','rice', 3, '/images/food_images/', '/images/food_images/', '/images/food_images', 221, 13.4, 
34.2, 10.2, '2017-07-21T17:32:28.000Z', 2.3, NULL, 1.2, 1.1, 3, NULL, 
25.3, 6.0, 4.2, 204.3, 1.2);

UPDATE food_objects
SET imgleft_filepath = concat('/images/food_objects/',food_objects.id,'/left/');


UPDATE food_objects
SET imgright_filepath = concat('/images/food_objects/',food_objects.id,'/right/');

UPDATE food_objects
SET imgtop_filepath = concat('/images/food_objects/',food_objects.id,'/top/');






