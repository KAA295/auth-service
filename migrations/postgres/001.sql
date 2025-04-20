
-- +migrate Up
CREATE TABLE tokens(
  id SERIAL PRIMARY KEY,
  token VARCHAR NOT NULL,
  user_id UUID NOT NULL, 
  expiration_time TIMESTAMP,
  pair_id UUID NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS tokens;
