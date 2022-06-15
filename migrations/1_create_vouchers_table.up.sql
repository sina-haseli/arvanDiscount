CREATE TABLE vouchers
(
    id         serial    not null,
    code       varchar   NOT NULL,
    amount     INT       NOT NULL,
    usable     INT       NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);