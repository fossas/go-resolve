-- +goose Up
CREATE TABLE jobs (
  job_id serial PRIMARY KEY,
  faktory_jid text NOT NULL,

  started timestamp NOT NULL,
  completed timestamp
);

CREATE TABLE logs (
  log_id serial PRIMARY KEY,
  job_id int NOT NULL REFERENCES jobs(job_id),
  entry text NOT NULL,
  timestamp timestamp NOT NULL
);

ALTER TABLE packages ADD COLUMN job_id int NOT NULL REFERENCES jobs(job_id);

-- +goose Down
ALTER TABLE packages DROP COLUMN job_id;
DROP TABLE logs;
DROP TABLE jobs;
