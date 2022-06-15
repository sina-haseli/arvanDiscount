CREATE TABLE redeemed_voucher
(
    id         serial    not null,
    user_id    INT       NOT NULL,
    voucher_id INT       NOT NULL,
    step       INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT voucher_id_user_id UNIQUE (voucher_id, user_id),
    CONSTRAINT limiter_unique UNIQUE (voucher_id, step),
    PRIMARY KEY (id)
);