--user
--STORING PLAIN TEXT PASSWORD VERY BAD DO NOT DO EVER
--WILL CHANGE LATER
PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE user (
  id integer primary key autoincrement,
  username varchar(100),
  password varchar(1000),
  email varchar(100)
);
INSERT INTO "user" VALUES(1,'david','temp','test@test.com');


--category
CREATE TABLE category(
  id integer primary key autoincrement,
  name varchar(1000) not null,
  user_id references user(id)
);

INSERT INTO "category" VALUES(1,'TaskApp',1);
COMMIT;

--status
PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE status (
  id integer primary key autoincrement,
  status varchar(50) not null
);
INSERT INTO "status" VALUES(1,'COMPLETE');
INSERT INTO "status" VALUES(2,'PENDING');
INSERT INTO "status" VALUES(3,'DELETED');
COMMIT;

--task

PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE task (
  id integer primary key autoincrement,
  title varchar(100),
  content text,
  created_date timestamp,
  last_modified_at timestamp,
  finish_date timestamp,
  priority integer,
  cat_id references category(id),
  task_status_id references status(id),
  due_date timestamp,
  user_id references user(id),
  hide int
);

CREATE TABLE comments(id integer primary key autoincrement, content ntext, taskID references task(id), created datetime, user_id references user(id));

CREATE TABLE files(name varchar(1000) not null, autoName varchar(255) not null, user_id references user(id), created_date timestamp);

COMMIT;
