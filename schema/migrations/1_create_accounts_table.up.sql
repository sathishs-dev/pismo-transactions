CREATE TABLE accounts (
    account_id SERIAL PRIMARY KEY,
    document_number VARCHAR(255) NOT NULL UNIQUE
);
