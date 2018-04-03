CREATE KEYSPACE IF NOT EXISTS meower WITH REPLICATION = {
  'class': 'SimpleStrategy',
  'replication_factor':'1'
};

USE meower;

CREATE TABLE IF NOT EXISTS users (
  google_id TEXT,
  name TEXT,
  PRIMARY KEY (google_id, name)
);

CREATE TABLE IF NOT EXISTS meows (
  id BIGINT,
  user_name TEXT,
  body TEXT,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS timeline (
  user_name TEXT,
  created_at TIMEUUID,
  meow_id BIGINT,
  PRIMARY KEY (user_name, created_at)
) WITH CLUSTERING ORDER BY (created_at DESC);
