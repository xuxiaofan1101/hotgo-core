create database hotgo;

CREATE USER 'hotgo'@'%' IDENTIFIED BY 'hg123456.';

grant select,insert,update,delete,create on hotgo.* to hotgo;

flush  privileges ;