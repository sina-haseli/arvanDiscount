CREATE TABLE redeemed_voucher
(
    id         serial    not null,
    user_id    VARCHAR    NOT NULL,
    voucher_id INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT voucher_id_user_id UNIQUE (voucher_id, user_id),
    PRIMARY KEY (id)
);