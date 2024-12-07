CREATE TABLE operation_types (
    operation_type_id INT PRIMARY KEY NOT NULL,
    description TEXT NOT NULL
);

INSERT INTO operation_types (operation_type_id, description)
VALUES (1, 'Normal Purchase'),
       (2, 'Purchase with installments'),
       (3, 'Withdrawal'),
       (4, 'Credit Voucher');
