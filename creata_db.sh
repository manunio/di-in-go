#!/usr/bin/env bash

rm -rf test.db

sqlite3 test.db 'CREATE TABLE people(id INTEGER PRIMARY KEY ASC, name TEXT, age INTEGER);'
sqlite3 test.db 'INSERT INTO people (name, age) VALUES ("michael", 35);'
sqlite3 test.db 'INSERT INTO people (name, age) VALUES ("lucifer", 29);'