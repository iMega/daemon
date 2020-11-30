CREATE DATABASE IF NOT EXISTS test;

use test;

CREATE TABLE test (
  title varchar(255) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

INSERT test (title) VALUES ('mysql1');
