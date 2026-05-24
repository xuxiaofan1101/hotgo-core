create database hotgo;

CREATE USER 'hotgo'@'%' IDENTIFIED BY '123456';

grant select,insert,update,delete,create on hotgo.* to hotgo;

flush  privileges;
