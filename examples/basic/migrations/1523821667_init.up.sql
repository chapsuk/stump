CREATE TABLE users(
  id SERIAL NOT NULL PRIMARY KEY,
  name VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  rating int DEFAULT 0
);