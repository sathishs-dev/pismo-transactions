CREATE TABLE transactions (
    transaction_id INT PRIMARY KEY,
    account_id INT NOT NULL REFERENCES accounts(account_id),
    operation_type_id INT NOT NULL REFERENCES operation_types(operation_type_id),
    amount DECIMAL NOT NULL,
    event_date TIMESTAMPTZ
);
