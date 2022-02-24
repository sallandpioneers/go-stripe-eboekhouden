create table customers (
    id BINARY(16) NOT NULL,
    stripe_id VARCHAR(20) NOT NULL,
    boekhouden_id BIGINT NOT NULL,
    boekhouden_code VARCHAR(15) NOT NULL
)ENGINE=INNODB DEFAULT COLLATE utf8mb4_unicode_ci;