CREATE TABLE IF NOT EXISTS organisations (
  id UUID PRIMARY KEY,
  name VARCHAR(50) NOT NULL,
  description VARCHAR(255)
);
