CREATE TABLE info
(
    id serial not null unique,
    user_id bigint not null,
    service_name varchar(255) not null,
    login varchar(255) not null,
    password varchar(255) not null
);