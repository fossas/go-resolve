-- +goose Up
CREATE TYPE vcs AS ENUM ('git', 'svn', 'hg', 'bzr');

CREATE TABLE repositories (
  -- `repository_id` is the primary key instead of `(type, url)` to make JOINs
  -- easier to work with.
  repository_id serial PRIMARY KEY,

  type vcs NOT NULL,
  url text NOT NULL,
  UNIQUE (type, url)
);

CREATE TABLE revisions (
  revision_id text PRIMARY KEY,
  repository_id int NOT NULL REFERENCES repositories(repository_id),
  timestamp timestamp NOT NULL
);

-- A package is an _instance_ of a package (i.e. one specific revision).
CREATE TABLE packages (
  package_id serial PRIMARY KEY,
  import_path text NOT NULL,

  repository_id int NOT NULL REFERENCES repositories(repository_id),
  revision_id text NOT NULL REFERENCES revisions(revision_id),

  -- `hash` is not unique because some packages will not change between
  -- revisions and therefore multiple packages may have the same hash.
  hash text NOT NULL,
  version text,

  last_updated timestamp NOT NULL,

  UNIQUE (import_path, revision_id)
);

-- +goose Down
DROP TABLE packages;
DROP TABLE revisions;
DROP TABLE repositories;
DROP TYPE vcs;
