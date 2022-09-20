-- +migrate Up
CREATE TABLE x
(
    id   UUID NOT NULL PRIMARY KEY,
    name TEXT NOT NULL
);
