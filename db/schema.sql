CREATE TYPE job_status AS ENUM ('in-progress', 'succeeded', 'failed');

CREATE TABLE revisions (
  package      text       NOT NULL,
  revision     text       NOT NULL,
  hash         text       NOT NULL,
  status       job_status NOT NULL UNIQUE,
  last_indexed timestamp  NOT NULL,
  PRIMARY KEY(package, revision)
);