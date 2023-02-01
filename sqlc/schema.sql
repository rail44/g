DROP schema public cascade;
CREATE schema public;

CREATE TABLE accounts (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY key,
  inserted_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,
  updated_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,
  name text
);

CREATE TABLE mints (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY key,
  amount DECIMAL NOT NULL,
  account BIGINT REFERENCES accounts
);

CREATE TABLE transfers (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY key,
  amount DECIMAL NOT NULL,
  from_account BIGINT REFERENCES accounts,
  to_account BIGINT REFERENCES accounts
);

CREATE TABLE transactions (
  id BIGINT generated BY DEFAULT AS IDENTITY PRIMARY key,
  inserted_at TIMESTAMP WITH TIME zone DEFAULT timezone('utc':: text, now()) NOT NULL,
  mint BIGINT REFERENCES mints UNIQUE,
  transfer BIGINT REFERENCES transfers UNIQUE,
  CONSTRAINT kind CHECK(num_nonnulls(mint, transfer) = 1)
);
