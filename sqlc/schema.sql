DROP schema public cascade;
CREATE schema public;

CREATE TABLE accounts (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY KEY,
  inserted_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,
  updated_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,
  name text
);

CREATE TABLE mints (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY key,
  amount DECIMAL NOT NULL
);

CREATE TABLE spends (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY key,
  amount DECIMAL NOT NULL
);

CREATE TABLE transfers (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY key,
  amount DECIMAL NOT NULL NOT NULL,
  recipient BIGINT REFERENCES accounts NOT NULL
);

CREATE TABLE transactions (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY key,
  account BIGINT REFERENCES accounts NOT NULL,
  inserted_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,

  mint BIGINT REFERENCES mints UNIQUE,
  spend BIGINT REFERENCES spends UNIQUE,
  transfer BIGINT REFERENCES transfers UNIQUE,
  CONSTRAINT kind CHECK(num_nonnulls(mint, transfer, spend) = 1)
);
CREATE INDEX ON transactions (account);
CREATE INDEX ON transactions (inserted_at);

CREATE TABLE balances (
  account BIGINT REFERENCES accounts PRIMARY KEY NOT NULL,
  balance DECIMAL NOT NULL
);
