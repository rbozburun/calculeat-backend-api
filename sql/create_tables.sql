CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS roles (
    id serial PRIMARY KEY NOT NULL,
    role_name character varying(30) NOT NULL,
    created_at timestamp default current_timestamp, 
    updated_at timestamp with time zone,
    deleted_on timestamp with time zone
);

CREATE TABLE IF NOT EXISTS users (
    id serial NOT NULL PRIMARY KEY,
    access_token character varying(64),
    password character varying(64),
    role_id integer REFERENCES roles(id),
    init_water_id integer NOT NULL,
    email character varying(75) NOT NULL UNIQUE,
    first_name character varying(40) NOT NULL,
    last_name character varying(40) NOT NULL,
    profile_img_filepath character varying(255),
    sign_in_provider character varying(75) NOT NULL,
    born DATE NOT NULL,
    gender CHAR NOT NULL,
    age integer NOT NULL,
    height integer NOT NULL,
    weight double precision NOT NULL,
    goal character varying(100),
    diseases text [],
    created_at timestamp default current_timestamp NOT NULL, 
    deleted_on timestamp with time zone,
    updated_at timestamp with time zone
);



CREATE TABLE IF NOT EXISTS connections (
    id serial NOT NULL PRIMARY KEY,
    client_user_id integer NOT NULL REFERENCES users(id),
    dietician_user_id integer NOT NULL REFERENCES users(id),
    is_accepted CHAR NOT NULL,
    created_at timestamp with time zone default current_timestamp NOT NULL,
    updated_at timestamp with time zone default current_timestamp,
    deleted_on timestamp with time zone
);

CREATE TABLE IF NOT EXISTS food_objects (
    id serial PRIMARY KEY NOT NULL,
    name character varying(75) NOT NULL,
    meal character varying(75) NOT NULL,
    user_id integer NOT NULL REFERENCES users(id),
    imgleft_filepath character varying(200) NOT NULL,
    imgright_filepath character varying(200) NOT NULL,
    imgtop_filepath character varying(200) NOT NULL,
    calorie integer NOT NULL,
    protein double precision NOT NULL,
    carbonhydrate double precision NOT NULL,
    fat double precision NOT NULL,
    magnesium double precision,
    calcium double precision,
    sodium double precision,
    iron double precision,
    vit_a integer,
    vit_c integer,
    vit_d integer,
    sugar double precision,
    saturated_fat double precision NOT NULL,
    unsaturated_fat double precision NOT NULL,
    gram double precision,
    created_at timestamp default current_timestamp NOT NULL, 
    updated_at timestamp with time zone,
    deleted_on timestamp with time zone
);

CREATE TABLE IF NOT EXISTS meetings (
    id serial PRIMARY KEY NOT NULL,
    client_user_id integer NOT NULL REFERENCES users(id),
    dietician_user_id integer NOT NULL REFERENCES users(id),
    meet_link character varying(100) NOT NULL,
    created_at timestamp default current_timestamp NOT NULL, 
    updated_at timestamp with time zone,
    deleted_on timestamp with time zone
);


CREATE TABLE IF NOT EXISTS messages (
    id serial PRIMARY KEY  NOT NULL,
    sender_id integer NOT NULL REFERENCES users(id),
    reciever_id integer NOT NULL REFERENCES users(id),
    message_text text NOT NULL,
    created_at timestamp default current_timestamp NOT NULL, 
    updated_at timestamp with time zone,
    deleted_on timestamp with time zone
);



CREATE TABLE IF NOT EXISTS sleep_objects (
    id serial PRIMARY KEY NOT NULL,
    user_id integer NOT NULL REFERENCES users(id),
    start_time timestamp with time zone,
    end_time timestamp with time zone,
    created_at timestamp default current_timestamp NOT NULL, 
    updated_at timestamp with time zone,
    deleted_on timestamp with time zone
);


CREATE TABLE IF NOT EXISTS water_objects (
    id serial PRIMARY KEY NOT NULL,
    user_id integer NOT NULL REFERENCES users(id),
    count integer NOT NULL,
    created_at timestamp default current_timestamp NOT NULL, 
    updated_at timestamp with time zone,
    deleted_on timestamp with time zone
);



