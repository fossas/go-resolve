-- +goose Up
CREATE TABLE packages (
  import_path text NOT NULL,
  revision text NOT NULL,
  hash text NOT NULL, -- This is not unique because some packages will not change between revisions and therefore might have the same hash.
  version text,

  last_updated timestamp NOT NULL,

  PRIMARY KEY (import_path, revision)
);

-- +goose Down
DROP TABLE packages;
