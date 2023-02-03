DROP schema public cascade;
CREATE schema public;

CREATE TABLE accounts (
  id SERIAL PRIMARY KEY,
  inserted_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,
  updated_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,
  name text
);

CREATE TABLE mints (
  id SERIAL PRIMARY KEY,
  amount DECIMAL NOT NULL
);

CREATE TABLE spends (
  id SERIAL PRIMARY KEY,
  amount DECIMAL NOT NULL
);

CREATE TABLE transfers (
  id SERIAL PRIMARY KEY,
  amount DECIMAL NOT NULL NOT NULL,
  recipient SERIAL REFERENCES accounts NOT NULL
);

CREATE TABLE transactions (
  id SERIAL PRIMARY KEY,
  account SERIAL REFERENCES accounts NOT NULL,
  inserted_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,

  mint SERIAL REFERENCES mints UNIQUE,
  spend SERIAL REFERENCES spends UNIQUE,
  transfer SERIAL REFERENCES transfers UNIQUE,
  CONSTRAINT kind CHECK(num_nonnulls(mint, transfer, spend) = 1)
);
CREATE INDEX ON transactions (account);
CREATE INDEX ON transactions (inserted_at);

CREATE TABLE balances (
  account SERIAL REFERENCES accounts PRIMARY KEY,
  balance DECIMAL NOT NULL
);
