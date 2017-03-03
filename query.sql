create database niec;

use niec;

create table user(id bigint unsigned primary key auto_increment, username varchar(25) unique not null, email varchar(255) character set latin1 collate latin1_swedish_ci unique not null, password varchar(255) not null, bio text, dp varchar(255), created_at datetime not null, website varchar(255), verified bool not null, verifyhash varchar(255) not null, public bool not null);

create table article(id bigint unsigned primary key auto_increment, title varchar(255) not null, created_at datetime not null, text text not null, user_id bigint not null references user(id), public bool not null, draft bool not null);