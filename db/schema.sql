CREATE TABLE revisions (
  package        text       NOT NULL,
  revision       text       NOT NULL,
  hash           text       NOT NULL,
  PRIMARY KEY(package, revision)
);