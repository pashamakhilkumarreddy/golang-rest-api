CREATE TABLE IF NOT EXISTS products (
  id SERIAL PRIMARY KEY,
  name VARCHAR(250) NOT NULL,
  price NUMERIC(10, 2) NOT NULL
);
